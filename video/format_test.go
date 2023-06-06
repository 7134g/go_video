package video

import "testing"

func TestCompleteURL(t *testing.T) {
	t.Log(CompleteURL("http://www.baidu.com"))
	t.Log(CompleteURL("https://www.baidu.com"))
	t.Log(CompleteURL("www.baidu.com"))

}
