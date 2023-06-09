package htmljson_test

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/nikolaydubina/htmljson"
)

//go:embed testdata/example-page.html
var examplePageHTML string

//go:embed testdata/example-page-color.html
var examplePageColorHTML string

func TestMarshalHTML(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	h := htmljson.DefaultPageMarshaler.Marshal(v)

	os.WriteFile("testdata/example-page.out.html", h, 0666)
	if strings.TrimSpace(examplePageHTML) != strings.TrimSpace(string(h)) {
		t.Errorf("wrong output: %s", string(h))
	}
}

func TestMarshalHTML_Color(t *testing.T) {
	var v any
	json.Unmarshal(exampleJSON, &v)

	s := htmljson.Marshaler{
		Null:   htmljson.NullHTML,
		Bool:   htmljson.BoolHTML,
		String: htmljson.StringHTML,
		Number: func(k string, v float64, s string) string {
			if k == "$.cakes.strawberry-cake.size" {
				return `<div class="json-value json-number" style="color:red;">` + s + `</div>`
			}
			if v > 10 {
				return `<div class="json-value json-number" style="color:blue;">` + s + `</div>`
			}
			return `<div class="json-value json-number">` + s + `</div>`
		},
		Array: htmljson.DefaultArrayHTML,
		Map:   htmljson.DefaultMapHTML,
		Row:   htmljson.DefaultRowHTML{Padding: 4}.Marshal,
	}

	m := htmljson.DefaultPageMarshaler
	m.Marshaler = &s

	htmlPage := m.Marshal(v)

	os.WriteFile("testdata/example-page-color.out.html", []byte(htmlPage), 0666)
	if strings.TrimSpace(examplePageColorHTML) != strings.TrimSpace(string(htmlPage)) {
		t.Errorf("wrong output: %s", string(htmlPage))
	}
}
