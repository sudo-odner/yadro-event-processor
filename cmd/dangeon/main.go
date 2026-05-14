package main

import (
	"fmt"

	"github.com/sudo-odner/yadro-event-processor/internal/config"
)

func main() {
	// Load config
	cfg := config.MustLoad("./config/config.json")
	fmt.Println(cfg)
}
