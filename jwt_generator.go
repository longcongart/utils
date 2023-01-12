package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Tokens struct to describe tokens object.
type Tokens struct {
	Access  string
	Refresh string
}

// GenerateNewTokens函数 用于生成新的访问和刷新令牌
func GenerateNewTokens(uid int64, credentials []string) (*Tokens, error) {
	// 生成JWT访问令牌.
	accessToken, err := generateNewAccessToken(uid, credentials)
	if err != nil {
		// 令牌生成返回错误
		return nil, err
	}
	// 生成JWT刷新令牌
	refreshToken, err := generateNewRefreshToken()
	if err != nil {
		// 令牌生成返回错误
		return nil, err
	}
	// 返回创建实例
	return &Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}

// 生成新的访问令牌
func generateNewAccessToken(uid int64, credentials []string) (string, error) {
	// 从.env文件中读取密钥key设置 JWT_SECRET_KEY
	secret := os.Getenv("JWT_SECRET_KEY")
	// 从.env文件中读取到期分钟计数设置 JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT
	minutesCount, _ := strconv.Atoi(os.Getenv("JWT_SECRET_KEY_EXPIRE_MINUTES_COUNT"))
	// 创建新声明断言字典
	claims := jwt.MapClaims{}
	// 设置公共声明:
	claims["uid"] = uid
	claims["expires"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()
	// claims["book:create"] = false
	// claims["book:update"] = false
	// claims["book:delete"] = false
	// 继续设置私有专用令牌凭据:
	for _, credential := range credentials {
		claims[credential] = true
	}
	// 创建带有声明的新JWT访问令牌 将声明部分进行she256哈希方法加密
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 生成令牌.
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		// 生成JWT令牌失败 返回错误
		return "", err
	}
	return t, nil
}

// 生成新的刷新令牌
func generateNewRefreshToken() (string, error) {
	// 创建新的SHA256哈希
	sha256 := sha256.New()
	// 加盐字符串 JWT_REFRESH_KEY+现在的日期时间
	refresh := os.Getenv("JWT_REFRESH_KEY") + time.Now().String()
	_, err := sha256.Write([]byte(refresh)) // See: https://pkg.go.dev/io#Writer.Write
	if err != nil {
		// 它刷新令牌生成失败返回错误
		return "", err
	}
	// 从.env文件中读取刷新key的多少小时到期计数设置 JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))
	// 设置到期时间
	expireTime := fmt.Sprint(time.Now().Add(time.Hour * time.Duration(hoursCount)).Unix())
	// 创新新的刷新令牌(sha256 string with salt + expire time).
	t := hex.EncodeToString(sha256.Sum(nil)) + "." + expireTime
	return t, nil
}

// 解析刷新令牌函数 对刷新令牌的第二参数进行解析
func ParseRefreshToken(refreshToken string) (int64, error) {
	return strconv.ParseInt(strings.Split(refreshToken, ".")[1], 0, 64)
}
