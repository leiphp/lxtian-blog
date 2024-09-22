package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

var key = []byte("0123456789abcdef")

// 加密
func Encrypt(text []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 填充明文
	text = PKCS7Padding(text, block.BlockSize())

	// 创建一个使用 AES 算法的 CBC 模式的加密器
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], text)

	return ciphertext, nil
}

// 解密
func Decrypt(ciphertext []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 提取 IV
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 创建一个使用 AES 算法的 CBC 模式的解密器
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去掉填充
	ciphertext = PKCS7Unpadding(ciphertext)

	return string(ciphertext), nil
}

// PKCS7Padding 在明文末尾填充 PKCS7
func PKCS7Padding(text []byte, blockSize int) []byte {
	padding := blockSize - len(text)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(text, padText...)
}

// PKCS7Unpadding 去掉末尾的 PKCS7 填充
func PKCS7Unpadding(text []byte) []byte {
	length := len(text)
	unpadding := int(text[length-1])
	return text[:(length - unpadding)]
}
