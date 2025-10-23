package utils

import "regexp"

// IsValidWeChatId 检查微信号是否有效
func IsValidWeChatId(wechatId string) bool {
	// 微信号规则正则：以字母开头，长度 6-20，只能包含字母、数字、下划线、减号
	pattern := `^[a-zA-Z][a-zA-Z0-9_-]{5,19}$`
	matched, err := regexp.MatchString(pattern, wechatId)
	if err != nil {
		return false
	}
	return matched
}

// IsValidChinesePhone 检查手机号是否有效
func IsValidChinesePhone(phone string) bool {
	// 手机号正则：以 1 开头，第二位是 3-9，总长 11 位
	pattern := `^1[3-9]\d{9}$`
	matched, err := regexp.MatchString(pattern, phone)
	if err != nil {
		return false
	}
	return matched
}
