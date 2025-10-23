package utils

import (
	"fmt"
	"strings"
)

// GroupBy 按照指定函数分组
func GroupBy[T any, R comparable](arr []T, fn func(T) R) map[R][]T {
	m := make(map[R][]T)
	for _, v := range arr {
		key := fn(v)
		m[key] = append(m[key], v)
	}
	return m
}

// ToMap 按照指定函数生成map
func ToMap[T any, R comparable](arr []T, fn func(T) R) map[R]T {
	m := make(map[R]T)
	for _, v := range arr {
		key := fn(v)
		m[key] = v
	}
	return m
}

// MergeMap 合并多个map,返回一个新的map(若存在key冲突，则后面的覆盖前面的)
func MergeMap[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// SliceFilter 过滤切片
func SliceFilter[T any](arr []T, fn func(T) bool) []T {
	result := make([]T, 0, len(arr))
	for _, v := range arr {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

// IntListToString 将int列表转换为字符串
func IntListToString[T int | int8 | int16 | int32 | int64](intList []T) string {
	var builder strings.Builder
	for i, num := range intList {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%d", num))
	}
	return builder.String()
}

// ExtractField 提取切片中的某一字段组成新切片
func ExtractField[T, R any](list []T, fn func(T) R) []R {
	res := make([]R, len(list))
	for i, num := range list {
		res[i] = fn(num)
	}
	return res
}
