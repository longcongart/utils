package utils

// AES CBC库+Base64库实现加解密
// 本文件节选自网络: http://www.zzvips.com/article/148213.html
/* 示例
func main() {
	var key = "321423u9y8d2fwfl"
	var iv = "1!34c6$9~76W43D1"
	encodeStr := "我勒个去"
	es, _ := AesCbcEncrypt(encodeStr, []byte(key), []byte(iv))
	fmt.Println(es)
	ds, _ := AesCbcDecrypt(es, []byte(key), []byte(iv))
	fmt.Println(string(ds))
}
*/

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func AesCbcEncrypt(encodeStr string, key []byte, iv []byte) (string, error) {
	encodeBytes := []byte(encodeStr)
	//根据key 生成密文
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	encodeBytes = PKCS5Padding(encodeBytes, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)
	return base64.StdEncoding.EncodeToString(crypted), nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//填充
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func AesCbcDecrypt(decodeStr string, key []byte, iv []byte) ([]byte, error) {
	//先解密base64
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
