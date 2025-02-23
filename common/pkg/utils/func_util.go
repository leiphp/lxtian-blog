package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"lxtian-blog/common/pkg/define"
	"net/http"
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
