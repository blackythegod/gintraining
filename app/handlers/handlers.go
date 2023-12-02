package handlers

import (
	"gintraining/database"
	"gintraining/models"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type IHandler interface {
	InitRoutes() *gin.Engine
	SignIn(c *gin.Context)
	SignUp(c *gin.Context)
	CheckMe(c *gin.Context)
	ReadDB(c *gin.Context)
	FindUser(c *gin.Context)
}
type Handler struct {
	db *database.Database
}

func (h Handler) ReadDB(c *gin.Context) {
	var users []*models.User
	h.db.Find(&users)
	c.JSON(200, &users)
}
func (h Handler) FindUser(c *gin.Context) {
	var user *models.User
	c.BindJSON(&user)
	h.db.Find(&user, "username = ?", user.Username)
}
func (h Handler) RegisterHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user *models.User
		err := c.BindJSON(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"content": "Failed to parse params"})
		}
		if h.db.CheckCredsForExisting(user) {
			h.db.Create(&user)
			c.JSON(200, &user)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"content": "Username is already in use"})
		}
	}
}
func (h Handler) SignUp(c *gin.Context) {
	var user *models.User
	err := c.BindJSON(&user)
	if err != nil {
		log.Printf("couldn't bind json: %s", err)
	}
	log.Println(user)
	h.db.Create(&user)
	c.JSON(200, &user)
}

func (h Handler) CheckMe(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(jwt.IdentityKey)
	c.JSON(200, gin.H{
		"userid":   claims[jwt.IdentityKey],
		"username": user.(*models.Login).Username,
		"text":     "Hello World.",
	})
}

func (h Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	authMiddleware, err := jwt.New(NewAuthMiddleware(h.db.DB, true))
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}
	db := router.Group("/db")
	{
		db.GET("/list", h.ReadDB)
		db.GET("/user", h.FindUser)
	}
	auth := router.Group("/auth")
	auth.POST("/sign-up", h.RegisterHandler())
	auth.POST("/sign-in", authMiddleware.LoginHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/me", h.CheckMe)
	}
	return router
}
func InitHandler(db *database.Database) *Handler {
	return &Handler{
		db,
	}
}
