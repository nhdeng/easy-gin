package easy

import (
	_ "github.com/go-sql-driver/mysql"
	Injector "github.com/shenyisyn/goft-ioc"
	"log"
	"xorm.io/xorm"
)

type XOrmAdapter struct {
	*xorm.Engine
}

func (this *XOrmAdapter) Name() string {
	return "XOrmAdapter"
}

func NewXOrmAdapter() *XOrmAdapter {
	var dsn string
	if config := Injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		dsn = config.(*SysConfig).Database.Uri
	} else {
		log.Fatal("链接数据库失败：dsn地址读取失败")
	}
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal("链接数据库失败", err)
	}
	engine.DB().SetMaxIdleConns(5)
	engine.DB().SetMaxOpenConns(10)
	return &XOrmAdapter{Engine: engine}
}
