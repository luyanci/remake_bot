package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewDatabase)

func NewDatabase(lc fx.Lifecycle, logger *zap.Logger) (*sql.DB, error) {
        user := os.Getenv("POSTGRES_USER")
        password := os.Getenv("POSTGRES_PASSWORD")
        host := os.Getenv("POSTGRES_HOST")
        port := os.Getenv("POSTGRES_PORT")
        dbname := os.Getenv("POSTGRES_DB")
        if user == "" {
        	user = "postgres"
        }
        if password == "" {
        	password = "114514"
        }
        if host == "" {
            	host = "localhost"
        }
        if port == "" {
           	port = "5432"
        }
        if dbname == "" {
        	dbname = "postgres"
        }
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// 在应用启动时验证连接
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			logger.Info("[db] Verifying connection...")
			return db.Ping()
		},
		OnStop: func(context.Context) error {
			logger.Info("[db] Closing connection...")
			return db.Close()
		},
	})

	return db, nil
}
