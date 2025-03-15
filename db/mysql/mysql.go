// db/mysql/mysql.go
package mysql

import (
	"database/sql"
	"errors"
)

type MySQLConfig struct {
	DSN string
}

func (c *MySQLConfig) Validate() error {
	if c.DSN == "" {
		return errors.New("empty DSN")
	}
	return nil
}

func (c *MySQLConfig) DriverName() string {
	return "mysql"
}

type MySQLDriver struct {
	db *sql.DB
}

func (m *MySQLDriver) Connect(config interface{}) error {
	cfg, ok := config.(*MySQLConfig)
	if !ok {
		return errors.New("invalid config type")
	}
	db, err := sql.Open("mysql", cfg.DSN)
	m.db = db
	return err
}

// 实现其他接口方法...
