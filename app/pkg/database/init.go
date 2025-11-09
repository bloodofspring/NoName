package database

import (
	"app/pkg/database/models"
	e "app/pkg/errors"
	"os"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

var (
	db   *pg.DB
	once sync.Once
)

// GetDB returns a singleton instance of the database connection
func GetDB() *pg.DB {
	once.Do(func() {
		db = pg.Connect(&pg.Options{
			Addr:     os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
			PoolSize: 20, // Устанавливаем разумный размер пула
		})
	})
	return db
}

func InitDb() *e.ErrorInfo {
	db := GetDB()

	models := []interface{}{
		&models.User{},
		&models.Chat{},
		&models.Thread{},
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return e.Error(err, "Error creating table").
				WithSeverity(e.Critical).
				WithData(map[string]any{
					"model": model,
				})
		}
	}

	return e.Nil()
}
