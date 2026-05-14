package eventparser

import (
	"testing"
	"time"

	"github.com/sudo-odner/yadro-event-processor/internal/domain"
)

func TestEventParser_ParseLine(t *testing.T) {
	parser := New()

	tests := []struct {
		name    string
		line    string
		want    *domain.Event
		wantErr bool
	}{
		{
			name: "valid registration",
			line: "[14:00:00] 1 1",
			want: &domain.Event{
				Time:     14 * time.Hour,
				PlayerID: 1,
				Type:     domain.EvRegister,
			},
			wantErr: false,
		},
		{
			name: "valid with extra param",
			line: "[14:27:00] 2 11 60",
			want: &domain.Event{
				Time:       14*time.Hour + 27*time.Minute,
				PlayerID:   2,
				Type:       domain.EvDamage,
				ExtraParam: "60",
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			line:    "invalid line",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid time",
			line:    "[14:60:00] 1 1",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "negative ID",
			line:    "[14:00:00] -1 1",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Time != tt.want.Time || got.PlayerID != tt.want.PlayerID || got.Type != tt.want.Type || got.ExtraParam != tt.want.ExtraParam {
					t.Errorf("ParseLine() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEventParser_MidnightTransition(t *testing.T) {
	parser := New()

	// First event at 23:59:59
	e1, _ := parser.ParseLine("[23:59:59] 1 1")
	if e1.Time != 23*time.Hour+59*time.Minute+59*time.Second {
		t.Errorf("Expected 23:59:59, got %v", e1.Time)
	}

	// Second event at 00:00:01 (next day)
	e2, _ := parser.ParseLine("[00:00:01] 1 2")
	expected := 24*time.Hour + 1*time.Second
	if e2.Time != expected {
		t.Errorf("Expected 24:00:01 (next day), got %v", e2.Time)
	}
}
