package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/pkg/define"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/shopspring/decimal"
)

// 用户服务：通过 HTTP 调用聊天服务接口
func SendMessageToChatService(host string, port int, userID, message string) error {
	url := fmt.Sprintf("http://%s:%d/send_msg", host, port)
	payload := map[string]string{
		"user_id": userID,
		"message": message,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message: %s", resp.Status)
	}
	return nil
}

func GetSocketMessage(token, msg string, userinfo define.User) (string, error) {
	// 初始化结构体数据
	loginResponse := define.LoginResponse{
		Type:     "login",
		Status:   "success",
		Msg:      msg,
		Token:    token,
		UserInfo: userinfo,
	}
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(loginResponse)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return "", err
	}
	// 打印 JSON 字符串
	jsonString := string(jsonData)
	fmt.Println(jsonString)
	return jsonString, nil
}

// RandomString 生成指定长度的随机字符串
func RandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	source := rand.NewSource(time.Now().UnixNano()) // 创建随机数种子源
	rng := rand.New(source)                         // 创建本地随机数生成器

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rng.Intn(len(charset))]
	}
	return string(b)
}

// GetRandomAvatar 根据性别随机获取头像
func GetRandomAvatar(gender int) (string, error) {
	// 定义 Avatar 结构体
	type Avatar struct {
		URL    string
		Gender int
	}

	// 初始化头像列表
	avatars := []Avatar{
		{URL: "https://img.100txy.com/blog/avatar/Liliana.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/Amaya.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Jameson.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Brian.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Nolan.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Ryan.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/George.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Andrea.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Adrian.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Easton.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Liam.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/Aidan.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Eliza.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/Wyatt.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Christopher.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Jade.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Sawyer.svg", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/Jessica.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/Jocelyn.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/Aiden.svg", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/1.png", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/2.png", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/3.png", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/4.png", Gender: 2},
		{URL: "https://img.100txy.com/blog/avatar/5.png", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/6.png", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/7.png", Gender: 1},
		{URL: "https://img.100txy.com/blog/avatar/8.png", Gender: 1},
	}

	// 根据性别过滤头像
	var filteredAvatars []Avatar
	for _, avatar := range avatars {
		if avatar.Gender == gender {
			filteredAvatars = append(filteredAvatars, avatar)
		}
	}

	// 如果没有符合性别的头像，返回错误
	if len(filteredAvatars) == 0 {
		return "", fmt.Errorf("没有找到符合性别的头像")
	}

	// 随机选择一个头像
	randomIndex := rand.Intn(len(filteredAvatars))
	return filteredAvatars[randomIndex].URL, nil
}

// ConvertByteFieldsToString 遍历列表中的字段，如果是 []byte 类型就转成 string
func ConvertByteFieldsToString(data []map[string]interface{}) {
	for i := range data {
		for k, v := range data[i] {
			if bv, ok := v.([]byte); ok {
				data[i][k] = string(bv)
			}
		}
	}
}

// 处理东八区字段为标准时间
// 支持 time.Time、RFC3339 字符串（如 2025-07-14T00:06:19+08:00）、常见 "2006-01-02 15:04:05" 字符串，以及 []byte。
func FormatTimeFields(data []map[string]interface{}, fields ...string) {
	const layout = "2006-01-02 15:04:05"
	for i := range data {
		for _, field := range fields {
			switch v := data[i][field].(type) {
			case time.Time:
				data[i][field] = v.Format(layout)
			case string:
				if parsed, err := parseTimeString(v); err == nil {
					data[i][field] = parsed.Format(layout)
				}
			case []byte:
				if parsed, err := parseTimeString(string(v)); err == nil {
					data[i][field] = parsed.Format(layout)
				}
			}
		}
	}
}

// parseTimeString 尝试解析常见的时间字符串格式
func parseTimeString(s string) (time.Time, error) {
	// 优先尝试 RFC3339（含时区）
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	// 兜底常用格式：yyyy-MM-dd HH:mm:ss
	return time.Parse("2006-01-02 15:04:05", s)
}

// 处理东八区字段为标准时间
func FormatTimeFieldsInMap(data map[string]interface{}, fields ...string) {
	for _, field := range fields {
		if t, ok := data[field].(time.Time); ok {
			data[field] = t.Format("2006-01-02 15:04:05")
		}
	}
}

// 批量处理字段0/1为true or flase
func FormatBoolFields(data []map[string]interface{}, fields ...string) {
	for i := range data {
		for _, field := range fields {
			if v, ok := data[i][field]; ok {
				switch val := v.(type) {
				case int:
					data[i][field] = val == 1
				case int8:
					data[i][field] = val == 1
				case uint8:
					data[i][field] = val == 1
				case int16:
					data[i][field] = val == 1
				case int32:
					data[i][field] = val == 1
				case int64:
					data[i][field] = val == 1
				case float64:
					data[i][field] = int(val) == 1
				case []byte:
					// 数据库返回 []byte 的情况，如 tinyint(1)
					if string(val) == "1" {
						data[i][field] = true
					} else {
						data[i][field] = false
					}
				default:
					if v == nil {
						data[i][field] = false
					} else {
						fmt.Printf("未匹配字段 [%s]: 类型[%T], 值=%v\n", field, v, v)
					}
				}
			}
		}
	}
}

// 生成唯一文件名
func GenerateFilename(original string) string {
	ext := path.Ext(original) // .png
	name := fmt.Sprintf("%d_%d%s", time.Now().UnixNano(), rand.Intn(1000), ext)
	prefix := time.Now().Format("200601")     // 生成 "202506"
	return fmt.Sprintf("%s/%s", prefix, name) // 返回 "202506/xxx.png"
}

type TreeNode struct {
	Id       int64      `json:"id"`
	IsGroup  int32      `json:"is_group"`
	Title    string     `json:"title"`
	Children []TreeNode `json:"children,omitempty"`
}

// 构建树结构
//func BuildTree(data []map[string]interface{}, parentId int64) []TreeNode {
//	var tree []TreeNode
//	for _, item := range data {
//		if item["parent_id"] == parentId {
//			node := TreeNode{
//				Id:      item["id"].(int64),
//				IsGroup: item["is_group"].(int32),
//				Label:   item["title"].(string),
//			}
//			children := BuildTree(data, item["id"].(int64))
//			if len(children) > 0 {
//				node.Children = children
//			}
//			tree = append(tree, node)
//		}
//	}
//	return tree
//}

func BuildTree(data []map[string]interface{}, parentId int64) []TreeNode {
	var tree []TreeNode
	for _, item := range data {
		// 安全获取 parent_id 并转换为 int64
		itemParentID, err := getInt64(item["parent_id"])
		if err != nil {
			continue // 或者处理错误
		}

		if itemParentID != parentId {
			continue
		}

		// 安全获取 id
		id, err := getInt64(item["id"])
		if err != nil {
			continue
		}

		// 安全获取 is_group
		isGroup, err := getInt32(item["is_group"])
		if err != nil {
			continue
		}

		// 安全获取 title
		title, ok := item["title"].(string)
		if !ok {
			title = ""
		}

		node := TreeNode{
			Id:      id,
			IsGroup: isGroup,
			Title:   title,
		}

		// 递归构建子树
		children := BuildTree(data, id)
		if len(children) > 0 {
			node.Children = children
		}

		tree = append(tree, node)
	}
	return tree
}

// 安全地将 interface{} 转换为 int64
func getInt64(val interface{}) (int64, error) {
	switch v := val.(type) {
	case int64:
		return v, nil
	case uint64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case uint:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}

// 安全地将 interface{} 转换为 int32
func getInt32(val interface{}) (int32, error) {
	switch v := val.(type) {
	case int32:
		return v, nil
	case uint32:
		return int32(v), nil
	case int64:
		return int32(v), nil
	case uint64:
		return int32(v), nil
	case int:
		return int32(v), nil
	case uint:
		return int32(v), nil
	case float64:
		return int32(v), nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}

func BuildTreeMap(data []map[string]interface{}, parentId int64) []map[string]interface{} {
	var tree []map[string]interface{}

	for _, item := range data {
		// 获取 parent_id
		itemParentID, err := getInt64(item["parent_id"])
		if err != nil || itemParentID != parentId {
			continue
		}

		// 获取 id
		id, err := getInt64(item["id"])
		if err != nil {
			continue
		}

		// 创建当前节点副本（避免原始数据污染）
		node := map[string]interface{}{
			"id":       item["id"],
			"title":    item["title"],
			"is_group": item["is_group"],
		}

		// 递归构建子节点
		children := BuildTreeMap(data, id)
		if len(children) > 0 {
			node["children"] = children
		}

		tree = append(tree, node)
	}

	return tree
}

// Int64ToString 将分转成字符串表示（保留整数分）
func Int64ToString(fen int64) string {
	return decimal.NewFromInt(fen).String()
}

// StringToInt64 将字符串的分转成int64
func StringToInt64(s string) (int64, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return 0, err
	}
	return d.IntPart(), nil
}
