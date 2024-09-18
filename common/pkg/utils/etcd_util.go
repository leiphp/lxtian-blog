package utils

import "strings"

// ParseHosts 接受一个包含多个主机地址的字符串（逗号分隔），并返回一个解析后的主机地址列表
func ParseHosts(hosts string) []string {
	var result []string
	// 拆分主机地址字符串为多个主机地址
	if hosts != "" {
		splitHosts := strings.Split(hosts, ",")
		result = append(result, splitHosts...)
	} else {
		result = append(result, "127.0.0.1:2379")
	}
	return result
}
