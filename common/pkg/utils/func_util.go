package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/pkg/define"
	"math/rand"
	"net/http"
	"time"
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
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Liliana", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Amaya", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Jameson", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Brian", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Nolan", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Ryan", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=George", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Andrea", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Adrian", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Easton", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Liam", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Aidan", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Eliza", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Wyatt", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Christopher", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Jade", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Sawyer", Gender: 1},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Jessica", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Jocelyn", Gender: 2},
		{URL: "https://api.dicebear.com/9.x/avataaars/svg?seed=Aiden", Gender: 2},
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
