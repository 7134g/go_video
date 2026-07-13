//go:build bindings

package main

import (
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

func main() {
	initShared()

	app := NewApp()

	err := wails.Run(&options.App{
		Title: "视频下载器",
		AssetServer: &assetserver.Options{
			Assets: nil,
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
