package webhelp

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/yosssi/gohtml"
)

// RenderHTML writes a templ component to the ResponseWriter.
// In dev mode, it formats the HTML output for readability.
func RenderHTML(ctx context.Context, w http.ResponseWriter, component templ.Component) error {
	w.Header().Set("Content-Type", "text/html")
	if DevMode() {
		var buf bytes.Buffer
		if err := component.Render(ctx, &buf); err != nil {
			return err
		}
		formatted := gohtml.Format(buf.String())
		_, err := io.WriteString(w, formatted)
		return err
	}
	return component.Render(ctx, w)
}
