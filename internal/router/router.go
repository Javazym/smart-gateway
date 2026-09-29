package router

import (
	"github.com/gin-gonic/gin"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/config"
	"smart-gateway/internal/handler"
	"smart-gateway/internal/middleware"
	"smart-gateway/internal/repository"
	"smart-gateway/internal/service"
)

// New 根据配置构建引擎并注册所有路由。
func New(db *clients.MysqlDB, rdb *clients.RedisDB, jwtCfg config.JwtConfig, secCfg config.SecurityConfig) *gin.Engine {
	r := gin.New()
	requestLogService := service.NewRequestLogService(db)
	// 全局请求日志：转发请求 + 管理 API + 登录 + 探活全部落 request_logs。
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(), middleware.NewRequestLog(requestLogService).Handle())
	r.GET("/healthz", handler.Health)

	// 认证
	authService := service.NewAuthService(db, jwtCfg.Secret, jwtCfg.Expire)
	authController := handler.NewAuthController(authService)
	authGroup := r.Group("/gateway_auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/logout", authController.Logout)
		authGroup.POST("/add_user", middleware.Auth(jwtCfg.Secret), middleware.RequireAdmin(), authController.AddUser)
	}

	// 各业务 service/controller
	userService := service.NewUserService(db)
	upstreamService := service.NewUpstreamService(db)
	routeService := service.NewRouteService(db)
	deviceService := service.NewDeviceService(db)
	banService := service.NewBanService(db)
	circuitBreakerService := service.NewCircuitBreakerService(db)
	rateLimitService := service.NewRateLimitService(db)
	auditLogService := service.NewAuditLogService(db)

	userController := handler.NewUserController(userService)
	upstreamController := handler.NewUpstreamController(upstreamService)
	routeController := handler.NewRouteController(routeService)
	deviceController := handler.NewDeviceController(deviceService)
	banController := handler.NewBanController(banService)
	circuitBreakerController := handler.NewCircuitBreakerController(circuitBreakerService)
	rateLimitController := handler.NewRateLimitController(rateLimitService)
	requestLogController := handler.NewRequestLogController(requestLogService)
	auditLogController := handler.NewAuditLogController(auditLogService)

	// 受保护的管理 API：鉴权 → 角色校验 → 审计。
	api := r.Group("/api", middleware.Auth(jwtCfg.Secret), middleware.RequireRole(), middleware.NewAudit(auditLogService).Handle())
	{
		// 用户
		api.GET("/users", userController.List)
		api.GET("/users/:id", userController.Get)
		api.POST("/users", userController.Create)
		api.PUT("/users/:id", userController.Update)
		api.DELETE("/users/:id", userController.Delete)

		// 上游分组
		api.GET("/upstream-groups", upstreamController.ListGroup)
		api.GET("/upstream-groups/:id", upstreamController.GetGroup)
		api.POST("/upstream-groups", upstreamController.CreateGroup)
		api.PUT("/upstream-groups/:id", upstreamController.UpdateGroup)
		api.DELETE("/upstream-groups/:id", upstreamController.DeleteGroup)

		// 上游节点
		api.GET("/upstreams", upstreamController.ListUpstream)
		api.GET("/upstreams/:id", upstreamController.GetUpstream)
		api.POST("/upstreams", upstreamController.CreateUpstream)
		api.PUT("/upstreams/:id", upstreamController.UpdateUpstream)
		api.DELETE("/upstreams/:id", upstreamController.DeleteUpstream)

		// 路由
		api.GET("/routes", routeController.List)
		api.GET("/routes/:id", routeController.Get)
		api.POST("/routes", routeController.Create)
		api.PUT("/routes/:id", routeController.Update)
		api.DELETE("/routes/:id", routeController.Delete)

		// 设备
		api.GET("/devices", deviceController.List)
		api.GET("/devices/:id", deviceController.Get)
		api.PUT("/devices/:id/status", deviceController.UpdateStatus)
		api.DELETE("/devices/:id", deviceController.Delete)

		// 封禁
		api.GET("/bans", banController.List)
		api.GET("/bans/:id", banController.Get)
		api.POST("/bans", banController.Create)
		api.POST("/bans/:id/unban", banController.Unban)
		api.DELETE("/bans/:id", banController.Delete)

		// 熔断
		api.GET("/circuit-breakers", circuitBreakerController.List)
		api.GET("/circuit-breakers/:id", circuitBreakerController.Get)
		api.POST("/circuit-breakers", circuitBreakerController.Create)
		api.PUT("/circuit-breakers/:id", circuitBreakerController.Update)
		api.DELETE("/circuit-breakers/:id", circuitBreakerController.Delete)

		// 限流
		api.GET("/rate-limits", rateLimitController.List)
		api.GET("/rate-limits/:id", rateLimitController.Get)
		api.POST("/rate-limits", rateLimitController.Create)
		api.PUT("/rate-limits/:id", rateLimitController.Update)
		api.DELETE("/rate-limits/:id", rateLimitController.Delete)

		// 日志（只读）
		api.GET("/request-logs", requestLogController.List)
		api.GET("/audit-logs", auditLogController.List)
	}

	// 请求转发管线：设备识别 → 封禁 → 路由/负载均衡 → 限流 → 熔断 → 反向代理。
	// （请求日志已提升为全局中间件，此处不再重复记录。）
	requestService := service.NewRequestService(
		repository.NewRouteRepo(db),
		repository.NewUpstreamRepo(db),
		repository.NewRateLimitPolicyRepo(db),
		repository.NewCircuitBreakerRepo(db),
		rdb,
		secCfg.SSRFProtection,
	)
	r.NoRoute(
		middleware.NewAddDevice(deviceService).HandleAddDevice(),
		middleware.NewBan(banService).HandleBan(),
		middleware.NewRobin(requestService, rdb).HandleRobin(),
		middleware.NewRateLimit(requestService, rdb).LimitRate(),
		middleware.NewCircuit(requestService, rdb).HandleCircuit(),
		middleware.NewProxy(requestService).HandleProxy(),
	)

	return r
}
