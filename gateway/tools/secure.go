package tools

import (
	"fmt"
	"gateway/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/martian/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const secret = "go-circle"

func GenerateToken(user model.User) (string, error) {
	user.Issuer = "go-business"
	user.NotBefore = time.Now().Unix() - 60
	user.ExpiresAt = time.Now().Unix() + 60*60*2
	return jwt.NewWithClaims(jwt.SigningMethodHS256, user).SignedString([]byte(secret))

}

func HashAndSalt(pwdStr string) (pwdHash []byte, err error) {
	pwdHash, err = bcrypt.GenerateFromPassword([]byte(pwdStr), bcrypt.MinCost)
	if err != nil {
		log.Errorf("加密算法发生错误", err)
	}

	return

}

func ComparePasswords(pwdHash string, pwdPlain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(pwdPlain))
	if err != nil {
		return false
	}
	return true
}

func ParserToken(tokenString string) (*model.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.User{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		log.Errorf(err.Error())
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token不可用")
				// ValidationErrorExpired表示Token过期
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token过期")
				// ValidationErrorNotValidYet表示无效token
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("无效的token")
			} else {
				return nil, fmt.Errorf("token不可用")
			}
		}
	}
	// 将token中的claims信息解析出来并断言成用户自定义的有效载荷结构
	if user, ok := token.Claims.(*model.User); ok && token.Valid {
		return user, nil
	}
	return nil, fmt.Errorf("token无效")
}
