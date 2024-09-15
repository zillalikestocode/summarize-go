package main

import (
	"fmt"

	"github.com/zillalikestocode/summarize-api/internal/api"
)

func main() {
	server := api.NewServer("localhost:8080")

	if err := server.Run(); err != nil {
		fmt.Printf("Failed to run server %v", err)
	}
}
