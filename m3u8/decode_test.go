package m3u8

import (
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	f, err := os.Open("index.m3u8")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	m, err := parse(f)
	if err != nil {
		t.Error(err)
	}

	t.Log(m)

	f, err = os.Open("test.m3u8")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()

	m1, err := parse(f)
	if err != nil {
		t.Error(err)
	}

	t.Log(m1)
}

func TestCalculationTime(t *testing.T) {
	t.Log(CalculationTime(36640.25))
	b := time.Now()
	t.Log(fmt.Sprintf("%s", time.Now().Add(time.Second*100).Sub(b)))
}

func TestName(t *testing.T) {
	u1 := `https://s7.fsvod1.com/20220717/WZ2rIvnp/index1.m3u8`
	u2 := `https://s7.fsvod1.com/hls/16/20221004/264860/plist-00001.ts`

	p1, _ := url.Parse(u1)
	p3, _ := p1.Parse(u2)
	fmt.Println(p3.String())
}
