package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/commitshark/notification-svc/internal/domain/ports"
	domain_template "github.com/commitshark/notification-svc/internal/domain/templates"
)

type GoTemplateRenderer struct {
	templates *template.Template
}

type EmailLayoutData struct {
	Subject        string
	Preheader      string
	UnsubscribeURL string
	Body           template.HTML
}

func NewGoTemplateRenderer(fsys fs.FS) (ports.TemplateRenderer, error) {
	tmpl, err := template.ParseFS(fsys, "emails/*.html")
	if err != nil {
		return nil, err
	}

	// Log all parsed templates
	for _, t := range tmpl.Templates() {
		fmt.Println("Parsed template:", t.Name())
	}

	return &GoTemplateRenderer{
		templates: tmpl,
	}, nil
}

func (r *GoTemplateRenderer) Render(
	templateName, subject string,
	data any,
	preHeader *string,
) (string, error) {
	templateNameFull := fmt.Sprintf("%s.html", templateName)

	d, ok := data.(domain_template.EmailTemplateData)
	if !ok {
		return "", fmt.Errorf("invalid template data type")
	}

	if r.templates.Lookup(templateNameFull) == nil {
		return "", fmt.Errorf("template %s not found", templateNameFull)
	}

	var buf bytes.Buffer
	err := r.templates.ExecuteTemplate(&buf, templateNameFull, d)
	if err != nil {
		return "", err
	}

	preHeaderStr := ""
	if preHeader != nil {
		preHeaderStr = *preHeader
	}

	layoutData := EmailLayoutData{
		Subject:        subject,
		Preheader:      preHeaderStr,
		UnsubscribeURL: "https://eventor.com/unsubscribe",
		Body:           template.HTML(buf.String()),
	}

	var layoutBuf bytes.Buffer
	err = r.templates.ExecuteTemplate(&layoutBuf, "layout.html", layoutData)
	if err != nil {
		return "", err
	}

	return layoutBuf.String(), nil
}
