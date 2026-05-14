package domain

import (
	"fmt"
	"time"
)

type Status string

const (
	StatusIncomplete Status = "INCOMPLETE" // Ситуация когда логи закончились, но игрок в подземелье(или на поверхности) и подземелье не закрыто
	StatusSuccess    Status = "SUCCESS"
	StatusFail       Status = "FAIL"
	StatusDisqual    Status = "DISQUAL"
)

type Player struct {
	ID           int
	HP           int
	Status       Status
	IsRegistered bool
	StartedAt    time.Duration // Время входа в подземелье
	FinishedAt   time.Duration // Время завершения подземелья (выход, смерть и т.д.)

	Floors       []FloorProgress
	CurrentFloor int
}

type FloorProgress struct {
	TotalTimeSpent time.Duration // Общее время нахождения на этаже
	LastEntryTime  time.Duration // Время последнего входа на этаж

	Number         int
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
	if p.Status == StatusIncomplete {
		return fmt.Sprintf("[%s] %d HP:%d", p.Status, p.ID, p.HP)
	}
	var timeInDungeon time.Duration
	if p.StartedAt > 0 {
		timeInDungeon = p.FinishedAt - p.StartedAt
	}

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
