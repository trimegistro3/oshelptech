package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"}) })

	// Rotas públicas: cadastro e login não exigem token.
	r.POST("/auth/register", proxy(env("AUTH_SERVICE_URL", "http://auth-service:8001")))
	r.POST("/auth/login", proxy(env("AUTH_SERVICE_URL", "http://auth-service:8001")))

	protected := r.Group("/")
	protected.Use(jwtMiddleware())
	{
		protected.Any("/clientes", proxy(env("CLIENTE_SERVICE_URL", "http://cliente-service:8002")))
		protected.Any("/clientes/:id", proxy(env("CLIENTE_SERVICE_URL", "http://cliente-service:8002")))
		protected.Any("/ordens", proxy(env("ORDEM_SERVICE_URL", "http://ordem-service:8003")))
		protected.Any("/ordens/:id", proxy(env("ORDEM_SERVICE_URL", "http://ordem-service:8003")))
		protected.Any("/pagamentos", proxy(env("FINANCEIRO_SERVICE_URL", "http://financeiro-service:8004")))
		protected.Any("/pagamentos/:id", proxy(env("FINANCEIRO_SERVICE_URL", "http://financeiro-service:8004")))
	}

	r.Run(":8080")
}

func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "envie Authorization: Bearer TOKEN"})
			return
		}
		tokenString := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(env("JWT_SECRET", "helptech_secret_dev")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
			return
		}
		c.Next()
	}
}

func proxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		url := target + c.Request.URL.RequestURI()
		req, err := http.NewRequest(c.Request.Method, url, bytes.NewReader(body))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao criar proxy"})
			return
		}
		req.Header = c.Request.Header.Clone()
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "serviço interno indisponível"})
			return
		}
		defer resp.Body.Close()
		for k, values := range resp.Header {
			for _, v := range values {
				c.Writer.Header().Add(k, v)
			}
		}
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
