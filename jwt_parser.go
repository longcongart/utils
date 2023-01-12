package utils

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      uint64 //uuid.UUID
	Credentials map[string]bool
	Expires     int64
}

// ExtractTokenMetadata func 从JWT中提取令牌元数据.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// 解析公共部分
		// User ID.
		userID := claims["uid"].(uint64) //采用snowFlake雪花算法 放弃uuid.Parse(claims["id"].(string))
		if userID <= 0 {
			return nil, err
		}
		// Expires time.
		expires := int64(claims["expires"].(float64))

		// User credentials.
		credentials := map[string]bool{
			// "book:create": claims["book:create"].(bool),
			// "book:update": claims["book:update"].(bool),
			// "book:delete": claims["book:delete"].(bool),
		}
		// 创建实例
		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}

	return nil, err
}

// 提取令牌
func extractToken(c *fiber.Ctx) string {
	bearToken := c.Get("Authorization") //提取 Authorization 参数值
	// 正常的授权HTTP标头部分.
	onlyToken := strings.Split(bearToken, " ") //对授权令牌分割字符串
	if len(onlyToken) == 2 {                   //判断
		return onlyToken[1] //取第二部分
	}
	return ""
}

// 验证令牌
func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)
	// 对令牌进行解析
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
