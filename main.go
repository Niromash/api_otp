package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"sync"
	"time"
)

var letterRunes = []rune("123456789")

func main() {
	codes := make(map[string]string)
	var mu sync.Mutex

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true

	router.Use(cors.New(config))

	router.POST("/otp/api/login", func(c *gin.Context) {
		var body struct {
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()

		codes[body.Email] = randomCode(6)
		c.JSON(200, gin.H{
			"message": codes[body.Email],
			"status":  200,
		})
	})

	router.POST("/otp/api/code", func(c *gin.Context) {
		var body struct {
			Code  string `json:"code"`
			Email string `json:"email"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if storedCode, ok := codes[body.Email]; ok && storedCode == body.Code {
			delete(codes, body.Email)
			c.JSON(200, gin.H{
				"message": "granted",
				"status":  200,
			})
		} else {
			c.JSON(200, gin.H{
				"message": "not allowed",
				"status":  200,
			})
		}
	})

	log.Println("Server started on port 8080")
	log.Fatalln(router.Run(":8080"))
}

func randomCode(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
