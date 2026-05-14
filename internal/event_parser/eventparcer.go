package eventparser

import (
	"fmt"
	"strings"
	"time"

	"github.com/sudo-odner/yadro-event-processor/internal/domain"
)

type EventParser struct {
	lastTime   time.Duration
	daysOffset time.Duration
}

func New() *EventParser {
	return &EventParser{}
}

func (p *EventParser) ParseLine(line string) (*domain.Event, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	var hh, mm, ss, playerID, eventID int
	n, err := fmt.Sscanf(line, "[%d:%d:%d] %d %d", &hh, &mm, &ss, &playerID, &eventID)

	// Валидация
	if err != nil || n < 5 {
		return nil, fmt.Errorf("invalid event format, expected [HH:MM:SS] ID Type: %w", err)
	}

	if hh < 0 || mm < 0 || mm >= 60 || ss < 0 || ss >= 60 {
		return nil, fmt.Errorf("invalid time values: %02d:%02d:%02d", hh, mm, ss)
	}

	// Валидация ID
	if playerID <= 0 {
		return nil, fmt.Errorf("player ID must be positive, got %d", playerID)
	}
	if eventID <= 0 {
		return nil, fmt.Errorf("event ID must be positive, got %d", eventID)
	}

	currentTime := time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute + time.Duration(ss)*time.Second
	if currentTime < p.lastTime {
		p.daysOffset += 24 * time.Hour
	}
	p.lastTime = currentTime

	parts := strings.Fields(line)
	var extra string
	if len(parts) > 3 {
		extra = strings.Join(parts[3:], " ")
	}

	return &domain.Event{
		Time:       currentTime + p.daysOffset,
		PlayerID:   playerID,
		Type:       domain.EventType(eventID),
		ExtraParam: extra,
	}, nil
}
