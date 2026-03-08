package proxy

type VideoTask struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    []byte
	Title   string
}
