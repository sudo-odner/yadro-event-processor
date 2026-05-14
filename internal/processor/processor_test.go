package processor

import (
	"strings"
	"testing"
	"time"

	"github.com/sudo-odner/yadro-event-processor/internal/config"
	"github.com/sudo-odner/yadro-event-processor/internal/domain"
)

func TestProcessor_BasicFlow(t *testing.T) {
	cfg := &config.Config{
		Floors:   2,
		Monsters: 1,
		OpenAt:   domain.DungeonTime{Duration: 14 * time.Hour},
		Duration: 2,
	}
	cfg.CloseAt = cfg.OpenAt.Duration + time.Duration(cfg.Duration)*time.Hour
	proc := New(cfg)

	// 1. Register
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour, PlayerID: 1, Type: domain.EvRegister})
	// 2. Enter
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour, PlayerID: 1, Type: domain.EvEntered})
	// 3. Kill monster
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour + 5*time.Minute, PlayerID: 1, Type: domain.EvKillMonster})
	// 4. Next floor
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour + 10*time.Minute, PlayerID: 1, Type: domain.EvNextFloor})
	// 5. Enter boss floor
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour + 10*time.Minute, PlayerID: 1, Type: domain.EvEnteredTheBossFloor})
	// 6. Kill boss
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour + 20*time.Minute, PlayerID: 1, Type: domain.EvKillBoss})
	// 7. Left dungeon
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour + 30*time.Minute, PlayerID: 1, Type: domain.EvLeftDungeon})

	report := proc.GetReport()
	if !strings.Contains(report, "[SUCCESS] 1 [00:30:00, 00:05:00, 00:10:00] HP:100") {
		t.Errorf("Unexpected report: %s", report)
	}
}

func TestProcessor_AutomaticFail(t *testing.T) {
	cfg := &config.Config{
		Floors:   2,
		Monsters: 1,
		OpenAt:   domain.DungeonTime{Duration: 14 * time.Hour},
		Duration: 1,
	}
	cfg.CloseAt = cfg.OpenAt.Duration + time.Duration(cfg.Duration)*time.Hour
	proc := New(cfg)

	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour, PlayerID: 1, Type: domain.EvRegister})
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour, PlayerID: 1, Type: domain.EvEntered})
	
	// Dungeon closes at 15:00:00
	proc.CloseDungeon()

	report := proc.GetReport()
	if !strings.Contains(report, "[FAIL] 1 [01:00:00") {
		t.Errorf("Expected FAIL after 1 hour, got: %s", report)
	}
}

func TestProcessor_DisqualifyUnregistered(t *testing.T) {
	cfg := &config.Config{
		Floors:   2,
		Monsters: 1,
		OpenAt:   domain.DungeonTime{Duration: 14 * time.Hour},
		Duration: 2,
	}
	cfg.CloseAt = cfg.OpenAt.Duration + time.Duration(cfg.Duration)*time.Hour
	proc := New(cfg)

	// Enter without registration
	proc.ProcessEvent(&domain.Event{Time: 14 * time.Hour, PlayerID: 5, Type: domain.EvEntered})

	report := proc.GetReport()
	if !strings.Contains(report, "[DISQUAL] 5") {
		t.Errorf("Expected DISQUAL for unregistered player, got: %s", report)
	}
}
