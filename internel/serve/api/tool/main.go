package main

import (
	"dv/internel/serve/api/internal/config"
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/svc/task_control"
	"dv/internel/serve/api/internal/util/model"
	"errors"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	configFile = flag.String("f", "etc/task_serve.yaml", "the config file")
	taskFile   = flag.String("t", "url.txt", "默认：url.txt文件")
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	db.InitSqlite(c.DB)
	taskDB := model.NewTaskModel(db.GetDB())
	if err := parseTaskList(taskDB); err != nil {
		logx.Error(err)
	}

	task_control.InitTaskConfig(c)
	core := task_control.NewTaskControl(c.TaskControlConfig.Concurrency)
	taskList, err := taskDB.List()
	if err != nil {
		logx.Error(err)
	}
	core.Run(taskList)
}

func parseTaskList(taskDB *model.TaskModel) error {
	bs, err := os.ReadFile(*taskFile)
	if err != nil {
		return err
	}

	content := string(bs)
	if content == "" {
		return errors.New("content is 0")
	}
	reHead, _ := regexp.Compile(`\s+`)
	content = reHead.ReplaceAllString(content, "\n")
	content = strings.TrimPrefix(content, "\n")
	content = strings.TrimSuffix(content, "\n")
	list := strings.Split(content, "\n")
	for i := 0; i < len(list); i++ {
		if i+1 == len(list) {
			break
		}
		key := list[i]
		value := list[i+1]
		if len(value) < 4 || "http" != value[:4] {
			log.Println("错误值：", value)
			continue
		}
		u, err := url.Parse(value)
		if err != nil {
			return err
		}
		index := strings.LastIndex(u.Path, ".")
		ext := u.Path[index+1:]
		switch ext {
		case model.VideoTypeMp4:
		case model.VideoTypeM3u8:
		default:
			return errors.New("ext error")
		}

		if err := taskDB.Insert(&model.Task{
			Name:      key,
			VideoType: ext,
			Type:      "url", // todo curl
			Data:      value,
		}); err != nil {
			return err
		}

		i++
	}

	return nil
}
