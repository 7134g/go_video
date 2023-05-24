package video

type VideoSetting struct {
	VideoExt      string // 视频格式（用于文件尾缀）
	VideoCategory string

	CryptoKey    []byte // 密匙
	CryptoMethod string // 类型
}
