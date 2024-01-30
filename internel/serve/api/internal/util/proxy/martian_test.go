package proxy

import (
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/util/model"
	"testing"
)

func TestMartian(t *testing.T) {
	db.InitSqlite("test.sqlite")
	taskDB = model.NewTaskModel(db.GetDB())
	SetServeProxyAddress("http://127.0.0.1:7890", "", "")
	OpenCert()
	if err := Martian(); err != nil {
		t.Fatal(err)
	}

}
