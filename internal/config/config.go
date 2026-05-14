package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sudo-odner/yadro-event-processor/internal/domain"
)

type Config struct {
	Floors   int                `json:"Floors"    env-required:"true"` // Number of floors in the dungeon
	Monsters int                `json:"Monsters"  env-required:"true"` // Number of monsters on each floor of the dungeon
	OpenAt   domain.DungeonTime `json:"OpenAt"    env-required:"true"` // Number of floors in the dungeon in HH:MM:SS format
	Duration int                `json:"Duration" env-required:"true"`  // Time until the dungeon closes in hours
	CloseAt  time.Duration
}

func MustLoad(configPath string) *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("FATAL: failed load config: %v", err)
	}

	cfg.CloseAt = cfg.OpenAt.Duration + time.Duration(cfg.Duration)*time.Hour
	return &cfg
}
