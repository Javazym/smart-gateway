package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Easonliuliang/helloprint/hash"
	"github.com/Easonliuliang/helloprint/parse"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
	"smart-gateway/internal/repository"
)

// ErrRouteNotFound 表示请求未命中任何启用的路由，用于与「上游不可用」区分。
var ErrRouteNotFound = errors.New("route not found")

type RequestService struct {
	rr  *repository.RouteRepo
	ur  *repository.UpstreamRepo
	rlr *repository.RateLimitPolicyRepo
	cb  *repository.CircuitBreakerRepo
	rdb *clients.RedisDB

	// ssrfProtection 是否拦截内网/回环/链路本地等上游地址。
	ssrfProtection bool
}

func NewRequestService(rr *repository.RouteRepo, ur *repository.UpstreamRepo, rlr *repository.RateLimitPolicyRepo, cb *repository.CircuitBreakerRepo, rdb *clients.RedisDB, ssrfProtection bool) *RequestService {
	return &RequestService{
		rr:             rr,
		ur:             ur,
		rlr:            rlr,
		cb:             cb,
		rdb:            rdb,
		ssrfProtection: ssrfProtection,
	}
}

// MatchRoute 按匹配类型与优先级从启用的路由中选出命中的路由。
// match_type 支持 path(精确路径)/prefix(前缀)/domain(域名)/header(请求头 Key=Value)，
// 空或未知类型按精确路径处理；多条命中时 priority 值小者优先。
func (r *RequestService) MatchRoute(ctx context.Context, req *http.Request) (*model.Route, error) {
	routes, err := r.rr.FindAllEnabled(ctx)
	if err != nil {
		return nil, err
	}
	for _, route := range routes {
		if matchRoute(route, req) {
			return &route, nil
		}
	}
	return nil, ErrRouteNotFound
}

func (r *RequestService) HandlerRateRequest(ctx context.Context, route *model.Route) (*model.RateLimitPolicy, error) {
	// 路由未关联限流策略，返回 nil 表示不限流。
	if route.RateLimitID == 0 {
		return nil, nil
	}
	return r.rlr.FindByID(ctx, route.RateLimitID)
}

func (r *RequestService) HandlerRobinRequest(ctx context.Context, route *model.Route) (*model.UpstreamGroup, error) {
	return r.ur.FindUpstreamGroupById(ctx, route.UpstreamGroupID)
}

func (r *RequestService) HandleCircuitRequest(ctx context.Context, route *model.Route) (*model.CircuitBreaker, error) {
	// 路由未关联熔断配置，直接返回 nil 表示放行。
	if route.CircuitBreakerID == 0 {
		return nil, nil
	}
	return r.cb.FindByID(ctx, route.CircuitBreakerID)
}

// matchRoute 依据路由匹配类型判断请求是否命中。
func matchRoute(route model.Route, req *http.Request) bool {
	switch route.MatchType {
	case "prefix":
		return strings.HasPrefix(req.URL.Path, route.MatchPattern)
	case "domain":
		return matchDomain(req.Host, route.MatchPattern)
	case "header":
		k, v := splitHeaderMatch(route.MatchPattern)
		if k == "" {
			return false
		}
		return req.Header.Get(k) == v
	default: // path 及未知类型按精确路径匹配
		return req.URL.Path == route.MatchPattern
	}
}

// matchDomain 比较 Host（忽略大小写），必要时剥离端口。
func matchDomain(host, pattern string) bool {
	if host == "" {
		return false
	}
	if strings.EqualFold(host, pattern) {
		return true
	}
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		return strings.EqualFold(h, pattern)
	}
	return false
}

// splitHeaderMatch 解析 "Key=Value" 或 "Key:Value" 形式的请求头匹配表达式。
func splitHeaderMatch(pattern string) (string, string) {
	if i := strings.Index(pattern, "="); i >= 0 {
		return strings.TrimSpace(pattern[:i]), strings.TrimSpace(pattern[i+1:])
	}
	if i := strings.Index(pattern, ":"); i >= 0 {
		return strings.TrimSpace(pattern[:i]), strings.TrimSpace(pattern[i+1:])
	}
	return "", ""
}

func (r *RequestService) NewProxy(targetURL string) *httputil.ReverseProxy {
	target, err := url.Parse(targetURL)
	if err != nil || target.Host == "" {
		return nil
	}
	// SSRF 防护（可选）：开启时拒绝内网/回环/链路本地等目标地址。
	if r.ssrfProtection && isForbiddenHost(target.Hostname()) {
		log.Printf("blocked proxy target (private address): %s", target.Hostname())
		return nil
	}
	proxy := httputil.NewSingleHostReverseProxy(target)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
	}

	// 转发失败（连接拒绝、超时等）统一返回 502，熔断中间件据此计入失败。
	proxy.ErrorHandler = func(w http.ResponseWriter, req *http.Request, err error) {
		log.Printf("proxy forward error: %v", err)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"code":502,"msg":"服务暂时不可用"}`))
	}
	return proxy
}

// isForbiddenHost 判断目标主机是否命中内网/回环/链路本地等禁止转发的地址。
func isForbiddenHost(host string) bool {
	if ip := net.ParseIP(host); ip != nil {
		return isForbiddenIP(ip)
	}
	// 主机名：解析后任一 IP 命中即拒绝，防 DNS 重绑定。
	ips, err := net.LookupIP(host)
	if err != nil {
		return true
	}
	for _, ip := range ips {
		if isForbiddenIP(ip) {
			return true
		}
	}
	return false
}

func isForbiddenIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsUnspecified()
}

func (r *RequestService) SavaDevice() {
	data, err := os.ReadFile("clienthello.bin")
	if err != nil {
		panic(err)
	}

	ch, err := parse.ParseClientHelloFromRecord(data)
	if err != nil {
		panic(err)
	}

	ja4Str := hash.JA4(ch)
	ja3Str := hash.JA3Hash(ch)

	fmt.Println("JA4_HASH =", ja4Str)
	fmt.Println("JA3_HASH =", ja3Str)
}
