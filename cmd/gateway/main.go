package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/config"
	"smart-gateway/internal/model"
	"smart-gateway/internal/router"
	"smart-gateway/internal/service"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	mysql, mysqlErr := config.LoadMysqlConfig()
	if mysqlErr != nil {
		log.Fatalf("load mysql config: %v", mysqlErr)
	}
	db, dbErr := clients.NewMysqlDB(mysql)
	if dbErr != nil {
		log.Fatalf("init mysql: %v", dbErr)
	}
	redis, redisErr := config.LoadRedisConfig()
	if redisErr != nil {
		log.Fatalf("load redis config: %v", redisErr)
	}
	rdb, rdbErr := clients.NewRedisDB(redis)
	if rdbErr != nil {
		log.Fatalf("init redis: %v", rdbErr)
	}
	defer func() {
		clients.CloseMysql(db)
	}()
	if err := db.AutoMigrate(
		&model.UpstreamGroup{},
		&model.Upstream{},
		&model.User{},
		&model.AuditLog{},
		&model.Ban{},
		&model.Device{},
		&model.CircuitBreaker{},
		&model.Route{},
		&model.RateLimitPolicy{},
		&model.RequestLog{},
	); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}
	if err := service.EnsureAdmin(db); err != nil {
		log.Fatalf("ensure admin: %v", err)
	}
	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router.New(db, rdb, cfg.Jwt, cfg.Security),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		log.Printf("gateway listening on %s", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server exited")
}
