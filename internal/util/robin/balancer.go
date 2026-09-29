package robin

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

// Node 后端实例。
type Node struct {
	ID     uint64 // 上游节点数据库主键，供请求日志关联使用
	Addr   string // 地址 host:port
	Scheme string // 协议 http/https
	Weight int    // 权重（加权轮询使用）
	Status int    // 1 健康可用
}

// Balancer 负载均衡器接口。
type Balancer interface {
	// Pick 根据请求选出目标节点。
	Pick(r *http.Request) (*Node, error)
	// ConnInc / ConnDec 维护活跃连接计数（最小连接数使用）。
	ConnInc(addr string)
	ConnDec(addr string)
}

// NewBalancer 根据分组配置构建负载均衡器。
func NewBalancer(group *model.UpstreamGroup, rdb *clients.RedisDB) Balancer {
	nodes := make([]*Node, 0, len(group.Upstreams))
	for _, up := range group.Upstreams {
		scheme := up.Scheme
		if scheme == "" {
			scheme = "http"
		}
		nodes = append(nodes, &Node{
			ID:     up.ID,
			Addr:   up.Host + ":" + strconv.Itoa(up.Port),
			Scheme: scheme,
			Weight: up.Weight,
			Status: up.Status,
		})
	}
	key := fmt.Sprintf("balancer:%d", group.ID)
	switch group.LBAlgorithm {
	case "weighted_rr", "weight_round_robin":
		return NewWeightRoundRobin(nodes, rdb, key+":wrr")
	case "ip_hash":
		return NewIpHash(nodes)
	case "least_conn":
		return NewLeastConn(nodes, rdb, key+":conn")
	case "consistent_hash":
		return NewConsistentHash(nodes)
	default: // round_robin 及未知算法均回落轮询
		return NewRoundRobin(nodes, rdb, key+":rr")
	}
}

// filterHealth 过滤下线节点。
func filterHealth(nodes []*Node) []*Node {
	res := make([]*Node, 0, len(nodes))
	for _, n := range nodes {
		if n.Status == 1 {
			res = append(res, n)
		}
	}
	return res
}

// GetClientIP 获取真实客户端 IP。
// 仅信任 TCP 对端地址（RemoteAddr），不读取 X-Forwarded-For/X-Real-IP，
// 避免客户端伪造转发头绕过 IP 维度的限流与封禁。
func GetClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
