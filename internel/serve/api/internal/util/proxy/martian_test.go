package proxy

import (
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/util/model"
	"testing"
)

func TestMartian(t *testing.T) {
	db.InitSqlite("test.sqlite")
	taskDB := model.NewTaskModel(db.GetDB())
	err := Martian(taskDB)
	if err != nil {
		t.Fatal(err)
	}

}
