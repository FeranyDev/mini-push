package api

import (
	"fmt"
	"github.com/gomarkdown/markdown/ast"
	"io"
	"net/http"
	"os"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
)

func markdownToHtml(c echo.Context) error {
	url := c.Param("url")
	parameters := c.QueryParams()
	if len(parameters) > 0 {
		url = url + "?" + parameters.Encode()
	}

	return c.HTML(http.StatusOK, markToHtml(url))
}

func markToHtml(url string) string {

	// Get Markdown
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error getting file:", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return ""
	}

	md := markdown.NormalizeNewlines(body)

	// create Markdown parser
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Autolink
	p := parser.NewWithExtensions(extensions)

	// parse markdown into AST tree
	doc := p.Parse(md)

	// see AST tree
	//fmt.Printf("%s", "--- AST tree:\n")
	ast.Print(os.Stdout, doc)

	// create HTML renderer
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.CompletePage
	opts := html.RendererOptions{
		Flags: htmlFlags,
		CSS:   "/markdown-to-html/style.css",
	}
	renderer := html.NewRenderer(opts)

	htmlText := markdown.Render(doc, renderer)

	return string(htmlText)
}
