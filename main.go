package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

var server *http.Server

//go:embed static/github.css
var githubCSS string

func main() {
	server = &http.Server{Addr: ":6969"}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedPath := r.URL.Path

		if requestedPath == "/" {
			filePath := "./index.md"
			serveMarkdownFile(filePath, w, r)
			return
		}

		requestedPath = strings.TrimPrefix(requestedPath, "/")
		requestedPath = strings.TrimSuffix(requestedPath, "/")

		fileInfo, err := os.Stat(requestedPath)
		if err == nil && fileInfo.IsDir() {
			filePath := filepath.Join(requestedPath, "index.md")
			serveMarkdownFile(filePath, w, r)
			return
		}

		requestedPath += ".md"

		filePath := filepath.Join(".", requestedPath)
		serveMarkdownFile(filePath, w, r)
	})

	log.Println("Starting server on http://localhost:6969")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", ":6969", err)
	}
}

func serveMarkdownFile(filePath string, w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
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
	opts := html.RendererOptions{
		Flags: html.CommonFlags,
	}

	renderer := html.NewRenderer(opts)

	htmlContent := markdown.ToHTML(mdContent, nil, renderer)

	htmlWithCSS := fmt.Sprintf("<style>%s</style>\n%s", githubCSS, htmlContent)

	return []byte(htmlWithCSS)
}
