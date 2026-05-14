package domain

import (
	"fmt"
	"time"
)

type DungeonTime struct {
	time.Duration
}

func (dt *DungeonTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	if len(s) < 2 {
		return fmt.Errorf("invalid time format")
	}
	s = s[1 : len(s)-1]

	var hh, mm, ss int
	if _, err := fmt.Sscanf(s, "%d:%d:%d", &hh, &mm, &ss); err != nil {
		return fmt.Errorf("failed to parse time: %w", err)
	}

	dt.Duration = time.Duration(hh)*time.Hour + time.Duration(mm)*time.Minute + time.Duration(ss)*time.Second
	return nil
}

func (dt *DungeonTime) String() string {
	h := dt.Duration / time.Hour
	m := (dt.Duration % time.Hour) / time.Minute
	s := (dt.Duration % time.Minute) / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
