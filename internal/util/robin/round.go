package robin

import (
	"fmt"
	"net/http"
	"time"

	"smart-gateway/internal/clients"
)

// RoundRobin 轮询。
type RoundRobin struct {
	nodes []*Node
	rdb   *clients.RedisDB
	key   string
}

func NewRoundRobin(nodes []*Node, rdb *clients.RedisDB, key string) *RoundRobin {
	return &RoundRobin{nodes: filterHealth(nodes), rdb: rdb, key: key}
}

func (rr *RoundRobin) Pick(r *http.Request) (*Node, error) {
	health := filterHealth(rr.nodes)
	if len(health) == 0 {
		return nil, fmt.Errorf("no healthy upstream node")
	}
	// INCR 原子自增，保证多实例轮询不重不漏。
	seq, err := rr.rdb.Incr(r.Context(), rr.key).Result()
	if err != nil {
		seq = time.Now().UnixNano()
	}
	idx := (seq - 1) % int64(len(health))
	return health[idx], nil
}

func (rr *RoundRobin) ConnInc(addr string) {}
func (rr *RoundRobin) ConnDec(addr string) {}
