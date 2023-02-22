package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

func validateAPIKey() gin.HandlerFunc {
	APIKey := os.Getenv("API_KEY")

	return func(c *gin.Context) {
		if c.Request.Header.Get("X-API-KEY") != APIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Authentication failed"})
			return
		}

	}
}

type Data struct {
	Style string `yaml:"style" json:"style"`
	Title struct {
		Text  string `yaml:"text" json:"text"`
		Style string `yaml:"style" json:"style"`
	}
	BackgroundColor string `yaml:"backgroundColor" json:"backgroundColor"`
	Items           []Item `yaml:"items" json:"items"`
}

type Item struct {
	Title  string `yaml:"title" json:"title"`
	Price  string `yaml:"price" json:"price"`
	Link   string `yaml:"link" json:"link"`
	Rating int64  `yaml:"rating" json:"rating"`
	Image  struct {
		Src string `yaml:"src" json:"src"`
		Alt string `yaml:"alt" json:"alt"`
	}
}

func getRecommandations(c *gin.Context) {

	buf, err := ioutil.ReadFile("config/data/" + c.Param("alias") + ".yaml")
	if err != nil {
		log.Fatal(err)
		// c.Status(http.StatusNotFound)
		c.IndentedJSON(http.StatusOK, nil)
		return
	}

	var data Data
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		log.Fatal(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.IndentedJSON(http.StatusOK, data)
}

func main() {

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://crossyshop-widget.herokuapp.com/widget/", "https://cdn.shopify.com"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"X-API-KEY"},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	router.GET("/recommandations/:alias", validateAPIKey(), getRecommandations)

	router.Run()
}
