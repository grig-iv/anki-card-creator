package ld

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func hasClass(node *html.Node, classes ...string) bool {
	for _, attr := range node.Attr {
		if attr.Key != "class" {
			continue
		}

		for _, class := range classes {
			if strings.Contains(attr.Val, class) == false {
				return false
			}

		}

		return true
	}

	return false
}

func innerTextTrim(node *html.Node) string {
	return strings.TrimSpace(htmlquery.InnerText(node))
}
