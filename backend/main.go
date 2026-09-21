package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {

	databaseURL := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://kt-crm-portfolio.tajimak1106.workers.dev"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// ヘルスチェック
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 顧客一覧
	r.GET("/api/customers", func(c *gin.Context) {
		rows, err := pool.Query(
			context.Background(),
			`SELECT id, name, email, company, created_at
			 FROM customers
			 ORDER BY id`,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		defer rows.Close()

		customers := []Customer{}

		for rows.Next() {
			var customer Customer

			err := rows.Scan(
				&customer.ID,
				&customer.Name,
				&customer.Email,
				&customer.Company,
				&customer.CreatedAt,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}
			customers = append(customers, customer)
		}

		c.JSON(http.StatusOK, customers)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
