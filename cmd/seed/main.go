package main

import (
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/config"
	"smart-gateway/internal/model"
)

func main() {
	mysqlCfg, err := config.LoadMysqlConfig()
	if err != nil {
		log.Fatalf("load mysql config: %v", err)
	}
	db, err := clients.NewMysqlDB(mysqlCfg)
	if err != nil {
		log.Fatalf("init mysql: %v", err)
	}
	defer clients.CloseMysql(db)

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

	ctx := context.Background()

	// 幂等：路由表已有数据则跳过，避免重复插入。
	var count int64
	if err := db.WithContext(ctx).Model(&model.Route{}).Count(&count).Error; err != nil {
		log.Fatalf("count routes: %v", err)
	}
	if count > 0 {
		log.Printf("routes 表已有 %d 条数据，跳过插入", count)
		return
	}

	if err := seedForwardingData(ctx, db); err != nil {
		log.Fatalf("seed: %v", err)
	}
	log.Println("20 组转发链路数据插入完成")
}

type svc struct {
	name  string
	match string
}

var services = []svc{
	{"user-service", "/api/user"},
	{"order-service", "/api/order"},
	{"payment-service", "/api/payment"},
	{"inventory-service", "/api/inventory"},
	{"notification-service", "/api/notification"},
	{"search-service", "/api/search"},
	{"analytics-service", "/api/analytics"},
	{"auth-service", "/api/auth"},
	{"billing-service", "/api/billing"},
	{"shipping-service", "/api/shipping"},
	{"review-service", "/api/review"},
	{"cart-service", "/api/cart"},
	{"recommendation-service", "/api/recommendation"},
	{"report-service", "/api/report"},
	{"scheduler-service", "/api/scheduler"},
	{"media-service", "/api/media"},
	{"chat-service", "/api/chat"},
	{"email-service", "/api/email"},
	{"logging-service", "/api/logging"},
	{"config-service", "/api/config"},
}

func seedForwardingData(ctx context.Context, db *clients.MysqlDB) error {
	lbAlgos := []string{"round_robin", "weighted_rr", "ip_hash", "least_conn"}
	rlAlgos := []string{"token_bucket", "fixed_window", "sliding_window", "leaky_bucket"}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, s := range services {
			group := &model.UpstreamGroup{
				Name:                s.name + "-group",
				LBAlgorithm:         lbAlgos[i%len(lbAlgos)],
				HealthCheckPath:     strPtr("/health"),
				HealthCheckInterval: 10,
				Status:              1,
			}
			if err := tx.Create(group).Error; err != nil {
				return fmt.Errorf("create group %s: %w", s.name, err)
			}

			for j := 0; j < 2; j++ {
				up := &model.Upstream{
					GroupID:      group.ID,
					Name:         fmt.Sprintf("%s-node-%d", s.name, j+1),
					Scheme:       "http",
					Host:         "127.0.0.1",
					Port:         9000 + i*10 + j,
					Weight:       j + 1,
					MaxConns:     0,
					HealthStatus: "unknown",
					Status:       1,
				}
				if err := tx.Create(up).Error; err != nil {
					return fmt.Errorf("create upstream %s: %w", s.name, err)
				}
			}

			var rateLimitID uint64
			if i%2 == 0 {
				rl := &model.RateLimitPolicy{
					Name:          s.name + "-rl",
					Algorithm:     rlAlgos[(i/2)%len(rlAlgos)],
					Rate:          100,
					Burst:         200,
					WindowSeconds: 1,
					KeyType:       "ip",
					Status:        1,
				}
				if err := tx.Create(rl).Error; err != nil {
					return fmt.Errorf("create rate limit %s: %w", s.name, err)
				}
				rateLimitID = rl.ID
			}

			var circuitBreakerID uint64
			if i%2 == 1 {
				cb := &model.CircuitBreaker{
					Name:                s.name + "-cb",
					FailureThreshold:    5,
					SuccessThreshold:    3,
					TimeoutMs:           3000,
					OpenDurationSeconds: 30,
					HalfOpenMaxRequests: 1,
					Status:              1,
				}
				if err := tx.Create(cb).Error; err != nil {
					return fmt.Errorf("create circuit breaker %s: %w", s.name, err)
				}
				circuitBreakerID = cb.ID
			}

			route := &model.Route{
				Name:             s.name,
				MatchType:        "prefix",
				MatchPattern:     s.match,
				UpstreamGroupID:  group.ID,
				RateLimitID:      rateLimitID,
				CircuitBreakerID: circuitBreakerID,
				StripPrefix:      "",
				Priority:         i + 1,
				Status:           1,
			}
			if err := tx.Create(route).Error; err != nil {
				return fmt.Errorf("create route %s: %w", s.name, err)
			}
		}
		return nil
	})
}

func strPtr(s string) *string { return &s }
