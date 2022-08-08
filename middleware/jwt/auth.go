package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const ErrMsgKey = "errMsg"
const DataKey = "data"

// JWTAuth 设置路由的中间件 这个是JWT校验的中间件
func JWTAuth() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		// 获取token 前端传过来的token存放在Authorization字段中
		token := ctx.GetHeader("Authorization")
		if token == "" {
			// 没有授权 401
			ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Not Authorized."})
			ctx.Abort()
			return
		}

		// 创建JWT实例
		j := NewJWT()

		// 解析token 得到claims
		claims, err := j.ParseToken(token)
		if err != nil {
			// 过期
			if err == TokenExpired {
				ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: "Authorization has expired."})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusUnauthorized, gin.H{ErrMsgKey: err.Error()})
			ctx.Abort()
			return
		}
		// 继续交由下一个路由处理,并将解析出的信息传递下去
		ctx.Set("claims", claims)
	}
}

// JWT 签名结构
type JWT struct {
	SigningKey []byte
}

// 一些常量
var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "Our Seckill Secret Key" // 用来签名的密钥
	Issuer 			 string = "this is a issuer"       // 发行人名字
)

// CustomClaims 载荷，可以加一些自己需要的信息
type CustomClaims struct {
	Username  string `json:"username"`
	Password string `json:"password"`
	Kind     string `json:"kind"`
	jwt.StandardClaims
}

// NewJWT 新建一个jwt实例 包含的内容是用于加密解密的对称密钥 HMAC算法
func NewJWT() *JWT {
	return &JWT{
		[]byte(GetSignKey()),
	}
}

// GetSignKey 获取signKey
func GetSignKey() string {
	return SignKey
}

// SetSignKey 这是SignKey
func SetSignKey(key string) string {
	SignKey = key
	return SignKey
}

// CreateToken 生成一个token token采用HS256
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回对称密钥
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

func (j *JWT) RefreshToken(tokenString string) (string, error) {
	// 这样子的话不是肯定不会过期吗？
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






