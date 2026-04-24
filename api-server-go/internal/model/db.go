package model

import (
	"fmt"
	"log"
	"time"

	"mochat-api-server/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB(cfg config.DBConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.Collation,
	)

	gormLogger := logger.Default.LogMode(logger.Info)
	if cfg.Driver != "mysql" {
		return fmt.Errorf("unsupported db driver: %s", cfg.Driver)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: &prefixNamingStrategy{
			TablePrefix: cfg.Prefix,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MinConns)
	sqlDB.SetMaxOpenConns(cfg.MaxConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxIdleTime) * time.Second)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.MaxIdleTime) * time.Second)

	log.Println("Database connected successfully")
	return nil
}

type prefixNamingStrategy struct {
	TablePrefix string
}

func (s *prefixNamingStrategy) TableName(table string) string {
	return s.TablePrefix + table
}

func (s *prefixNamingStrategy) SchemaName(table string) string {
	return ""
}

func (s *prefixNamingStrategy) ColumnName(table, column string) string {
	return schema.NamingStrategy{}.ColumnName(table, column)
}

func (s *prefixNamingStrategy) JoinTableName(table string) string {
	return s.TablePrefix + table
}

func (s *prefixNamingStrategy) RelationshipFKName(rel schema.Relationship) string {
	return ""
}

func (s *prefixNamingStrategy) CheckerName(table, column string) string {
	return ""
}

func (s *prefixNamingStrategy) IndexName(table, column string) string {
	return ""
}

func (s *prefixNamingStrategy) UniqueName(table, column string) string {
	return ""
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
