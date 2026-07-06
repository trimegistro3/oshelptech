package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Nome         string    `json:"nome"`
	Email        string    `json:"email" gorm:"uniqueIndex"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type registerRequest struct{ Nome, Email, Password string }
type loginRequest struct{ Email, Password string }

func (User) TableName() string { return "users" }

var db *gorm.DB

func main() {
	db = connectDB()
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "auth-service"}) })
	r.POST("/auth/register", register)
	r.POST("/auth/login", login)
	log.Println("auth-service rodando na porta 8001")
	log.Fatal(r.Run(":8001"))
}

func connectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo", env("POSTGRES_HOST", "postgres"), env("POSTGRES_USER", "helptech"), env("POSTGRES_PASSWORD", "helptech123"), env("POSTGRES_DB", "helptech_os"), env("POSTGRES_PORT", "5432"))
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("erro ao conectar no banco: ", err)
	}
	return database
}

func register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Nome == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nome, email e password são obrigatórios"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar senha"})
		return
	}
	user := User{Nome: req.Nome, Email: req.Email, PasswordHash: string(hash)}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "não foi possível criar usuário; verifique se o email já existe"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": user.ID, "nome": user.Nome, "email": user.Email})
}

func login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	var user User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciais inválidas"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID, "email": user.Email, "exp": time.Now().Add(24 * time.Hour).Unix()})
	signed, err := token.SignedString([]byte(env("JWT_SECRET", "helptech_secret_dev")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed, "type": "Bearer"})
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
