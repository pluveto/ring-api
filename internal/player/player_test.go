package player

import (
	"encoding/binary"
	"errors"
	"math"
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

	sampleRate := 44100
	duration := 1    // seconds
	frequency := 440 // Hz, A4 note

	numSamples := sampleRate * duration

	// WAV file header (same as before, but adjusted for data size)
	header := []byte{
		0x52, 0x49, 0x46, 0x46, // ChunkID: "RIFF"
		36, 0, 0, 0, // ChunkSize: Will be updated later
		0x57, 0x41, 0x56, 0x45, // Format: "WAVE"
		0x66, 0x6d, 0x74, 0x20, // Subchunk1ID: "fmt "
		16, 0, 0, 0, // Subchunk1Size: 16 for PCM
		1, 0, // AudioFormat: 1 (PCM)
		1, 0, // NumChannels: 1 (Mono)
		0x44, 0xac, 0, 0, // SampleRate: 44100
		0x88, 0x58, 0x01, 0, // ByteRate: SampleRate * NumChannels * BitsPerSample/8
		2, 0, // BlockAlign: NumChannels * BitsPerSample/8
		16, 0, // BitsPerSample: 16
		0x64, 0x61, 0x74, 0x61, // Subchunk2ID: "data"
		0, 0, 0, 0, // Subchunk2Size: Will be updated later
	}

	// Calculate data size and update header
	dataSize := numSamples * 2                                      // 2 bytes per sample (16-bit)
	binary.LittleEndian.PutUint32(header[4:8], uint32(36+dataSize)) // Update ChunkSize
	binary.LittleEndian.PutUint32(header[40:44], uint32(dataSize))  // Update Subchunk2Size

	if _, err := f.Write(header); err != nil {
		t.Fatal(err)
	}

	// Generate sine wave data
	for i := 0; i < numSamples; i++ {
		// Generate a sine wave
		sample := math.Sin(2 * math.Pi * float64(frequency) * float64(i) / float64(sampleRate))

		// Convert to 16-bit signed integer (range: -32768 to 32767) and reduce volume to 1/4
		intValue := int16(sample * 32767 * 0.25)

		// Write the sample to the file
		err := binary.Write(f, binary.LittleEndian, &intValue)
		if err != nil {
			t.Fatal(err)
		}
	}

	return f.Name()
}
