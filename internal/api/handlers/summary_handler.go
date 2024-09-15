package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type SummaryHandler struct {
}

type SummaryRequestBody struct {
	URL string `json:"url"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewSummaryHandler() *SummaryHandler {
	return &SummaryHandler{}
}

func (s *SummaryHandler) Summarize(g *gin.Context) {
	var body SummaryRequestBody
	var content strings.Builder
	var wg sync.WaitGroup

	done := make(chan bool)

	if err := g.BindJSON(&body); err != nil {
		fmt.Printf("An error occured while parsing data %v", err)
	}

	c := colly.NewCollector()

	c.OnHTML("body", func(h *colly.HTMLElement) {
		content.WriteString(h.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting the site: ", r.URL)
	})

	wg.Add(1)

	go func() {
		defer wg.Done()
		err := c.Visit(body.URL)
		if err != nil {
			fmt.Println("an error occured")
		}
	}()
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		openAIreq := OpenAIRequest{
			Model: "gpt-3.5-turbo",
			Messages: []Message{
				{
					Content: "You are a helpful assistant that summarizes web content",
					Role:    "system",
				},
				{
					Role:    "user",
					Content: fmt.Sprintf("summarize this website: \n\n %s", content.String()),
				},
			},
		}
		summary, err := openAIreq.Summarize(content.String())
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		g.JSON(http.StatusOK, gin.H{"summary": summary})
	case <-g.Request.Context().Done():
		g.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timeout"})
	}
}

func (openAIreq *OpenAIRequest) Summarize(content string) ([]genai.Part, error) {
	apiKey := "AIzaSyCLZ8Pj7noz1j8efL6CqbALiDtsR1Zd7K0"

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text("summarize this web content: "+content))

	if err != nil {
		return nil, err
	}
	formated := formatResponse(resp)
	return formated, nil

}

func formatResponse(resp *genai.GenerateContentResponse) []genai.Part {
	var arrayofParts []genai.Part
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				arrayofParts = append(arrayofParts, part)
			}
		}
	}
	return arrayofParts
}
