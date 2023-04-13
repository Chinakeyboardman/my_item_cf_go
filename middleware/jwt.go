package middleware

//middleware/middleware.go

import (
	"errors"
	"fmt"
	"my_item_cf_go/context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	jwt "github.com/dgrijalva/jwt-go"
)

// func M1(c *context.Context) {

// 	fmt.Println("我是1")

// }

var (
	TokenExpired     error = errors.New("Token is expired")
	TokenNotValidYet error = errors.New("Token is not active yet")
	TokenMalformed   error = errors.New("Malformed token")
	TokenInvalid     error = errors.New("can't handle this token")
)

const (
	// 这个是需要保密的一段信息
	SignKey         string = "a87x80wfebei90f8532f16f423b125616dea9b75"
	Gin_Context_Key string = "claims"
)

func JWTAuth(c *context.Context) {

	fmt.Println("JWTAuth")

	//过滤是否验证token， login结构直接放行，这里为了简单起见，直接判断路径中是否带login，携带login直接放行
	if strings.Contains(c.Request.RequestURI, "login") {
		return
	}

	token := c.Request.Header.Get("token")
	if token == "" {
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    "请求未携带token，无权限访问",
		})
		c.Abort()
		return
	}

	// k8s日志
	// klog.Infof("gotten token:%s", token)

	j := NewJWT()
	// parse token, get the user and role info
	claims, err := j.ParseToken(token)
	if err != nil {
		if err == TokenExpired {
			c.JSON(http.StatusOK, gin.H{
				"status": -1,
				"msg":    "授权已过期",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": -1,
			"msg":    err.Error(),
		})
		c.Abort()
		return
	}
	// 继续交由下一个路由处理,并将解析出的信息传递下去
	c.Set(Gin_Context_Key, claims)
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 载荷，可添加自己需要的一些信息
type CustomClaims struct {
	UserId   int64  `json:"userId"`
	UserName string `json:"userName"`
	RoleId   string `json:"roleId"`
	jwt.StandardClaims
}

// 新建一个jwt实例
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// 获取signKey
func GetSignKey() string {
	return SignKey
}

// CreateToken 生成一个token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, TokenInvalid
}

// 更新token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return j.CreateToken(*claims)
	}
	return "", TokenInvalid
}
