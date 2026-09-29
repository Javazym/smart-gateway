package robin

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"smart-gateway/internal/clients"
)

// LeastConn 最小连接数。
// 各节点活跃连接数存于 Redis hash，跨实例共享。
type LeastConn struct {
	nodes []*Node
	rdb   *clients.RedisDB
	key   string
}

func NewLeastConn(nodes []*Node, rdb *clients.RedisDB, key string) *LeastConn {
	return &LeastConn{nodes: filterHealth(nodes), rdb: rdb, key: key}
}

func (lc *LeastConn) Pick(r *http.Request) (*Node, error) {
	health := filterHealth(lc.nodes)
	if len(health) == 0 {
		return nil, fmt.Errorf("no healthy upstream node")
	}
	counts, err := lc.rdb.HGetAll(r.Context(), lc.key).Result()
	if err != nil {
		return health[0], nil
	}
	best := health[0]
	bestConn := int64(1) << 60
	for _, node := range health {
		var c int64
		if raw, ok := counts[node.Addr]; ok {
			c, _ = strconv.ParseInt(raw, 10, 64)
		}
		if c < bestConn {
			bestConn = c
			best = node
		}
	}
	return best, nil
}

func (lc *LeastConn) ConnInc(addr string) {
	_ = lc.rdb.HIncrBy(context.Background(), lc.key, addr, 1).Err()
}

func (lc *LeastConn) ConnDec(addr string) {
	_ = lc.rdb.HIncrBy(context.Background(), lc.key, addr, -1).Err()
}
