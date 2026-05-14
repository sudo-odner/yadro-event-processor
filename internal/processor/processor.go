package processor

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/sudo-odner/yadro-event-processor/internal/config"
	"github.com/sudo-odner/yadro-event-processor/internal/domain"
)

type Processor struct {
	cfg     *config.Config
	players map[int]*domain.Player
	events  []string // Буфер чтобы потом вывести все ивенты
}

func New(cfg *config.Config) *Processor {
	return &Processor{
		cfg:     cfg,
		players: make(map[int]*domain.Player),
	}
}

func (p *Processor) ProcessEvent(ev *domain.Event) {
	player, ok := p.players[ev.PlayerID]
	if !ok {
		player = &domain.Player{
			ID:     ev.PlayerID,
			HP:     100,
			Status: domain.StatusIncomplete,
		}
		p.players[ev.PlayerID] = player
	}

	if player.Status != "" && player.Status != domain.StatusIncomplete {
		// Игрок зкончил свой путь странствий
		return
	}

	// Если пытаемся зайти в подземелье, когда уже закрыта
	if ev.Time > p.cfg.CloseAt {
		p.finishPlayer(player, ev.Time, domain.StatusFail)
		return
	}

	// Если не зарегистрированный человек пытается войти в подземелье
	if ev.Type != domain.EvRegister && !player.IsRegistered {
		p.disqualify(player, ev.Time)
		return
	}

	switch ev.Type {
	case domain.EvRegister:
		player.IsRegistered = true
		p.logEvent(ev.Time, "Player [%d] registered", player.ID)

	case domain.EvEntered:
		if player.StartedAt != 0 {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		if ev.Time < p.cfg.OpenAt.Duration {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		player.StartedAt = ev.Time
		player.Status = domain.StatusIncomplete
		player.CurrentFloor = 0
		player.Floors = make([]domain.FloorProgress, p.cfg.Floors)
		for i := 0; i < len(player.Floors)-1; i++ {
			player.Floors[i] = domain.FloorProgress{
				Number:   i + 1,
				Monsters: p.cfg.Monsters,
			}
		}
		player.Floors[len(player.Floors)-1] = domain.FloorProgress{
			Number: 1 + len(player.Floors) - 1,
			IsBoss: true,
		}
		player.Floors[0].LastEntryTime = ev.Time
		p.logEvent(ev.Time, "Player [%d] entered the dungeon", player.ID)

	case domain.EvKillMonster:
		if player.StartedAt == 0 || player.CurrentFloor >= p.cfg.Floors || player.Floors[player.CurrentFloor].IsBoss {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		f := &player.Floors[player.CurrentFloor]
		if f.IsCompleted {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		f.MonstersKilled++
		p.logEvent(ev.Time, "Player [%d] killed the monster", player.ID)
		if f.MonstersKilled >= f.Monsters {
			f.IsCompleted = true
			f.TotalTimeSpent += ev.Time - f.LastEntryTime
		}

	case domain.EvNextFloor:
		if player.StartedAt == 0 || player.CurrentFloor >= p.cfg.Floors-1 {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		f := &player.Floors[player.CurrentFloor]
		if !f.IsCompleted {
			f.TotalTimeSpent += ev.Time - f.LastEntryTime
		}
		player.CurrentFloor++
		player.Floors[player.CurrentFloor].LastEntryTime = ev.Time
		p.logEvent(ev.Time, "Player [%d] went to the next floor", player.ID)

	case domain.EvPreviousFloor:
		if player.StartedAt == 0 || player.CurrentFloor <= 0 {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		f := &player.Floors[player.CurrentFloor]
		if !f.IsCompleted {
			f.TotalTimeSpent += ev.Time - f.LastEntryTime
		}
		player.CurrentFloor--
		player.Floors[player.CurrentFloor].LastEntryTime = ev.Time
		p.logEvent(ev.Time, "Player [%d] went to the previous floor", player.ID)

	case domain.EvEnteredTheBossFloor:
		if player.StartedAt == 0 || player.CurrentFloor != p.cfg.Floors-1 {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		p.logEvent(ev.Time, "Player [%d] entered the boss's floor", player.ID)

	case domain.EvKillBoss:
		if player.StartedAt == 0 || player.CurrentFloor != p.cfg.Floors-1 || !player.Floors[player.CurrentFloor].IsBoss {
			p.inpossibleMove(player, ev.Time, ev.Type)
			return
		}
		f := &player.Floors[player.CurrentFloor]
		f.IsCompleted = true
		f.TotalTimeSpent = ev.Time - f.LastEntryTime
		p.logEvent(ev.Time, "Player [%d] killed the boss", player.ID)

	case domain.EvLeftDungeon:
		status := domain.StatusFail
		if p.isDungeonComplete(player) {
			status = domain.StatusSuccess
		}
		p.finishPlayer(player, ev.Time, status)
		p.logEvent(ev.Time, "Player [%d] left the dungeon", player.ID)

	case domain.EvReason:
		if player.StartedAt == 0 {
			p.disqualify(player, ev.Time)
		} else {
			p.finishPlayer(player, ev.Time, domain.StatusFail)
		}

	case domain.EvHealth:
		val, _ := strconv.Atoi(ev.ExtraParam)
		player.HP += val
		if player.HP > 100 {
			player.HP = 100
		}
		p.logEvent(ev.Time, "Player [%d] has restored [%s] of health", player.ID, ev.ExtraParam)

	case domain.EvDamage:
		val, _ := strconv.Atoi(ev.ExtraParam)
		player.HP -= val
		p.logEvent(ev.Time, "Player [%d] recieved [%s] of damage", player.ID, ev.ExtraParam)
		if player.HP <= 0 {
			player.HP = 0
			p.finishPlayer(player, ev.Time, domain.StatusFail)
			p.logEvent(ev.Time, "Player [%d] is dead", player.ID)
		}
	}
}

func (p *Processor) isDungeonComplete(player *domain.Player) bool {
	for _, f := range player.Floors {
		if !f.IsCompleted {
			return false
		}
	}
	return true
}

func (p *Processor) logEvent(t time.Duration, format string, args ...any) {
	hh := t / time.Hour
	mm := (t % time.Hour) / time.Minute
	ss := (t % time.Minute) / time.Second
	timeStr := fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
	p.events = append(p.events, fmt.Sprintf("%s %s", timeStr, fmt.Sprintf(format, args...)))
}

func (p *Processor) inpossibleMove(player *domain.Player, t time.Duration, evType domain.EventType) {
	p.logEvent(t, "Player [%d] makes imposible move [%d]", player.ID, evType)
}

func (p *Processor) disqualify(player *domain.Player, t time.Duration) {
	player.Status = domain.StatusDisqual
	player.FinishedAt = t
	p.logEvent(t, "Player [%d] is disqualified", player.ID)
}

func (p *Processor) finishPlayer(player *domain.Player, t time.Duration, status domain.Status) {
	if player.Status != domain.StatusIncomplete {
		return
	}
	player.Status = status
	player.FinishedAt = t
	if 0 < player.CurrentFloor && player.CurrentFloor < len(player.Floors) {
		f := &player.Floors[player.CurrentFloor]
		if !f.IsCompleted {
			f.TotalTimeSpent += t - f.LastEntryTime
		}
	}
}

func (p *Processor) GetEvents() []string {
	return p.events
}

func (p *Processor) GetReport() string {
	var report string
	report += "Final report:\n"

	ids := make([]int, 0, len(p.players))
	for id := range p.players {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	for _, id := range ids {
		report += p.players[id].Report() + "\n"
	}
	return report
}
