package utils

// 辅助方法：从 map 中获取 bool 值
func GetBoolValue(m map[string]interface{}, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// 辅助方法：从 map 中获取 int32 值
func GetInt32Value(m map[string]interface{}, key string) int32 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int32:
			return val
		case int:
			return int32(val)
		case int64:
			return int32(val)
		case float64:
			return int32(val)
		}
	}
	return 0
}

// 辅助方法：从 map 中获取 int64 值
func GetInt64Value(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case int64:
			return val
		case int:
			return int64(val)
		case int32:
			return int64(val)
		case float64:
			return int64(val)
		}
	}
	return 0
}

// 辅助方法：从 map 中获取 string 值
func GetStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
