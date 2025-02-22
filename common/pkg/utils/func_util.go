package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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
