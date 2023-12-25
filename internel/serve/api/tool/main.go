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
	if err := c.SetUp(); err != nil {
		logx.Error(err)
		return
	}
	db.InitSqlite(c.DB)
	taskDB := model.NewTaskModel(db.GetDB())
	taskList, err := parseTaskList(taskDB)
	if err != nil {
		logx.Error(err)
		return
	}

	task_control.InitTaskConfig(c)
	core := task_control.NewTaskControl(c.TaskControlConfig.Concurrency)
	core.Run(taskList)
}

func parseTaskList(taskDB *model.TaskModel) ([]model.Task, error) {
	taskList := make([]model.Task, 0)
	bs, err := os.ReadFile(*taskFile)
	if err != nil {
		return nil, err
	}

	content := string(bs)
	if content == "" {
		return nil, errors.New("content is 0")
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
			return nil, err
		}
		index := strings.LastIndex(u.Path, ".")
		ext := u.Path[index+1:]
		switch ext {
		case model.VideoTypeMp4:
		case model.VideoTypeM3u8:
		default:
			return nil, errors.New("ext error")
		}

		t := model.Task{
			Name:      key,
			VideoType: ext,
			Type:      "url", // todo curl
			Data:      value,
		}
		if err := taskDB.Insert(&t); err != nil {
			return nil, err
		}
		taskList = append(taskList, t)

		i++
	}

	return taskList, nil
}
