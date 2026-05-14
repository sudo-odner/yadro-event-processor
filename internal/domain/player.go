package domain

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusIncomple Status = "INCOMPLETE" // Ситуация когда логи закончились, но иргок в подземелье(или на поаерхности) и подземелье не закрыто
	StatusSuccess  Status = "SUCCESS"
	StatusFail     Status = "FAIL"
	StatusDisqual  Status = "DISQUAL"
)

type Player struct {
	ID         int
	HP         int
	Status     Status
	IsRegisted bool
	StartedAt  time.Duration // Время входа в полземелье
	FinishedAt time.Duration // Время завершения подземельея (выход, смерть и т.д.)

	Floors       []FloorProgress
	CurrentFloor int
}

type FloorProgress struct {
	TotalTimeSpent time.Duration // Общее время находения на этаже
	LastEntryTime  time.Duration // Время последнего входа на этаж

	Numver         int
	Monsters       int
	MonstersKilled int
	IsBoss         bool
	IsCompleted    bool
}

func durationToString(t time.Duration) string {
	hh := t / time.Hour
	mm := (t % time.Hour) / time.Minute
	ss := (t % time.Minute) / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", hh, mm, ss)
}

func (p *Player) Report() string {
	// Случай когда игрок в подземелье(или на поверхности) и подземелье еще не закрыто
	if p.Status == StatusIncomple {
		return fmt.Sprintf("[%s] %d HP:%d", p.Status, p.ID, p.HP)
	}
	timeInDungeon := p.FinishedAt - p.StartedAt

	var sumTimeCompletedFloor time.Duration
	var timeCompletedBoss time.Duration
	var countCompletedFloor int
	for _, floor := range p.Floors {
		if floor.IsCompleted {
			if floor.IsBoss {
				timeCompletedBoss = floor.TotalTimeSpent
			} else {
				sumTimeCompletedFloor += floor.TotalTimeSpent
				countCompletedFloor++
			}
		}
	}

	var avgTimeCompletedFloor time.Duration
	if countCompletedFloor > 0 {
		avgTimeCompletedFloor = sumTimeCompletedFloor / time.Duration(countCompletedFloor)
	}

	return fmt.Sprintf("[%s] %d [%s, %s, %s] HP:%d",
		p.Status,
		p.ID,
		durationToString(timeInDungeon),
		durationToString(avgTimeCompletedFloor),
		durationToString(timeCompletedBoss),
		p.HP,
	)
}
