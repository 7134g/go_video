package proxy

import (
	"fmt"
	"testing"
)

func TestGetVideo(t *testing.T) {
	fmt.Println(GetVideo("https://cableav.video/index/data/dplayer.html?video=https://t0.97img.com/a1002239/a.m3u8"))
	fmt.Println(GetVideo("https://cableav.video/index/data/dplayer.html?z=https://t0.97img.com/a1002239/a.m3u8?r=https://t0.97img.com/a10sss2239/a.m3u8"))
	fmt.Println(GetVideo("https://cableav.video/index/data/dplayer.html?x=https://t0.97img.com/a1002239/a.m3u8?y=https://t0.97img.com/a10sss2239/a.m3u8"))
	fmt.Println(GetVideo("https://t0.97img.com/a10sss2239/a.m3u8"))
	fmt.Println(GetVideo("https://t0.97img.com/a10sss2239/a.mp4"))
}
