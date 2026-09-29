package robin

import (
	"fmt"
	"hash/fnv"
	"net/http"
)

// IP 哈希 ip_hash
type IpHash struct {
	nodes []*Node
}

func NewIpHash(nodes []*Node) *IpHash {
	return &IpHash{nodes: filterHealth(nodes)}
}

func (ih *IpHash) Pick(r *http.Request) (*Node, error) {
	health := filterHealth(ih.nodes)
	if len(health) == 0 {
		return nil, fmt.Errorf("no healthy node")
	}
	ip := GetClientIP(r)
	h := fnv.New32a()
	h.Write([]byte(ip))
	hashVal := h.Sum32()
	idx := int(hashVal) % len(health)
	return health[idx], nil
}
func (ih *IpHash) ConnInc(addr string) {}
func (ih *IpHash) ConnDec(addr string) {}
