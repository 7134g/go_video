package m3u8

import (
	"testing"
)

func TestDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	if len(key) != 16 {
		t.Error("key should be 16 bytes")
	}
}

func TestParseIV(t *testing.T) {
	iv, err := ParseIV("", 1)
	if err != nil {
		t.Fatalf("parse empty IV failed: %v", err)
	}
	if len(iv) != 16 {
		t.Errorf("expected 16 bytes IV, got %d", len(iv))
	}

	iv, err = ParseIV("0x00000000000000000000000000000001", 0)
	if err != nil {
		t.Fatalf("parse hex IV failed: %v", err)
	}
	if iv[15] != 1 {
		t.Errorf("expected last byte to be 1, got %d", iv[15])
	}
}
