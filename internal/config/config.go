package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	PgDsn           string
	Port            string
	SongsInfoAPIURL string
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	dns := os.Getenv("PG_DSN")
	if dns == "" {
		log.Fatal("PG_DSN env var not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT env var not set")
	}

	songsInfoAPIURL := os.Getenv("SONGS_INFO_API_URL")
	if songsInfoAPIURL == "" {
		log.Fatal("SONGS_INFO_API_URL env var not set")
	}

	return &Config{
		PgDsn:           dns,
		Port:            port,
		SongsInfoAPIURL: songsInfoAPIURL,
	}
}
