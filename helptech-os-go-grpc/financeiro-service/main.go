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

type Pagamento struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	OrdemID        uint      `json:"ordem_id"`
	Valor          float64   `json:"valor"`
	FormaPagamento string    `json:"forma_pagamento"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (Pagamento) TableName() string { return "pagamentos" }

var db *gorm.DB

func main() {
	db = connectDB()
	// AutoMigrate cria/atualiza a tabela deste serviço automaticamente no PostgreSQL.
	if err := db.AutoMigrate(&Pagamento{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "financeiro-service"}) })
	r.POST("/pagamentos", create)
	r.GET("/pagamentos", list)
	r.GET("/pagamentos/:id", getByID)
	log.Println("financeiro-service rodando na porta 8004")
	log.Fatal(r.Run(":8004"))
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
	var item Pagamento
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	db.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func list(c *gin.Context) {
	var items []Pagamento
	db.Find(&items)
	c.JSON(http.StatusOK, items)
}

func getByID(c *gin.Context) {
	var item Pagamento
	if err := db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro não encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
