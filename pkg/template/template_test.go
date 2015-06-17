package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
	"testing"

	_ "github.com/GoogleCloudPlatform/kubernetes/pkg/api/latest"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util"

	"github.com/projectatomic/appinfra-next/pkg/api/latest"
	"github.com/projectatomic/appinfra-next/pkg/api/v1beta3"
	"github.com/projectatomic/appinfra-next/pkg/template/api"
	"github.com/projectatomic/appinfra-next/pkg/template/generator"
)

func makeParameter(name, value, generate string) api.Parameter {
	return api.Parameter{
		Name:     name,
		Value:    value,
		Generate: generate,
	}
}

func TestAddParameter(t *testing.T) {
	var template api.Template

	jsonData, _ := ioutil.ReadFile("../../test/templates/fixtures/guestbook.json")
	json.Unmarshal(jsonData, &template)

	AddParameter(&template, makeParameter("CUSTOM_PARAM", "1", ""))
	AddParameter(&template, makeParameter("CUSTOM_PARAM", "2", ""))

	if p := GetParameterByName(&template, "CUSTOM_PARAM"); p == nil {
		t.Errorf("Unable to add a custom parameter to the template")
	} else {
		if p.Value != "2" {
			t.Errorf("Unable to replace the custom parameter value in template")
		}
	}
}

type FooGenerator struct {
}

func (g FooGenerator) GenerateValue(expression string) (interface{}, error) {
	return "foo", nil
}

type ErrorGenerator struct {
}

func (g ErrorGenerator) GenerateValue(expression string) (interface{}, error) {
	return "", fmt.Errorf("error")
}

func TestParameterGenerators(t *testing.T) {
	tests := []struct {
		parameter  api.Parameter
		generators map[string]generator.Generator
		shouldPass bool
		expected   api.Parameter
	}{
		{ // Empty generator, should pass
			makeParameter("PARAM", "X", ""),
			map[string]generator.Generator{},
			true,
			makeParameter("PARAM", "X", ""),
		},
		{ // Foo generator, should pass
			makeParameter("PARAM", "", "foo"),
			map[string]generator.Generator{"foo": FooGenerator{}},
			true,
			makeParameter("PARAM", "foo", ""),
		},
		{ // Invalid generator, should fail
			makeParameter("PARAM", "", "invalid"),
			map[string]generator.Generator{"invalid": nil},
			false,
			makeParameter("PARAM", "", "invalid"),
		},
		{ // Error generator, should fail
			makeParameter("PARAM", "", "error"),
			map[string]generator.Generator{"error": ErrorGenerator{}},
			false,
			makeParameter("PARAM", "", "error"),
		},
	}

	for i, test := range tests {
		processor := NewProcessor(test.generators)
		template := api.Template{Parameters: []api.Parameter{test.parameter}}
		err := processor.GenerateParameterValues(&template)
		if err != nil && test.shouldPass {
			t.Errorf("test[%v]: Unexpected error %v", i, err)
		}
		if err == nil && !test.shouldPass {
			t.Errorf("test[%v]: Expected error", i)
		}
		actual := template.Parameters[0]
		if actual.Value != test.expected.Value {
			t.Errorf("test[%v]: Unexpected value: Expected: %#v, got: %#v", i, test.expected.Value, test.parameter.Value)
		}
	}
}

func TestProcessValueEscape(t *testing.T) {
	var template api.Template
	if err := latest.Codec.DecodeInto([]byte(`{
		"kind":"Template", "apiVersion":"v1",
		"objects": [
			{
				"kind": "Service", "apiVersion": "v1beta3${VALUE}",
				"metadata": {
					"labels": {
						"key1": "${VALUE}",
						"key2": "$${VALUE}"
					}
				}
			}
		]
	}`), &template); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	generators := map[string]generator.Generator{
		"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(1337))),
	}
	processor := NewProcessor(generators)

	// Define custom parameter for the transformation:
	AddParameter(&template, makeParameter("VALUE", "1", ""))

	// Transform the template config into the result config
	errs := processor.Process(&template)
	if len(errs) > 0 {
		t.Fatalf("unexpected error: %v", errs)
	}
	result, err := v1beta3.Codec.Encode(&template)
	if err != nil {
		t.Fatalf("unexpected error during encoding Config: %#v", err)
	}
	expect := `{"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},"objects":[{"apiVersion":"v1beta31","kind":"Service","metadata":{"labels":{"key1":"1","key2":"$1"}}}],"parameters":[{"name":"VALUE","value":"1"}]}`
	if expect != string(result) {
		t.Errorf("unexpected output: %s", util.StringDiff(expect, string(result)))
	}
}

var trailingWhitespace = regexp.MustCompile(`\n\s*`)

