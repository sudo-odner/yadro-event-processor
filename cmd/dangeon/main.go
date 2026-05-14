package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/sudo-odner/yadro-event-processor/internal/config"
	eventparser "github.com/sudo-odner/yadro-event-processor/internal/event_parser"
)

func main() {
	configPath := flag.String("config", "", "Path configuration file (required)")
	eventsPath := flag.String("events", "", "Path to events log file (optional)")

	flag.Parse()
	if *configPath == "" {
		fmt.Println("ERROR: --config is requered")
		flag.Usage()
		os.Exit(1)
	}

	cfg := config.MustLoad(*configPath)
	parser := eventparser.New()

	if *eventsPath != "" {
		// если есть читаем файл и выводим ответ
		file, err := os.Open(*eventsPath)
		if err != nil {
			fmt.Printf("ERROR: falied to open file: %v\n", err)
			return
		}
		defer file.Close()

		processScanner(bufio.NewScanner(file), cfg, parser)
	} else {
		// real-time ввод
		processScanner(bufio.NewScanner(os.Stdin), cfg, parser)
	}
}

func processScanner(scanner *bufio.Scanner, cfg *config.Config, pr *eventparser.EventParser) {
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		fmt.Println(pr.ParceLine(line))
		// System analiz
		_ = cfg
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR: failed during scanning: %v\n", err)
	}

	// Report
}
