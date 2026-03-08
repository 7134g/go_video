package proxy

import (
	"bytes"
	"fmt"
	"net/url"
	"regexp"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseHtml(t *testing.T) {
	data := bytes.NewBuffer([]byte(`
<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Go 语言环境安装 | 菜鸟教程</title>
		<title>不要的</title>
		<title>多余的</title>
	</head>
</html>
`))
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		t.Fatal(err)
	}
	node := doc.Find("title").First().Text()
	t.Log(node)
}

func TestExtractFilename(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/video/test.m3u8", "test.m3u8"},
		{"https://example.com/path/to/video.mp4", "video.mp4"},
		{"https://example.com/v/f230.m3u8", "f230.m3u8"},
	}

	for _, tt := range tests {
		u, err := url.Parse(tt.url)
		if err != nil {
			t.Fatal(err)
		}
		result := ExtractFilename(u.Path)
		if result != tt.expected {
			t.Errorf("ExtractFilename(%s) = %s, want %s", tt.url, result, tt.expected)
		}
	}
}

func TestRegexMatch(t *testing.T) {
	u, err := url.Parse("https://example.com/v/f230.m3u8")
	if err != nil {
		t.Fatal(err)
	}

	reg := regexp.MustCompile(`([^/]+)(\.m3u8|\.mp4)$`)
	result := reg.FindString(u.Path)
	t.Log(result)

	html := `<html><body>f230.m3u8</body></html>`
	regKey := regexp.MustCompile(fmt.Sprintf("(%s)|(%s)|(%s)",
		"https://example.com/v/f230.m3u8", "f230.m3u8", "f230"))
	t.Log(regKey.MatchString(html))
}
