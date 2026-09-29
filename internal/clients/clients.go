package clients

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"smart-gateway/internal/config"
	"time"
)

type MysqlDB struct {
	*gorm.DB
}

type RedisDB struct {
	*redis.Client
}

func NewMysqlDB(conf *config.MysqlConfig) (*MysqlDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Database,
		conf.Charset)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	dur, _ := time.ParseDuration(conf.ConnMaxLifetime)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dur)

	// ping探测
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("mysql ping failed: %w", err)
	}
	log.Println("[database] mysql connected")
	return &MysqlDB{
		DB: db,
	}, nil
}

func CloseMysql(db *MysqlDB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB.DB()
	if err == nil {
		_ = sqlDB.Close()
		log.Println("[database] mysql closed")
	}
}

func NewRedisDB(conf *config.RedisConfig) (*RedisDB, error) {
	cxt := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		Password:     "",
		DB:           conf.DB,
		PoolSize:     conf.PoolSize,
		MinIdleConns: conf.MinIdleConns,
	})
	_, err := rdb.Ping(cxt).Result()
	if err != nil {
		return nil, err
	}
	return &RedisDB{
		Client: rdb,
	}, nil
}
