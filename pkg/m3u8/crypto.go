package m3u8

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// Decrypt 使用 AES-128-CBC 解密数据
func Decrypt(data, key, iv []byte) ([]byte, error) {
	if len(key) != 16 {
		return nil, errors.New("key must be 16 bytes for AES-128")
	}
	if len(iv) != 16 {
		return nil, errors.New("iv must be 16 bytes")
	}
	if len(data) == 0 || len(data)%aes.BlockSize != 0 {
		return nil, errors.New("invalid data length")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(data))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, data)

	return pkcs7Unpad(decrypted)
}

// pkcs7Unpad 移除 PKCS7 填充
func pkcs7Unpad(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	padding := int(data[len(data)-1])
	if padding > aes.BlockSize || padding == 0 {
		return nil, errors.New("invalid padding")
	}
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding")
		}
	}
	return data[:len(data)-padding], nil
}

// ParseIV 解析 IV 字符串，支持 0x 前缀的十六进制格式
func ParseIV(ivStr string, segmentIndex uint64) ([]byte, error) {
	if ivStr == "" {
		// 默认使用 segment 序号作为 IV
		iv := make([]byte, 16)
		for i := 15; i >= 8; i-- {
			iv[i] = byte(segmentIndex)
			segmentIndex >>= 8
		}
		return iv, nil
	}

	// 移除 0x 前缀
	if len(ivStr) >= 2 && ivStr[:2] == "0x" {
		ivStr = ivStr[2:]
	}

	iv, err := hex.DecodeString(ivStr)
	if err != nil {
		return nil, fmt.Errorf("invalid IV hex: %w", err)
	}
	if len(iv) != 16 {
		return nil, errors.New("IV must be 16 bytes")
	}
	return iv, nil
}

// KeyCache 密钥缓存
type KeyCache struct {
	mu    sync.RWMutex
	cache map[string][]byte
}

// NewKeyCache 创建密钥缓存
func NewKeyCache() *KeyCache {
	return &KeyCache{cache: make(map[string][]byte)}
}

// Get 获取缓存的密钥
func (c *KeyCache) Get(uri string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	key, ok := c.cache[uri]
	return key, ok
}

// Set 缓存密钥
func (c *KeyCache) Set(uri string, key []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[uri] = key
}

// DownloadKey 下载密钥文件
func DownloadKey(keyURL string, baseURL *url.URL, header http.Header) ([]byte, error) {
	fullURL := ResolveURL(keyURL, baseURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download key failed: %s", resp.Status)
	}

	key, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(key) != 16 {
		return nil, fmt.Errorf("invalid key length: %d", len(key))
	}
	return key, nil
}

// ResolveURL 将相对 URL 转换为绝对 URL
func ResolveURL(uri string, baseURL *url.URL) string {
	if uri == "" {
		return ""
	}
	parsed, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	if parsed.IsAbs() {
		return uri
	}
	return baseURL.ResolveReference(parsed).String()
}
