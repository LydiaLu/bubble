package data

import (
	"bubble/internal/conf"
	"errors"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewTodoRepo, NewDB)

// Data .
type Data struct {
	// TODO wrapped database client
	db *gorm.DB
	// db sqlx.DB
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 如果在这里直接使用gorm连接DB，就不符合控制反转要求
	// db, err := gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{})
	// 正确方法应该是使用依赖注入，将依赖参数传进来
	return &Data{db: db}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	//根据配置文件中指定的driver连接不同的数据库
	switch strings.ToLower(c.Database.Driver) {
	case "mysql":
		return gorm.Open(mysql.Open(c.Database.Dsn), &gorm.Config{})
	case "sqllite":
		return gorm.Open(sqlite.Open(c.Database.Dsn), &gorm.Config{})
	}
	return nil, errors.New("invalid driver")
}
