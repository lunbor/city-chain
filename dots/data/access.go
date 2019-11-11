package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/scryinfo/citychain/dots/data/dao"
	"github.com/scryinfo/dot/dot"
)

const (
	TypeId = "b9124432-a9a1-4b55-ba83-69a94d987a7f"
)

type config struct {
	//sample:  "user:password@/dbname?charset=utf8&parseTime=True&loc=Local"
	//see https://github.com/go-sql-driver/mysql#parameters
	DbParameters string `json:"dbParameters"`
}

type Access struct {
	conf config
	Db   *gorm.DB

	Block    dao.Block
	Purchase dao.Purchase
	Tx20     dao.Tx20
	Tx721    dao.Tx721
}

func (c *Access) Destroy(ignore bool) error {
	if c.Db != nil {
		c.Db.Close()
		c.Db = nil
	}

	return nil
}

func (c *Access) Create(l dot.Line) (err error) {
	c.Db, err = gorm.Open("mysql", c.conf.DbParameters)
	c.Db.LogMode(true) //todo for test
	if err != nil {
		c.Db = nil
	} else {
		c.Db.SingularTable(true) //不使用表名复数
	}
	return err
}

func newData(conf interface{}) (d dot.Dot, err error) {
	var bs []byte
	if bt, ok := conf.([]byte); ok {
		bs = bt
	} else {
		return nil, dot.SError.Parameter
	}
	dconf := &config{}
	err = dot.UnMarshalConfig(bs, dconf)
	if err != nil {
		return nil, err
	}

	d = &Access{conf: *dconf}

	return d, err
}

//TypeLives
func TypeLives() []*dot.TypeLives {
	lives := []*dot.TypeLives{
		{
			Meta: dot.Metadata{TypeId: TypeId, NewDoter: func(conf interface{}) (dot dot.Dot, err error) {
				return newData(conf)
			}},
		},
	}
	return lives
}

//ConfigTypeLives
func ConfigTypeLives() *dot.ConfigTypeLives {
	return &dot.ConfigTypeLives{
		TypeIdConfig: TypeId,
		ConfigInfo:   &config{},
	}
}
