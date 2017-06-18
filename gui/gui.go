// Host graphical front-end via http
package gui

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// only supports inputGuiBlocks for now
type GuiBlock interface {
	Put([]float64)
	Get() []float64
}

// Global collection of gui input fields
var GuiBlocks = make(map[string]GuiBlock)

var DefaultGuiHtml = "<!DOCTYPE html>\n" +
	"<html>\n" +
	"<head>\n" +
	"</head>\n" +
	"<body bgcolor=\"#E5E4E2\">\n" +
	"no gui elements specified" +
	"</body>\n" +
	"</html>"

type Gui struct {
	Html   string
	Server *http.Server
}

func generateFormHtml(fieldName string, vs []string) string {
	html := fieldName + ": "

	for _, v := range vs {
		html = html + "<input type = \"text\" name=\"" +
			fieldName + "\" value=\"" + v + "\">"
	}

	html = html + "<br>"

	return html
}

func getBlockValueStrings(b GuiBlock) []string {
	vs := b.Get()
	strs := []string{}

	for _, v := range vs {
		str := strconv.FormatFloat(v, 'g', 10, 64)
		strs = append(strs, str)
	}

	return strs
}

func makeGui() string {
	html := DefaultGuiHtml

	if len(GuiBlocks) > 0 {
		header := "<!DOCTYPE html>\n" +
			"<html>\n" +
			"<head>\n" +
			"</head>\n" +
			"<body bgcolor=\"#E5E4E2\">\n" +
			"<form action=\"/\" method=\"get\">\n"

		var body string

		footer := "<br><input type=\"submit\" value=\"Submit\">\n" +
			"</form>\n" +
			"</body>\n" +
			"</html>"

		for key, b := range GuiBlocks {
			strs := getBlockValueStrings(b)
			formHtml := generateFormHtml(key, strs)
			body = body + formHtml + "\n"
		}
		html = header + body + footer
	}

	return html
}

func (g *Gui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	isUpdated := false
	for key, values := range r.Form {
		fieldName := key

		if _, ok := GuiBlocks[fieldName]; ok {
			x := GuiBlocks[fieldName].Get()

			for i, v := range values {
				parsedFloat, err := strconv.ParseFloat(v, 64)

				if (err == nil) && (i < len(x)) {
					x[i] = parsedFloat
				}
			}

			GuiBlocks[fieldName].Put(x)
			isUpdated = true
		}

	}

	if isUpdated {
		g.Html = makeGui()
	}

	fmt.Fprintf(w, g.Html)
}

func LaunchGui() {
	guiHtml := makeGui()

	g := &Gui{
		Html: guiHtml,
	}

	g.Server = &http.Server{
		Addr:           ":8888", // TODO: specify by user
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go g.Server.ListenAndServe()
}
