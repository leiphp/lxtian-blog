package dohttp

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// FormData 用于处理 application/x-www-form-urlencoded 格式的数据
type FormData struct {
	data map[string]string
}

// NewFormData 创建一个新的 FormData 实例
func NewFormData() *FormData {
	return &FormData{
		data: make(map[string]string),
	}
}

// Set 设置键值对
func (f *FormData) Set(key, value string) {
	f.data[key] = value
}

// Encode 将数据编码为 application/x-www-form-urlencoded 格式
func (f *FormData) Encode() string {
	var buf strings.Builder
	for key, value := range f.data {
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(url.QueryEscape(key) + "=" + url.QueryEscape(value))
	}
	return buf.String()
}

// Get 发送 GET 请求
func Get(requestURL string, headers map[string]string, queryParams map[string]string) (*http.Response, error) {
	client := &http.Client{}

	// 构建查询参数并添加到 URL
	if len(queryParams) > 0 {
		params := url.Values{}
		for key, value := range queryParams {
			params.Add(key, value)
		}
		requestURL += "?" + params.Encode()
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return client.Do(req)
}

// Post 发送 POST 请求，支持 form-data, raw, 和 application/x-www-form-urlencoded
func Post(requestURL string, headers map[string]string, body interface{}, contentType string) (*http.Response, error) {
	client := &http.Client{}
	var requestBody io.Reader

	switch contentType {
	case "application/json":
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		requestBody = bytes.NewBuffer(jsonData)

	case "application/x-www-form-urlencoded":
		data := url.Values{}
		for key, value := range body.(map[string]string) {
			data.Set(key, value)
		}
		requestBody = strings.NewReader(data.Encode())

	case "multipart/form-data":
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)
		for key, value := range body.(map[string]string) {
			_ = writer.WriteField(key, value)
		}
		// 如果 body 包含文件，你可以使用 writer.CreateFormFile 来添加文件
		// fileWriter, err := writer.CreateFormFile("file", "filename.txt")
		// if err != nil {
		//     return nil, err
		// }
		// _, _ = io.Copy(fileWriter, file)

		writer.Close()
		requestBody = buffer
		headers["Content-Type"] = writer.FormDataContentType()

	default:
		// 默认处理为 raw 数据
		requestBody = bytes.NewBufferString(body.(string))
	}

	req, err := http.NewRequest("POST", requestURL, requestBody)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	if contentType != "multipart/form-data" {
		req.Header.Set("Content-Type", contentType)
	}

	return client.Do(req)
}
