package robin

import (
	"fmt"
	"hash/fnv"
	"net/http"
	"sort"
)

// ConsistentHash 一致性哈希负载均衡，用于缓存亲和性路由。
// 每个节点映射多个虚拟节点到哈希环，相同 key（客户端 IP）始终落到同一节点。
type ConsistentHash struct {
	nodes  []*Node
	ring   []hashEntry
	vnodes int
}

type hashEntry struct {
	hash uint32
	addr string
}

const defaultVirtualNodes = 100

func NewConsistentHash(nodes []*Node) *ConsistentHash {
	ch := &ConsistentHash{nodes: filterHealth(nodes), vnodes: defaultVirtualNodes}
	ch.buildRing()
	return ch
}

func (ch *ConsistentHash) buildRing() {
	ch.ring = make([]hashEntry, 0, len(ch.nodes)*ch.vnodes)
	for _, n := range ch.nodes {
		for i := 0; i < ch.vnodes; i++ {
			ch.ring = append(ch.ring, hashEntry{
				hash: hashKey(fmt.Sprintf("%s#%d", n.Addr, i)),
				addr: n.Addr,
			})
		}
	}
	sort.Slice(ch.ring, func(i, j int) bool { return ch.ring[i].hash < ch.ring[j].hash })
}

func (ch *ConsistentHash) Pick(r *http.Request) (*Node, error) {
	if len(ch.nodes) == 0 {
		return nil, fmt.Errorf("no healthy upstream node")
	}
	h := hashKey(GetClientIP(r))
	idx := sort.Search(len(ch.ring), func(i int) bool { return ch.ring[i].hash >= h })
	if idx == len(ch.ring) {
		idx = 0
	}
	addr := ch.ring[idx].addr
	for _, n := range ch.nodes {
		if n.Addr == addr {
			return n, nil
		}
	}
	return ch.nodes[0], nil
}

func (ch *ConsistentHash) ConnInc(addr string) {}
func (ch *ConsistentHash) ConnDec(addr string) {}

func hashKey(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
