package task_control

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/util/model"
	"github.com/zeromicro/go-zero/core/logx"
	"testing"
)

func TestNewTaskControl(t *testing.T) {
	c := NewTaskControl(3)
	db.InitSqlite("test.sqlite")
	cfg := config.Config{
		HttpConfig: config.HttpConfig{
			Headers: map[string]string{
				"user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			},
			Proxy:       "",
			ProxyStatus: false,
		},
		TaskControlConfig: config.TaskControlConfig{
			Concurrency:       2,
			ConcurrencyM3u8:   10,
			SaveDir:           "./download",
			TaskErrorMaxCount: 20,
			TaskErrorDuration: 1,
			UseFfmpeg:         false,
			FfmpegPath:        "",
		},
	}
	cfg.Log.Encoding = "plain"
	cfg.Log.Mode = "console"
	cfg.Log.Level = "debug"
	_ = logx.SetUp(cfg.Log)
	InitTask(cfg)

	tasks := []model.Task{
		//{
		//	ID:        1,
		//	Name:      "测试mp4_1",
		//	VideoType: "mp4",
		//	Type:      "url",
		//	Data:      "http://clips.vorwaerts-gmbh.de/big_buck_bunny.mp4",
		//	Status:    0,
		//},
		{
			ID:        2,
			Name:      "测试m3u8_1",
			VideoType: "m3u8",
			Type:      "url",
			Data:      "https://1257120875.vod2.myqcloud.com/0ef121cdvodtransgzp1257120875/3055695e5285890780828799271/v.f230.m3u8",
			Status:    0,
		},
	}
	c.Run(tasks)
}
