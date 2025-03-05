package utils

import (
	"gorm.io/gorm"
)

// 从数据库中随机获取一个形容词
func getRandomAdjective(db *gorm.DB, gender int) (string, error) {
	var name string
	query := "SELECT name FROM txy_nickname WHERE type = 1 and gender = ? ORDER BY RAND() LIMIT 1"
	err := db.Raw(query, gender).Debug().Scan(&name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

// 从数据库中随机获取一个名词
func getRandomNoun(db *gorm.DB, gender int) (string, error) {
	var name string
	query := "SELECT name FROM txy_nickname WHERE type = 2 and gender = ? ORDER BY RAND() LIMIT 1"
	err := db.Raw(query, gender).Debug().Scan(&name).Error
	if err != nil {
		return "", err
	}
	return name, nil
}

// 生成随机昵称
func GenerateNickname(db *gorm.DB, gender int) (string, error) {
	adjective, err := getRandomAdjective(db, gender)
	if err != nil {
		return "", err
	}
	noun, err := getRandomNoun(db, gender)
	if err != nil {
		return "", err
	}
	// 生成 4 个随机字符
	randomStr := RandomString(4)
	return adjective + noun + randomStr, nil
}
