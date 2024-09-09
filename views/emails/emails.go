package emails

import (
	"context"
	"embed"
	"io"

	"github.com/a-h/templ"
)

//go:embed *.txt
var textTemplates embed.FS

type MailTemplateHandler interface {
	GenerateTextVersion() (string, error)
	GenerateHtmlVersion() (string, error)
	Render(ctx context.Context, w io.Writer) error
}

func unsafe(html string) templ.Component {
	return templ.ComponentFunc(
		func(ctx context.Context, w io.Writer) (err error) {
			_, err = io.WriteString(w, html)
			return
		},
	)
}
