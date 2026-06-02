package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Estructura de Producto alineada con la Base de Datos
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
}

var db *sql.DB

// Middleware para habilitar CORS (Crucial para que Svelte se pueda conectar)
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRifier, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func main() {
	// 1. Lectura de Variables de Entorno (Principio Stateless / Cloud-Native)
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "storeuser")
	dbPassword := getEnv("DB_PASSWORD", "storepassword")
	dbName := getEnv("DB_NAME", "nonamestore")

	// 2. Configurar la cadena de conexión (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Error configurando el driver de la BD: %v", err)
	}
	defer db.Close()

	// 3. Inicializar Gin Gonic
	r := gin.Default()
	r.Use(CORSMiddleware())

	// Endpoint 1: Health Check (Liveness/Readiness Probe para Kubernetes)
	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "DOWN",
				"database": "UNREACHABLE",
				"error":    err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "UP",
			"database": "CONNECTED",
		})
	})

	// Endpoint 2: Listar todos los productos
	r.GET("/products", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, description, price, image_url FROM products")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			products = append(products, p)
		}

		c.JSON(http.StatusOK, products)
	})

	// Endpoint 3: Detalle de un producto específico por ID
	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p Product

		err := db.QueryRow("SELECT id, name, description, price, image_url FROM products WHERE id = ?", id).
			Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL)

		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Producto no encontrado"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, p)
	})

	log.Println("🚀 Servidor Backend corriendo en http://localhost:8080")
	r.Run(":8080")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
