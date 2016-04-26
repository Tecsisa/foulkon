package postgresql

import (
	_ "github.com/lib/pq"

	"github.com/jinzhu/gorm"
	"time"
)

type PostgresRepo struct {
	Dbmap *gorm.DB
}

func InitDb(datasourcename string) (*gorm.DB, error) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := gorm.Open("postgres", datasourcename)
	if err != nil {
		return nil, err
	}

	// construct a gorp DbMap
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.DB().SetConnMaxLifetime(5 * time.Minute)

	// Check connection
	err = db.DB().Ping()
	if err != nil {
		return nil, err
	}

	// Create tables if not exist
	err = db.AutoMigrate(&User{}).Error
	if err != nil {
		return nil, err
	}

	// Activate sql logger
	db.LogMode(true)

	return db, nil
}
