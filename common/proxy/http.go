package proxy

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
)

func ParseHtmlTitle(r io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return "", err
	}

	node := doc.Find("title").First()
	if node == nil {
		return "", errors.New("cannot find title")
	}

	return node.Text(), nil
}

func ExtractFilename(urlPath string) string {
	return regUrl.FindString(urlPath)
}
