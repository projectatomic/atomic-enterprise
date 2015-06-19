package templates

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	kapi "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/rest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/apiserver"
	kclient "github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/master"
	"github.com/GoogleCloudPlatform/kubernetes/plugin/pkg/admission/admit"

	"github.com/projectatomic/appinfra-next/pkg/api/latest"
	osclient "github.com/projectatomic/appinfra-next/pkg/client"
	templateapi "github.com/projectatomic/appinfra-next/pkg/template/api"
	templateregistry "github.com/projectatomic/appinfra-next/pkg/template/registry"
)

func walkJSONFiles(inDir string, fn func(name, path string, data []byte)) error {
	err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != inDir {
			return filepath.SkipDir
		}
		name := filepath.Base(path)
		ext := filepath.Ext(name)
		if ext != "" {
			name = name[:len(name)-len(ext)]
		}
		if !(ext == ".json" || ext == ".yaml") {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		fn(name, path, data)
		return nil
	})
	return err
}

func TestTemplateTransformationFromConfig(t *testing.T) {
	osMux := http.NewServeMux()
	server := httptest.NewServer(osMux)
	defer server.Close()

	osClient := osclient.NewOrDie(&kclient.Config{Host: server.URL, Version: latest.Version})

	storage := map[string]rest.Storage{
		"processedTemplates": templateregistry.NewREST(),
	}
	for k, v := range storage {
		delete(storage, k)
		storage[strings.ToLower(k)] = v
	}

	interfaces, _ := latest.InterfacesFor(latest.Version)
	handlerContainer := master.NewHandlerContainer(osMux)
	version := apiserver.APIGroupVersion{
		Root:    "/oapi",
		Version: latest.Version,

		Mapper: latest.RESTMapper,

		Storage: storage,
		Codec:   interfaces.Codec,

		Creater:   kapi.Scheme,
		Typer:     kapi.Scheme,
		Convertor: kapi.Scheme,
		Linker:    interfaces.MetadataAccessor,

		Admit:   admit.NewAlwaysAdmit(),
		Context: kapi.NewRequestContextMapper(),
	}
	if err := version.InstallREST(handlerContainer); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	walkJSONFiles("fixtures", func(name, path string, data []byte) {
		template, err := interfaces.Codec.Decode(data)
		if err != nil {
			t.Errorf("%q: unexpected error: %v", path, err)
			return
		}
		config, err := osClient.TemplateConfigs("default").Create(template.(*templateapi.Template))
		if err != nil {
			t.Errorf("%q: unexpected error: %v", path, err)
			return
		}
		if len(config.Objects) == 0 {
			t.Errorf("%q: no items in config object", path)
			return
		}
		t.Logf("tested %q", path)
	})
}
