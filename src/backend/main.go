package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mandrigin/gin-spa/spa"
)

var (
	aiSearchEndpoint       = ""
	aiSearchApiKey         = ""
	aiSearchSemanticConfig = ""
	gptFullEndpoint        = ""
	gptApiKey              = ""
	gptModelName           = ""
	client                 = &http.Client{}
)

func init() {
	godotenv.Load()

	aiSearchEndpoint = os.Getenv("AI_SEARCH_ENDPOINT")
	aiSearchApiKey = os.Getenv("AI_SEARCH_API_KEY")
	aiSearchSemanticConfig = os.Getenv("AI_SEARCH_SEMANTIC_CONFIG")
	gptFullEndpoint = os.Getenv("GPT_FULL_ENDPOINT")
	gptApiKey = os.Getenv("GPT_API_KEY")
	gptModelName = os.Getenv("GPT_MODEL_NAME")

	if aiSearchEndpoint == "" || aiSearchApiKey == "" || aiSearchSemanticConfig == "" || gptFullEndpoint == "" || gptApiKey == "" || gptModelName == "" {
		log.Fatal("One or more environment variables are missing: AI_SEARCH_ENDPOINT, AI_SEARCH_API_KEY, AI_SEARCH_SEMANTIC_CONFIG, GPT_FULL_ENDPOINT, GPT_API_KEY, GPT_MODEL_NAME")
	}
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.Use(gin.Logger())

	r.POST("api/ragbot", func(c *gin.Context) {

		// Call AI Search
		var payload RAGRequest
		err := c.ShouldBindJSON(&payload)
		if err != nil {
			// return bad request
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if payload.K == 0 {
			payload.K = 3
		}
		if payload.Relevance == 0 {
			payload.Relevance = 0.75
		}
		if payload.Temperature == 0 {
			payload.Temperature = 0.2
		}

		intent, err := Intent(payload.Input)

		if err != nil {
			// return bad request
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if intent == "Other" {
			messages := append(payload.Messages, Message{Role: "user", Content: payload.Input})
			gptPayload := GPTRequest{Messages: messages, MaxTokens: payload.MaxTokens, Temperature: payload.Temperature}
			gptResponse, err := CallGPT(gptPayload)
			if err != nil {
				// return bad request
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gptResponse)
			return
		}

		// Call AI Search
		aiSearchRequest := AISearchRequest{
			Search:                payload.Input,
			VectorQueries:         []VectorQuery{{Text: payload.Input, Kind: "text", K: payload.K, Fields: "vector"}},
			SemanticConfiguration: aiSearchSemanticConfig,
			Top:                   payload.K,
			QueryType:             "semantic",
			Select:                "chunk_id,chunk,title",
			QueryLanguage:         "en",
		}
		aiSearchResults, err := CallAISearch(aiSearchRequest)
		if err != nil {
			log.Fatal(err)
		}

		sb := strings.Builder{}
		for _, result := range aiSearchResults.Value {
			sb.WriteString(result.Chunk)
		}

		// CALL GPT with the prompt and the concatenated chunks
		input := fmt.Sprintf("%s\nText\"\"\":\n%s\n\"\"\"\nRespond with the provided data only.\n", payload.Input, sb.String())
		messages := append(payload.Messages, Message{Role: "user", Content: input})

		gptPayload := GPTRequest{Messages: messages, MaxTokens: payload.MaxTokens, Temperature: payload.Temperature}
		gptResponse, err := CallGPT(gptPayload)
		if err != nil {
			// return bad request
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gptResponse)
	})

	r.Use(spa.Middleware("/", "./wwwroot"))
	r.Run()
}
