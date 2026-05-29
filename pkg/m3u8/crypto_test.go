package m3u8

import (
	"crypto/aes"
	"crypto/cipher"
	"testing"
)

func TestDecrypt_RoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef")
	iv := []byte("1234567890abcdef")
	plain := []byte("hello hls aes-128 cbc, padded by pkcs7.")

	// 加密（PKCS7 + AES-CBC）
	pad := aes.BlockSize - len(plain)%aes.BlockSize
	for i := 0; i < pad; i++ {
		plain = append(plain, byte(pad))
	}
	block, _ := aes.NewCipher(key)
	enc := make([]byte, len(plain))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(enc, plain)

	dec, err := Decrypt(enc, key, iv)
	if err != nil {
		t.Fatalf("decrypt: %v", err)
	}
	if string(dec) != "hello hls aes-128 cbc, padded by pkcs7." {
		t.Errorf("round-trip mismatch: %q", dec)
	}
}

func TestDecrypt_InvalidKeyLen(t *testing.T) {
	if _, err := Decrypt(make([]byte, 16), []byte("short"), make([]byte, 16)); err == nil {
		t.Error("expected error for short key")
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
