package types

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gen/field"
)

type PageTurner interface {
	GetPager() (page, size int)
	GetWhere() map[string]interface{}
	ParseMongo() (map[string]interface{}, *options.FindOptions)
	ParseMysql() (where map[string]interface{}, offset, limit int)
	ParseMysqlOrderBy(tableName string) field.Expr
}

const (
	Desc = "desc"
	Asc  = "asc"
)

func (d *DbQueryList) GetPager() (page int, size int) {
	page = d.Page
	if d.Page <= 0 {
		page = 1
	}

	size = d.Size
	if d.Size <= 0 {
		size = 1
	}

	return page, size
}

func (d *DbQueryList) GetWhere() map[string]interface{} {
	var err error
	where := map[string]interface{}{}
	for k, v := range d.Where {
		switch v.(type) {
		case json.Number:
			where[k], err = v.(json.Number).Int64()
			if err != nil {
				where[k], _ = v.(json.Number).Float64()
			}
		case string:
			valueStr := v.(string)
			if valueStr == "" {
				continue
			}
			where[k] = valueStr
		default:
			where[k] = v
		}
	}
	return where
}

func (d *DbQueryList) orderMongo() bson.D {
	order := bson.D{}
	if d.Order != "" && d.OrderKey != "" {
		seq := 0
		key := d.OrderKey
		switch d.Order {
		case Desc:
			seq = -1
		case Asc:
			seq = 1
		default:
			seq = 1
		}

		if key == "id" {
			key = "_id"
		}

		order = append(order, bson.E{
			Key: key, Value: seq,
		})
	}

	return order
}

func (d *DbQueryList) ParseMongo() (map[string]interface{}, *options.FindOptions) {
	page, size := d.GetPager()
	offset := int64((page - 1) * size)
	limit := int64(size)

	opt := &options.FindOptions{}
	opt.Skip = &offset
	opt.Limit = &limit
	opt.Sort = d.orderMongo()

	where := d.GetWhere()

	return where, opt
}

func (d *DbQueryList) ParseMysql() (where map[string]interface{}, offset, limit int) {
	page, size := d.GetPager()

	where = d.GetWhere()
	offset = (page - 1) * size
	limit = size

	return where, offset, limit
}

func (d *DbQueryList) ParseMysqlOrderBy(tableName string) field.Expr {
	var orderExpr field.Expr
	key := field.NewField(tableName, d.OrderKey)
	if d.Order == Desc {
		orderExpr = key.Desc()
	} else {
		orderExpr = key
	}

	return orderExpr
}
