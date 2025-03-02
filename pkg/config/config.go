package config

import (
	"os"
	"path/filepath"
	"runtime"
)

type Config struct {
	AudioPath string
	HTTPPort  string
}

func Load() *Config {
	return &Config{
		AudioPath: getEnv("AUDIO_PATH", defaultAudioPath()),
		HTTPPort:  getEnv("HTTP_PORT", "8080"),
	}
}

func defaultAudioPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("C:\\", "audio", "ring.wav")
	}
	return "/var/lib/ring-api/ring.wav"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
