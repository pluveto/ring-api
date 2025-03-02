package player

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewPlayer(t *testing.T) {
	p := New()
	if p == nil {
		t.Fatal("New() returned nil")
	}
}

func TestPlaybackFlow(t *testing.T) {
	p := New()
	testFile := createTestAudio(t)
	defer os.Remove(testFile)

	if err := p.Play(testFile); err != nil {
		t.Fatalf("First play failed: %v", err)
	}

	if err := p.Play(testFile); !errors.Is(err, ErrPlaybackInProgress) {
		t.Errorf("Expected ErrPlaybackInProgress, got %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := p.Play(testFile); err != nil {
		t.Errorf("Second play failed: %v", err)
	}
}

func TestConcurrentPlay(t *testing.T) {
	p := New()
	testFile := createTestAudio(t)
	defer os.Remove(testFile)

	var wg sync.WaitGroup
	errors_ch := make(chan error, 10)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p.Play(testFile); err != nil && !errors.Is(err, ErrPlaybackInProgress) {
				errors_ch <- err
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := p.Play(testFile)
		if err != nil && !errors.Is(err, ErrPlaybackInProgress) {
			errors_ch <- err
		}
	}()

	go func() {
		wg.Wait()
		close(errors_ch)
	}()

	for err := range errors_ch {
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	}
}

func TestInvalidAudioPath(t *testing.T) {
	p := New()

	tests := []struct {
		name     string
		path     string
		expected error
	}{
		{"Empty Path", "", ErrInvalidAudioPath},
		{"Relative Path", "../test.wav", ErrInvalidAudioPath},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := p.createCommand(tt.path); !errors.Is(err, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, err)
			}
		})
	}
}

func TestPlatformCommands(t *testing.T) {
	p := New()
	testFile := createTestAudio(t)

	t.Run("Linux Command", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skip("Skipping Linux test")
		}
		cmd, err := p.createCommand(testFile)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Base(cmd.Path) != "aplay" {
			t.Errorf("Expected aplay, got %s", cmd.Path)
		}
	})

	t.Run("Windows Command", func(t *testing.T) {
		if runtime.GOOS != "windows" {
			t.Skip("Skipping Windows test")
		}
		cmd, err := p.createCommand(testFile)
		if err != nil {
			t.Fatal(err)
		}
		if filepath.Base(cmd.Path) != "ffplay.exe" {
			t.Errorf("Expected ffplay.exe, got %s", cmd.Path)
		}
	})
}

func createTestAudio(t *testing.T) string {
	f, err := os.CreateTemp("", "test-*.wav")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	return f.Name()
}
