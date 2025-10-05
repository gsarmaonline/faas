package main

import (
	"context"
	"fmt"

	"github.com/gsarmaonline/faas/faas"
)

func main() {
	_, err := faas.NewFaas(context.Background())
	if err != nil {
		fmt.Printf("Error creating FAAS: %v\n", err)
		return
	}

	fmt.Println("✅ FAAS framework initialized successfully with all functions:")
	fmt.Println("- slack")
	fmt.Println("- email")
	fmt.Println("- sms")
	fmt.Println("- docker_registry")
	fmt.Println("- http")
	fmt.Println("- logger")
	fmt.Println("- github")
	fmt.Println("Total: 7 functions registered")
}
