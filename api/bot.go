package api

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Bot struct {
	HttpClient *http.Client
	Debug      bool
}

// NewBot 创建一个新的Bot实例
func NewBot(proxyUrl string) (*Bot, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{}
	if proxyUrl != "" {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxy)
	}
	return &Bot{
		HttpClient: &http.Client{
			Jar:       jar,
			Transport: transport,
			Timeout:   30 * time.Second, // 添加30秒超时
		},
	}, nil
}

// SendPostRequest 发送请求
func (b *Bot) SendPostRequest(apiUrl string, data []byte, headerMap map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headerMap {
		req.Header.Set(k, v)
	}
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//读取body数据
	if resp.StatusCode != http.StatusOK {
		// 处理非 200 的情况
		all, _ := io.ReadAll(resp.Body)
		log.Printf("status code: %d, body: %s\n", resp.StatusCode, string(all))
	}
	return resp, nil
}

// SendGetRequest 发送请求
func (b *Bot) SendGetRequest(apiUrl string, headerMap map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headerMap {
		req.Header.Set(k, v)
	}
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//读取body数据
	if resp.StatusCode != http.StatusOK {
		// 处理非 200 的情况
		all, _ := io.ReadAll(resp.Body)
		log.Printf("status code: %d, body: %s\n", resp.StatusCode, string(all))
	}
	return resp, nil
}

// 格式刷
func FormatBrush(str string) (ret string) {
	code := 0
	//第一位不能为0 非数字字母后的第一个字母也不能为0
	no := true
	for _, char := range []rune(str) {
		code = int(char)
		if code >= 48 && code <= 57 { //数字
			if no {
				ret += string(rune(RandInt64(49, 57)))
				no = false
			} else {
				ret += string(rune(RandInt64(48, 57)))
			}
		} else if code >= 97 && code <= 102 { //a-f
			ret += string(rune(RandInt64(97, 102)))
		} else if code >= 103 && code <= 122 { //g-z
			ret += string(rune(RandInt64(103, 122)))
		} else if code >= 65 && code <= 70 { //A-F
			ret += string(rune(RandInt64(65, 70)))
		} else if code >= 71 && code <= 90 { //G-Z
			ret += string(rune(RandInt64(71, 90)))
		} else {
			//其他字符
			ret += string(rune(code))
			no = true
		}
	}
	return ret
}

func RandInt64(min, max int64) int64 {
	if min > max {
		return min // 或者其他错误处理
	}
	return rand.Int63n(max-min+1) + min
}

func encrypt(key string) string {
	// 将字符串转换为字节数组
	v := []byte(key)

	// 遍历字节数组并进行位异或操作
	for i := 0; i < len(v); i++ {
		v[i] = v[i] ^ 5
	}

	// 将字节数组转换为十六进制字符串
	hexString := hex.EncodeToString(v)

	// 返回小写形式的十六进制字符串
	return strings.ToLower(hexString)
}
