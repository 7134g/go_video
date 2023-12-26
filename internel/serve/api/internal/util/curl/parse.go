package curl

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

func Parse(content string) (_url string, header http.Header, err error) {
	content = strings.TrimSpace(content)
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		err = errors.New("curl parse content error: " + content)
		return
	}

	regData, _ := regexp.Compile(`'(.*?)'`)
	findResult := regData.FindStringSubmatch(lines[0])
	if len(findResult) < 2 {
		err = errors.New("curl parse url error: " + lines[0])
		return
	}
	_url = findResult[1]
	header = http.Header{}
	for _, s := range lines[1:] {
		data := regData.FindStringSubmatch(s)
		if len(data) != 2 {
			continue
		}
		part := strings.SplitN(data[1], ":", 2)
		if len(part) != 2 {
			continue
		}
		key := part[0]
		value := strings.TrimSpace(part[1])
		header.Set(key, value)
	}

	return
}
