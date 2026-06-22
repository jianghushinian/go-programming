package config

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig defines config for mysql database.
type MySQLConfig struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	// +optional
	Logger logger.Interface
}

// DSN return DSN from MySQLConfig.
func (o *MySQLConfig) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`,
		o.Username,
		o.Password,
		o.Host,
		o.Database,
		true,
		"Local")
}

// NewMySQL create a new gorm db instance with the given config.
func NewMySQL(cfg *MySQLConfig) (*gorm.DB, error) {
	// Set default values to ensure all fields in cfg are available.
	setDefaults(cfg)

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{
		// PrepareStmt executes the given query in cached statement.
		// This can improve performance.
		PrepareStmt: true,
		Logger:      cfg.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(cfg.MaxConnectionLifeTime)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	return db, nil
}

// setDefaults set available default values for some fields.
func setDefaults(cfg *MySQLConfig) {
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1:3306"
	}
	if cfg.MaxIdleConnections == 0 {
		cfg.MaxIdleConnections = 100
	}
	if cfg.MaxOpenConnections == 0 {
		cfg.MaxOpenConnections = 100
	}
	if cfg.MaxConnectionLifeTime == 0 {
		cfg.MaxConnectionLifeTime = time.Duration(10) * time.Second
	}
	if cfg.Logger == nil {
		cfg.Logger = logger.Default
	}
}
