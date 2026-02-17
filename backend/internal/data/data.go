package data

import (
	"database/sql"
	"fmt"

	"backend/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewGreeterRepo,
	NewArticleRepo,
	NewDictionaryRepo,
	NewWordRepo,
	NewLearnRecordRepo,
	NewUploadTaskRepo,
	NewUserRepo,
	NewRefreshTokenRepo,
)

// Data .
type Data struct {
	db *sql.DB
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	d := &Data{db: db}
	cleanup := func() {
		log.Info("closing the data resources")
		db.Close()
	}
	return d, cleanup, nil
}
