package db

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/pkg/errors"
)

func NewDB(gormDB *gorm.DB) *sql.DB {
	return gormDB.DB()
}

func NewGormDB(o *Options) (*gorm.DB, error) {
	if o.Type == "postgres" {
		return newPostgres(o)
	} else if o.Type == "mysql" {
		return newMySQL(o)
	}
	return nil, errors.New("Not supported db")
}

func newPostgres(o *Options) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		o.Host,
		o.Port,
		o.Username,
		o.Password,
		o.Database)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func newMySQL(o *Options) (*gorm.DB, error) {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Asia%%2FShanghai&parseTime=true",
		o.Username,
		o.Password,
		o.Host,
		o.Port,
		o.Database,
		o.ChartSet,
	)

	db, err := gorm.Open("mysql", mysqlInfo)
	if err != nil {
		return nil, err
	}
	return db, err
}
