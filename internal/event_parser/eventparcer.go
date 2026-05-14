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

func (p *EventParser) ParceLine(line string) (*domain.Event, error) {
	var hh, mm, ss, playerID, eventID int

	if _, err := fmt.Sscanf(line, "[%d:%d:%d] %d %d", &hh, &mm, &ss, &playerID, &eventID); err != nil {
		return nil, fmt.Errorf("failed parce event line: %w", err)
	}

	currectTime := time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute + time.Duration(ss)*time.Second
	if currectTime < p.lastTime {
		p.daysOffset += 24 * time.Hour
	}

	parts := strings.Fields(line)
	var extra string
	if len(parts) > 3 {
		extra = strings.Join(parts[3:], " ")
	}

	return &domain.Event{
		Time:       currectTime + p.daysOffset,
		PlayerID:   playerID,
		Type:       domain.EventType(eventID),
		ExtraParam: extra,
	}, nil
}
