package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/sudo-odner/yadro-event-processor/internal/config"
	eventparser "github.com/sudo-odner/yadro-event-processor/internal/event_parser"
	"github.com/sudo-odner/yadro-event-processor/internal/processor"
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
	proc := processor.New(cfg)

	var scanner *bufio.Scanner
	if *eventsPath != "" {
		// если есть читаем файл и выводим ответ
		file, err := os.Open(*eventsPath)
		if err != nil {
			fmt.Printf("ERROR: falied to open file: %v\n", err)
			return
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		// real-time ввод
		scanner = bufio.NewScanner(os.Stdin)
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ev, err := parser.ParceLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing line: %v\n", err)
			continue
		}
		proc.ProcessEvent(ev)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("ERROR: failed during scanning: %v\n", err)
	}

	// Ивенты
	for _, event := range proc.GetEvents() {
		fmt.Println(event)
	}

	// Report
	fmt.Print(proc.GetReport())
}
