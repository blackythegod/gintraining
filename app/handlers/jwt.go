package handlers

import (
	"bytes"
	"encoding/json"
	"gintraining/models"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"time"
)

func NewAuthMiddleware(db *gorm.DB, SendCoockie bool) *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		SendCookie:  SendCoockie,
		CookieName:  "jwt",
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		Key:         []byte("secret key"),
		TimeFunc:    time.Now,
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,

		Authenticator: authFunc(db),
		PayloadFunc:   payloadFunc,
		LoginResponse: loginResponseFunc(db),

		Authorizator:    authorizatorFunc,
		Unauthorized:    unauthriziedFunc,
		IdentityHandler: identityFunc,
		LogoutResponse:  logoutResponseFunc,
	}
}
func authFunc(db *gorm.DB) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var login *models.Login
		b, _ := io.ReadAll(c.Request.Body)
		err := json.Unmarshal(b, &login)
		if err != nil {
			log.Print(err)
		}
		err = db.
			Where("username = ? AND password = ?", login.Username, login.Password).
			First(&models.User{}).
			Error
		if err != nil {
			return "", jwt.ErrFailedAuthentication
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))
		return &models.Login{
			Username: login.Username,
		}, nil
	}
}
func loginResponseFunc(db *gorm.DB) func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, message string, time time.Time) {
		var user *models.User
		err := c.BindJSON(&user)
		if err != nil {
			log.Printf("error: %s", err)
		}
		c.Writer.Header().Add("Access-Token", message)
		c.Writer.Header().Add("Expire-Token", time.Format("2006-01-02 15:04:05"))
		c.JSON(code, "logged in")
	}
}
func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*models.Login); ok {
		return jwt.MapClaims{
			jwt.IdentityKey: v.Username,
		}
	}
	return jwt.MapClaims{}
}
func identityFunc(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &models.Login{
		Username: claims[jwt.IdentityKey].(string),
	}
}
func authorizatorFunc(data interface{}, c *gin.Context) bool {
	return true
}
func unauthriziedFunc(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
func logoutResponseFunc(c *gin.Context, code int) {
	c.JSON(code, "logged out")
}
