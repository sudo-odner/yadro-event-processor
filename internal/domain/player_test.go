package domain

import (
	"strings"
	"testing"
	"time"
)

func TestPlayer_Report(t *testing.T) {
	tests := []struct {
		name   string
		player Player
		want   string
	}{
		{
			name: "Success report",
			player: Player{
				ID:         1,
				HP:         35,
				Status:     StatusSuccess,
				StartedAt:  14 * time.Hour,
				FinishedAt: 14*time.Hour + 24*time.Minute,
				Floors: []FloorProgress{
					{Number: 1, IsCompleted: true, TotalTimeSpent: 5 * time.Minute},
					{Number: 2, IsBoss: true, IsCompleted: true, TotalTimeSpent: 11 * time.Minute},
				},
			},
			want: "[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35",
		},
		{
			name: "Fail report",
			player: Player{
				ID:         2,
				HP:         0,
				Status:     StatusFail,
				StartedAt:  14*time.Hour + 10*time.Minute,
				FinishedAt: 14*time.Hour + 29*time.Minute,
				Floors: []FloorProgress{
					{Number: 1, IsCompleted: false, TotalTimeSpent: 19 * time.Minute},
					{Number: 2, IsBoss: true, IsCompleted: false},
				},
			},
			want: "[FAIL] 2 [00:19:00, 00:00:00, 00:00:00] HP:0",
		},
		{
			name: "Disqual report",
			player: Player{
				ID:         3,
				HP:         100,
				Status:     StatusDisqual,
				StartedAt:  0,
				FinishedAt: 0,
			},
			want: "[DISQUAL] 3 [00:00:00, 00:00:00, 00:00:00] HP:100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.player.Report()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Report() = %q, want %q", got, tt.want)
			}
		})
	}
}
