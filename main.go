package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type User struct {
	Name 		string `json:"name"`
	Email		string `json:"email"`
	Stack		string `json:"stack"`
}


type CatFactResponse struct {
	Fact		string `json:"fact"`
	Length		int		`json:"length"`
}


type MyAPIResponse struct {
	Status			string		`json:"status"`
	User			User		`json:"user"`
	Timestamp		string		`json:"timestamp"`
	Fact			string		`json:"fact"`
}


var me = User{
	Name:	"Tobi Ade",
	Email: "saintcleverley@gmail.com",
	Stack: "go, ts, postgres",
}

const externalAPIURL = "https://catfact.ninja/fact"


func getRandomCatFact() (*CatFactResponse, error) {
	client := http.Client{
		// Crucial best practice: Always set a timeout for external requests
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(externalAPIURL)
	if err != nil {
		log.Printf("Error fetching fact from external API: %v", err)
		return nil, fmt.Errorf("failed to reach Cat API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("External API returned non-OK: %d, Body: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("cat facts reponded with %d", resp.StatusCode)
	}

	var factResponse CatFactResponse
	if err := json.NewDecoder(resp.Body).Decode(&factResponse); err != nil {
		log.Printf("Error decoding external API JSON: %v", err)
		return nil, fmt.Errorf("failed to decode response from Cat service")
	}

	return &factResponse, nil
}


func getMe(c *gin.Context) {
	factData, err := getRandomCatFact()

	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Could not retrieve me",
			"details": err.Error(),
		})
		return
	}

	utcTimeISO := time.Now().UTC().Format(time.RFC3339)

	
	myResponse := MyAPIResponse{
		Status: "success",
		User:	me, 
		Fact: factData.Fact,
		Timestamp: utcTimeISO,
	}

	c.JSON(http.StatusOK, myResponse)
}

func main() {
	router := gin.Default()
	router.GET("/me", getMe)

	fmt.Println("Server running on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
