package verify

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claim struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var MySecret = []byte("云升不是处")

/**
签发token
*/
func SignToken(userId string) (string, error) {
	c := Claim{
		userId, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 356).Unix(), // 过期时间
			Issuer:    "hita",                                      // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*Claim, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claim); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
