package utils

import (
	"bytes"
	"net/url"
	"runtime"
	"strconv"
	"strings"
)

func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	i := bytes.IndexByte(b, ' ')
	gid, _ := strconv.ParseUint(string(b[:i]), 10, 64)
	return gid
}

func Condition[T any](cond bool, a, b T) T {
	if cond {
		return a
	} else {
		return b
	}
}

// EncodeUrl 编码 URL
func EncodeUrl(rawURL string) (string, error) {
	// 解析 URL
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	// 将 URL 的路径部分按 "/" 分割，对每个分段进行编码
	segments := strings.Split(u.Path, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	encodedPath := strings.Join(segments, "/")
	// 保留 u.Path 的原始值，仅设置 u.RawPath 为编码后的字符串
	u.RawPath = encodedPath

	// 输出编码后的 URL
	return u.String(), err
}
