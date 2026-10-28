package repository

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"io"
	"log"
	"time"
)

type MySQL struct{ db *sql.DB }

func Open(addr, name, user, password string, timeout time.Duration) (*MySQL, error) {
	cfg := mysql.NewConfig()
	cfg.Logger = log.New(io.Discard, "", 0)
	cfg.User = user
	cfg.Passwd = password
	cfg.Net = "tcp"
	cfg.Addr = addr
	cfg.DBName = name
	cfg.Timeout = timeout
	cfg.ReadTimeout = timeout
	cfg.WriteTimeout = timeout
	connector, e := mysql.NewConnector(cfg)
	if e != nil {
		return nil, e
	}
	db := sql.OpenDB(connector)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(3 * time.Minute)
	return &MySQL{db}, nil
}
func (m *MySQL) Ping(ctx context.Context) error { return m.db.PingContext(ctx) }
func (m *MySQL) Close() error                   { return m.db.Close() }
