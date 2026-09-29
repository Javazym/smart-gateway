package util

import "net/url"

// SplitUrlByHostSlash 把url拆成 scheme://host 和 host之后第一个/开始的后缀
func SplitUrlByHostSlash(raw string) (prefix string, suffix string, err error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", "", err
	}
	prefix = u.Scheme + "://" + u.Host
	suffix = u.RequestURI()
	if suffix == "" {
		suffix = "/"
	}
	return
}
