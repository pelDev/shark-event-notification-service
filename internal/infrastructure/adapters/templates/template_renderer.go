package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/commitshark/notification-svc/internal/domain/ports"
)

type GoTemplateRenderer struct {
	templates *template.Template
}

type EmailLayoutData struct {
	Subject        string        // {{ subject }}
	Preheader      string        // {{ preheader }}
	UnsubscribeURL string        // {{ unsubscribeUrl }}
	Body           template.HTML // {{ body }} (already-rendered HTML partial)
}

func NewGoTemplateRenderer(fsys fs.FS) (ports.TemplateRenderer, error) {
	tmpl, err := template.ParseFS(fsys, "*.html")
	if err != nil {
		return nil, err
	}

	return &GoTemplateRenderer{
		templates: tmpl,
	}, nil
}

func (r *GoTemplateRenderer) Render(
	templateName, subject string,
	data any,
) (string, error) {
	d, ok := data.(EmailTemplateData)
	if !ok {
		return "", fmt.Errorf("invalid template data type")
	}

	if r.templates.Lookup(templateName) == nil {
		return "", fmt.Errorf("template %s not found", templateName)
	}

	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, templateName, d)
	if err != nil {
		return "", err
	}

	layoutData := EmailLayoutData{
		Subject:        subject,
		Preheader:      "",
		UnsubscribeURL: "https://eventor.com/unsubscribe",
		Body:           template.HTML(buf.String()),
	}

	var layoutBuf bytes.Buffer
	err = r.templates.ExecuteTemplate(&layoutBuf, "layout", layoutData)
	if err != nil {
		return "", err
	}

	return layoutBuf.String(), nil
}
