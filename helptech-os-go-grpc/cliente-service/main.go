package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Cliente struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Nome      string    `json:"nome"`
	Telefone  string    `json:"telefone"`
	Email     string    `json:"email"`
	CpfCnpj   string    `json:"cpf_cnpj" gorm:"column:cpf_cnpj"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Cliente) TableName() string { return "clientes" }

var db *gorm.DB

func main() {
	db = connectDB()
	// AutoMigrate cria/atualiza a tabela deste serviço automaticamente no PostgreSQL.
	if err := db.AutoMigrate(&Cliente{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "cliente-service"}) })
	r.POST("/clientes", create)
	r.GET("/clientes", list)
	r.GET("/clientes/:id", getByID)
	r.PUT("/clientes/:id", update)
	r.DELETE("/clientes/:id", remove)
	log.Println("cliente-service rodando na porta 8002")
	log.Fatal(r.Run(":8002"))
}

func connectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Sao_Paulo", env("POSTGRES_HOST", "postgres"), env("POSTGRES_USER", "helptech"), env("POSTGRES_PASSWORD", "helptech123"), env("POSTGRES_DB", "helptech_os"), env("POSTGRES_PORT", "5432"))
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("erro ao conectar no banco: ", err)
	}
	return database
}

func create(c *gin.Context) {
	var item Cliente
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	db.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func list(c *gin.Context) {
	var items []Cliente
	db.Find(&items)
	c.JSON(http.StatusOK, items)
}

func getByID(c *gin.Context) {
	var item Cliente
	if err := db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro não encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func update(c *gin.Context) {
	var item Cliente
	if err := db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro não encontrado"})
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	db.Save(&item)
	c.JSON(http.StatusOK, item)
}

func remove(c *gin.Context) {
	if err := db.Delete(&Cliente{}, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao remover"})
		return
	}
	c.Status(http.StatusNoContent)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
