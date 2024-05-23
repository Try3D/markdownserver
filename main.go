package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed static/github.css
var githubCSS string

func main() {
	server := &http.Server{Addr: ":6969"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveMarkdownFile("."+r.URL.Path, w, r)
	})

	log.Println("Starting server on http://localhost:6969")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", ":6969", err)
	}
}

func serveMarkdownFile(filePath string, w http.ResponseWriter, r *http.Request) {
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) || fileInfo.IsDir() {
		http.NotFound(w, r)
		return
	}

	mdContent, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	htmlContent := convertMarkdownToHTML(mdContent)
	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

func convertMarkdownToHTML(mdContent []byte) []byte {
	var buf bytes.Buffer

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	if err := md.Convert(mdContent, &buf); err != nil {
		log.Fatalf("Error converting markdown to HTML: %v", err)
	}

	htmlWithCSS := fmt.Sprintf("<style>%s</style>\n%s", githubCSS, buf.String())
	return []byte(htmlWithCSS)
}