func TestEvaluateLabels(t *testing.T) {
	testCases := map[string]struct {
		Input  string
		Output string
		Labels map[string]string
	}{
		"no labels": {
			Input: `{
				"kind":"Template", "apiVersion":"v1",
				"objects": [
					{
						"kind": "Service", "apiVersion": "v1beta3",
						"metadata": {"labels": {"key1": "v1", "key2": "v2"}	}
					}
				]
			}`,
			Output: `{
				"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},
				"objects":[
					{
						"apiVersion":"v1beta3","kind":"Service","metadata":{
						"labels":{"key1":"v1","key2":"v2"}}
					}
				]
			}`,
		},
		"one different label": {
			Input: `{
				"kind":"Template", "apiVersion":"v1",
				"objects": [
					{
						"kind": "Service", "apiVersion": "v1beta3",
						"metadata": {"labels": {"key1": "v1", "key2": "v2"}	}
					}
				]
			}`,
			Output: `{
				"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},
				"objects":[
					{
						"apiVersion":"v1beta3","kind":"Service","metadata":{
						"labels":{"key1":"v1","key2":"v2","key3":"v3"}}
					}
				],
				"labels":{"key3":"v3"}
			}`,
			Labels: map[string]string{"key3": "v3"},
		},
		"when the root object has labels and no metadata": {
			Input: `{
				"kind":"Template", "apiVersion":"v1",
				"objects": [
					{
						"kind": "Service", "apiVersion": "v1beta1",
						"labels": {
							"key1": "v1",
							"key2": "v2"
						}
					}
				]
			}`,
			Output: `{
				"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},
				"objects":[
					{
						"apiVersion":"v1beta1","kind":"Service",
						"labels":{"key1":"v1","key2":"v2","key3":"v3"}
					}
				],
				"labels":{"key3":"v3"}
			}`,
			Labels: map[string]string{"key3": "v3"},
		},
		"when the root object has labels and metadata": {
			Input: `{
				"kind":"Template", "apiVersion":"v1",
				"objects": [
					{
						"kind": "Service", "apiVersion": "v1beta1",
						"metadata": {},
						"labels": {
							"key1": "v1",
							"key2": "v2"
						}
					}
				]
			}`,
			Output: `{
				"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},
				"objects":[
					{
						"apiVersion":"v1beta1","kind":"Service",
						"labels":{"key1":"v1","key2":"v2"},
						"metadata":{"labels":{"key3":"v3"}}
					}
				],
				"labels":{"key3":"v3"}
			}`,
			Labels: map[string]string{"key3": "v3"},
		},
		"overwrites label": {
			Input: `{
				"kind":"Template", "apiVersion":"v1",
				"objects": [
					{
						"kind": "Service", "apiVersion": "v1beta3",
						"metadata": {"labels": {"key1": "v1", "key2": "v2"}	}
					}
				]
			}`,
			Output: `{
				"kind":"Template","apiVersion":"v1beta3","metadata":{"creationTimestamp":null},
				"objects":[
					{
						"apiVersion":"v1beta3","kind":"Service","metadata":{
						"labels":{"key1":"v1","key2":"v3"}}
					}
				],
				"labels":{"key2":"v3"}
			}`,
			Labels: map[string]string{"key2": "v3"},
		},
	}

	for k, testCase := range testCases {
		var template api.Template
		if err := latest.Codec.DecodeInto([]byte(testCase.Input), &template); err != nil {
			t.Errorf("%s: unexpected error: %v", k, err)
			continue
		}

		generators := map[string]generator.Generator{
			"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(1337))),
		}
		processor := NewProcessor(generators)

		template.ObjectLabels = testCase.Labels

		// Transform the template config into the result config
		errs := processor.Process(&template)
		if len(errs) > 0 {
			t.Errorf("%s: unexpected error: %v", k, errs)
			continue
		}
		result, err := v1beta3.Codec.Encode(&template)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", k, err)
			continue
		}
		expect := testCase.Output
		expect = trailingWhitespace.ReplaceAllString(expect, "")
		if expect != string(result) {
			t.Errorf("%s: unexpected output: %s", k, util.StringDiff(expect, string(result)))
			continue
		}
	}
}

func TestProcessTemplateParameters(t *testing.T) {
	var template, expectedTemplate api.Template
	jsonData, _ := ioutil.ReadFile("../../test/templates/fixtures/guestbook.json")
	if err := latest.Codec.DecodeInto(jsonData, &template); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedData, _ := ioutil.ReadFile("../../test/templates/fixtures/guestbook_list.json")
	if err := latest.Codec.DecodeInto(expectedData, &expectedTemplate); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	generators := map[string]generator.Generator{
		"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(1337))),
	}
	processor := NewProcessor(generators)

	// Define custom parameter for the transformation:
	AddParameter(&template, makeParameter("CUSTOM_PARAM1", "1", ""))

	// Transform the template config into the result config
	errs := processor.Process(&template)
	if len(errs) > 0 {
		t.Fatalf("unexpected error: %v", errs)
	}
	result, err := v1beta3.Codec.Encode(&template)
	if err != nil {
		t.Fatalf("unexpected error during encoding Config: %#v", err)
	}
	exp, _ := v1beta3.Codec.Encode(&expectedTemplate)

	if string(result) != string(exp) {
		t.Errorf("unexpected output: %s", util.StringDiff(string(exp), string(result)))
	}
}
