package m3u8

import (
	"os"
	"testing"
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
}
