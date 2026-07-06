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

type Ordem struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	ClienteID   uint      `json:"cliente_id"`
	Equipamento string    `json:"equipamento"`
	Marca       string    `json:"marca"`
	Modelo      string    `json:"modelo"`
	Defeito     string    `json:"defeito"`
	Diagnostico string    `json:"diagnostico"`
	Status      string    `json:"status"`
	Valor       float64   `json:"valor"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Ordem) TableName() string { return "ordens" }

var db *gorm.DB

func main() {
	db = connectDB()
	// AutoMigrate cria/atualiza a tabela deste serviço automaticamente no PostgreSQL.
	if err := db.AutoMigrate(&Ordem{}); err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "ordem-service"}) })
	r.POST("/ordens", create)
	r.GET("/ordens", list)
	r.GET("/ordens/:id", getByID)
	r.PUT("/ordens/:id", update)
	r.DELETE("/ordens/:id", remove)
	log.Println("ordem-service rodando na porta 8003")
	log.Fatal(r.Run(":8003"))
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
	var item Ordem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "json inválido"})
		return
	}
	db.Create(&item)
	c.JSON(http.StatusCreated, item)
}

func list(c *gin.Context) {
	var items []Ordem
	db.Find(&items)
	c.JSON(http.StatusOK, items)
}

func getByID(c *gin.Context) {
	var item Ordem
	if err := db.First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "registro não encontrado"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func update(c *gin.Context) {
	var item Ordem
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
	if err := db.Delete(&Ordem{}, c.Param("id")).Error; err != nil {
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
