package utils

import (
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dataBase struct {
	*gorm.DB
}

var (
	once sync.Once
	DB   *dataBase
)

func GetDBInstance() {
	once.Do(func() {
		dsn := `host=localhost user=postgres password=root 
		dbname=task_test port=5432 sslmode=disable TimeZone=Asia/Shanghai`
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Failed to connect to the database")
		}
		DB = &dataBase{db}
	})
}
