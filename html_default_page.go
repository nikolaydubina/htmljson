package htmljson

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed html_default_page.html
var defaultPageTemplate []byte

var DefaultPageMarshaler = PageMarshaler{
	Title:            "htmljson",
	Template:         defaultPageTemplate,
	TemplateTitleKey: `{{.Title}}`,
	TemplateJSONKey:  `{{.HTMLJSON}}`,
	Marshaler:        &DefaultMarshaler,
}

// PageMarshaler encodes JSON via marshaller into HTML page by placing Title and content appropriately.
type PageMarshaler struct {
	Title            string
	Template         []byte
	TemplateTitleKey string
	TemplateJSONKey  string

	Marshaler interface {
		MarshalTo(w io.Writer, v any) error
	}

	idxTitle    int
	idxHTMLJSON int
}

func (m *PageMarshaler) Marshal(v any) []byte {
	b := bytes.Buffer{}
	m.MarshalTo(&b, v)
	return b.Bytes()
}

func (m *PageMarshaler) parseTemplate() {
	if m.idxTitle == 0 || m.idxHTMLJSON == 0 {
		m.idxTitle = bytes.Index(m.Template, []byte(m.TemplateTitleKey))
		m.idxHTMLJSON = bytes.Index(m.Template, []byte(m.TemplateJSONKey))
	}
}

func (m *PageMarshaler) MarshalTo(w io.Writer, v any) error {
	m.parseTemplate()

	var s int

	if f := m.idxTitle; f > 0 {
		w.Write(m.Template[s:f])
		s = f + len(m.TemplateTitleKey)
		w.Write([]byte(m.Title))
	}

	if f := m.idxHTMLJSON; f > 0 {
		w.Write(m.Template[s:f])
		s = f + len(m.TemplateJSONKey)
		m.Marshaler.MarshalTo(w, v)
	}

	w.Write(m.Template[s:])

	return nil
}
