package robin

import (
	"fmt"
	"net/http"
	"time"

	"smart-gateway/internal/clients"
)

// WeightRoundRobin 加权轮询。
// 通过 Redis INCR 取得全局自增序号，按权重区间确定性落点，避免多实例下的读改写竞态。
type WeightRoundRobin struct {
	nodes       []*Node
	totalWeight int
	rdb         *clients.RedisDB
	key         string
}

func NewWeightRoundRobin(nodes []*Node, rdb *clients.RedisDB, key string) *WeightRoundRobin {
	health := filterHealth(nodes)
	total := 0
	for _, n := range health {
		total += n.Weight
	}
	return &WeightRoundRobin{nodes: health, totalWeight: total, rdb: rdb, key: key}
}

func (w *WeightRoundRobin) Pick(r *http.Request) (*Node, error) {
	if len(w.nodes) == 0 || w.totalWeight <= 0 {
		return nil, fmt.Errorf("no healthy upstream node")
	}
	seq, err := w.rdb.Incr(r.Context(), w.key).Result()
	if err != nil {
		seq = time.Now().UnixNano()
	}
	x := (seq - 1) % int64(w.totalWeight)
	acc := 0
	for _, node := range w.nodes {
		acc += node.Weight
		if x < int64(acc) {
			return node, nil
		}
	}
	return w.nodes[len(w.nodes)-1], nil
}

func (w *WeightRoundRobin) ConnInc(addr string) {}
func (w *WeightRoundRobin) ConnDec(addr string) {}
