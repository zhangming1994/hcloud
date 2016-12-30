package utils

import (
	"encoding/base64"
	"strings"
)

func MainEncrypt(endata string) string {

	var data = Base64encodeUrl(base64.StdEncoding.EncodeToString([]byte(endata)))
	return data
}

func MainDecrypt(dedata string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(Base64decodeUrl(dedata))
	return result, err
}

/// <summary>
/// 从普通字符串转换为适用于URL的Base64编码字符串
/// </summary>
func Base64decodeUrl(base64String string) string {
	str := strings.Replace(base64String, "-", "+", -1)
	return str
}

/// <summary>
/// 从普通字符串转换为适用于URL的Base64编码字符串
/// </summary>
func Base64encodeUrl(base64String string) string {
	str := strings.Replace(base64String, "+", "-", -1)
	return str
}
