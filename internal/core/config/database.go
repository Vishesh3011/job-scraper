package config

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"syscall"
)

type dbConfig struct {
	hostname                     string
	port                         int
	username, password, database string
}

func newDBConfig() (*dbConfig, error) {
	hostname, found := syscall.Getenv("DB_HOSTNAME")
	if !found {
		err := errors.New("DB_HOSTNAME environment variable not set")
		return nil, err
	}

	port, found := syscall.Getenv("DB_PORT")
	if !found {
		err := errors.New("DB_PORT environment variable not set")
		return nil, err
	}

	username, found := syscall.Getenv("DB_USERNAME")
	if !found {
		err := errors.New("DB_USERNAME environment variable not set")
		return nil, err
	}

	password, found := syscall.Getenv("DB_PASSWORD")
	if !found {
		err := errors.New("DB_PASSWORD environment variable not set")
		return nil, err
	}

	database, found := syscall.Getenv("DB_DATABASE")
	if !found {
		err := errors.New("DB_DATABASE environment variable not set")
		return nil, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("DB_PORT must be an integer: %v", err))
	}

	return &dbConfig{
		hostname: hostname,
		port:     portInt,
		username: username,
		password: password,
		database: database,
	}, nil
}

func (db *dbConfig) Hostname() string {
	return db.hostname
}

func (db *dbConfig) Port() int {
	return db.port
}

func (db *dbConfig) Username() string {
	return db.username
}

func (db *dbConfig) Password() string {
	return db.password
}

func (db *dbConfig) Database() string {
	return db.database
}

func (db *dbConfig) DBConn() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", db.Username(), db.Password(), db.Hostname(), db.Port(), db.Database())
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		err = fmt.Errorf("failed to connect to database: %v", err)
		return nil, err
	}
	return conn, nil
}
