package origin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful"
)

const asset_off_notice = "You need to upgrade to OpenShift in order to take advantage of this feature"

func TestInitializeOpenshiftAPIVersionRouteHandler(t *testing.T) {
	service := new(restful.WebService)
	initAPIVersionRoute(service, "osapi", "v1beta3")

	if len(service.Routes()) != 1 {
		t.Fatalf("Exp. the OSAPI route but found none")
	}
	route := service.Routes()[0]
	if !contains(route.Produces, restful.MIME_JSON) {
		t.Fatalf("Exp. route to produce mimetype json")
	}
	if !contains(route.Consumes, restful.MIME_JSON) {
		t.Fatalf("Exp. route to consume mimetype json")
	}
}

// assetServerOffNotice Case 1: Request / and expect assetServerOffNotice to not handle the request
func Test_assetServerOffNoticeDoesntHandleRoot(t *testing.T) {
	handler, recorder := setUpAssetServer()
	req, _ := http.NewRequest("GET", "/", nil)
	handler.ServeHTTP(recorder, req)
	if recorder.Code != 404 {
		t.Fatalf("assetServerOffNotice returned unexpected response: %i != 404", recorder.Code)
	}
}

// assetServerOffNotice Case 2: Request /login and /logout an off notice
func Test_assetServerOffNoticeHandlesLoginLogoutConsole(t *testing.T) {
	paths := []string{"/login", "/logout", "/console", "/console/java"}
	for i := range paths {
		handler, recorder := setUpAssetServer()
		req, _ := http.NewRequest("GET", paths[i], nil)
		handler.ServeHTTP(recorder, req)
		if recorder.Code != 200 {
			t.Fatalf("assetServerOffNotice returned unexpected response: %i != 200", recorder.Code)
		}
		body := recorder.Body.String()
		if body != asset_off_notice {
			t.Fatalf("assetServerOffNotice did not return the proper body. %s", body)
		}
	}
}

func contains(list []string, value string) bool {
	for _, entry := range list {
		if entry == value {
			return true
		}
	}
	return false
}

// Setup for the assetServer tests
func setUpAssetServer() (http.Handler, *httptest.ResponseRecorder) {
	handler := assetServerOffNotice(http.NotFoundHandler())
	recorder := httptest.NewRecorder()
	return handler, recorder
}
