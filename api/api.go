package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Api[T any] struct {
	apiUrl     string            // API URL
	apiMethod  string            // API 方法
	apiParams  url.Values        // API 参数
	apiHeaders map[string]string // API 头部
	apiBody    []byte            // API 请求体
	apiResp    T                 // API 响应
}

func NewApi[T any](method, apiUrl string) *Api[T] {
	//从apiUrl中解析出query
	u, err := url.Parse(apiUrl)
	if err != nil {
		return nil
	}
	queryParams := u.Query() // 获取查询参数
	u.RawQuery = ""          // 清除查询参数

	return &Api[T]{
		apiUrl:     u.String(),
		apiMethod:  method,
		apiParams:  queryParams,
		apiHeaders: make(map[string]string),
	}
}

func NewGetApi[T any](apiUrl string) *Api[T] {
	return NewApi[T](http.MethodGet, apiUrl)
}
func NewPostApi[T any](apiUrl string) *Api[T] {
	return NewApi[T](http.MethodPost, apiUrl)
}
func NewPutApi[T any](apiUrl string) *Api[T] {
	return NewApi[T](http.MethodPut, apiUrl)
}
func NewDeleteApi[T any](apiUrl string) *Api[T] {
	return NewApi[T](http.MethodDelete, apiUrl)
}

func (a *Api[T]) AddParam(k, v string) *Api[T] {
	a.apiParams.Add(k, v)
	return a
}
func (a *Api[T]) AddParamArray(k string, vArray []string) *Api[T] {
	for _, v := range vArray {
		a.apiParams.Add(k, v)
	}
	return a
}
func (a *Api[T]) DelParam(k string) *Api[T] {
	a.apiParams.Del(k)
	return a
}
func (a *Api[T]) AddHeader(k, v string) *Api[T] {
	a.apiHeaders[k] = v
	return a
}
func (a *Api[T]) DelHeader(k string) *Api[T] {
	delete(a.apiHeaders, k)
	return a
}
func (a *Api[T]) SetBody(body []byte) *Api[T] {
	a.apiBody = body
	return a
}

// Do 执行请求并返回响应体
func (a *Api[T]) Do(b *Bot) (*T, error) {
	if b == nil {
		var err error
		if b, err = NewBot(""); err != nil {
			return nil, err
		}
	}
	queryString := getQueryParams(a.apiParams)
	req, err := http.NewRequest(a.apiMethod, a.apiUrl+queryString, bytes.NewReader(a.apiBody))
	if err != nil {
		return nil, err
	}
	for k, v := range a.apiHeaders {
		req.Header.Set(k, v)
	}
	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(resp.Body)
	if b.Debug {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		fmt.Printf("响应体：%s\n", string(respBody))
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}
	//检查状态码是否为2xx
	if resp.StatusCode/100*100 != http.StatusOK {
		// 处理非 200 的情况
		all, _ := io.ReadAll(resp.Body)
		log.Printf("status code: %d, body: %s\n", resp.StatusCode, string(all))
		// 重新设置Body以便解码器可以再次读取
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(all))
	}
	r := new(T)
	if err = json.NewDecoder(resp.Body).Decode(r); err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	return r, nil
}

// DoRaw 执行请求并返回原始HTTP响应，不进行JSON解码
// 主要用于下载文件等需要处理原始数据的场景
func (a *Api[T]) DoRaw(b *Bot) (*http.Response, error) {
	if b == nil {
		var err error
		if b, err = NewBot(""); err != nil {
			return nil, err
		}
	}
	queryString := getQueryParams(a.apiParams)
	req, err := http.NewRequest(a.apiMethod, a.apiUrl+queryString, bytes.NewReader(a.apiBody))
	if err != nil {
		return nil, err
	}
	for k, v := range a.apiHeaders {
		req.Header.Set(k, v)
	}

	resp, err := b.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//检查状态码是否为2xx
	if resp.StatusCode/100*100 != http.StatusOK {
		// 处理非 200 的情况
		all, _ := io.ReadAll(resp.Body)
		log.Printf("status code: %d, body: %s\n", resp.StatusCode, string(all))
		// 重新设置Body以便调用者可以继续使用
		resp.Body = io.NopCloser(bytes.NewReader(all))
	}
	return resp, nil
}

// getQueryParams 获取query字符串
func getQueryParams(params url.Values) string {
	query := params.Encode()
	if query != "" {
		query = "?" + query
	}
	return query
}
