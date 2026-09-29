package util

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sort"
	"strings"
)

// volatileHeaders 参与设备指纹时需剔除的易变/敏感请求头，避免指纹随会话或链路变化。
var volatileHeaders = map[string]struct{}{
	"cookie":          {},
	"authorization":   {},
	"x-forwarded-for": {},
	"x-real-ip":       {},
	"x-request-id":    {},
}

// HeaderHash 对请求头规范化（键小写、按名+值排序、剔除易变头）后取 sha256，作为设备画像的一部分。
// 注意：纯 HTTP 层无 TLS 指纹，仅凭头特征区分度有限。
func HeaderHash(r *http.Request) string {
	type kv struct{ k, v string }
	pairs := make([]kv, 0, len(r.Header))
	for name, values := range r.Header {
		lower := strings.ToLower(name)
		if _, skip := volatileHeaders[lower]; skip {
			continue
		}
		for _, v := range values {
			pairs = append(pairs, kv{lower, v})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].k != pairs[j].k {
			return pairs[i].k < pairs[j].k
		}
		return pairs[i].v < pairs[j].v
	})

	h := sha256.New()
	for _, p := range pairs {
		h.Write([]byte(p.k))
		h.Write([]byte{'='})
		h.Write([]byte(p.v))
		h.Write([]byte{'\n'})
	}
	return hex.EncodeToString(h.Sum(nil))
}

// FingerprintID 组合 User-Agent 与头哈希生成设备综合指纹。
// 后续接入 TLS ClientHello（JA4/JA3）与 OS 指纹后可在此追加分量。
func FingerprintID(userAgent, headerHash string) string {
	h := sha256.New()
	h.Write([]byte(userAgent))
	h.Write([]byte{'|'})
	h.Write([]byte(headerHash))
	return hex.EncodeToString(h.Sum(nil))
}
