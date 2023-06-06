package video

import (
	"dv/base"
	"testing"
)

func TestCompleteURL(t *testing.T) {
	t.Log(base.CompleteURL("http://www.baidu.com"))
	t.Log(base.CompleteURL("https://www.baidu.com"))
	t.Log(base.CompleteURL("www.baidu.com"))

}
