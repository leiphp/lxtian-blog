package utils

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zeromicro/go-zero/core/logc"
	"reflect"
	"strings"
)

// StructSliceToMapSlice 将任何结构体或结构体指针切片转换为 []map[string]interface{}
func StructSliceToMapSlice(slice interface{}) ([]map[string]interface{}, error) {
	// 通过反射获取传入切片的值
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice {
		err := errors.New("input is not a slice")
		logc.Errorf(context.Background(), "Error: %v", err)
		return nil, err
	}

	// 初始化结果切片
	result := make([]map[string]interface{}, 0, v.Len())

	// 遍历每个结构体或结构体指针
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)

		// 如果是指针，需要解引用
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}

		// 检查是否为结构体
		if item.Kind() != reflect.Struct {
			err := errors.New("slice elements are not structs or struct pointers")
			logc.Errorf(context.Background(), "Error: %v", err)
			return nil, err
		}

		// 获取结构体的类型
		itemType := item.Type()

		// 创建 map 来存储结构体字段和值
		m := make(map[string]interface{})
		for j := 0; j < item.NumField(); j++ {
			field := item.Field(j)
			fieldType := itemType.Field(j)

			// 只处理导出的字段（大写字母开头的字段）
			if fieldType.PkgPath != "" {
				continue // 未导出的字段，跳过
			}

			// 将字段名转换为小写
			fieldName := strings.ToLower(fieldType.Name)
			m[fieldName] = field.Interface()
		}
		result = append(result, m)
	}

	return result, nil
}

// StructSliceToMapSliceUsingJSON 用json将任何结构体或结构体指针切片转换为 []map[string]interface{
func StructSliceToMapSliceUsingJSON(slice interface{}) ([]map[string]interface{}, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil, err
	}

	return result, nil
}
