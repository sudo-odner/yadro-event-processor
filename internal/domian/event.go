package domain

import "time"

type EventType int

const (
	EvRegister            EventType = 1
	EvEntered             EventType = 2
	EvKillMonster         EventType = 3
	EvNextFloor           EventType = 4
	EvPreviousFloor       EventType = 5
	EvEnteredTheBossFloor EventType = 6
	EvKillBoss            EventType = 7
	EvLeftDungeon         EventType = 8
	EvReason              EventType = 9
	EvHealth              EventType = 10
	EvDemage              EventType = 11
	EvDisqual             EventType = 31
	EvDead                EventType = 32
	EvImposibleMove       EventType = 33
)

type Event struct {
	Time       time.Duration
	PlayerID   int
	Type       EventType
	ExtraParam string
}
