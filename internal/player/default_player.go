package player

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
)

var (
	ErrPlaybackInProgress = errors.New("playback in progress")
	ErrUnsupportedOS      = errors.New("unsupported operating system")
	ErrInvalidAudioPath   = errors.New("invalid audio path")
)

type DefaultPlayer struct {
	mu        sync.Mutex
	isPlaying bool
}

func New() *DefaultPlayer {
	return &DefaultPlayer{}
}

func (p *DefaultPlayer) Play(audioPath string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isPlaying {
		return ErrPlaybackInProgress
	}

	if cmd, err := p.createCommand(audioPath); err != nil {
		return err
	} else {
		p.isPlaying = true
		go func() {
			defer p.setPlayStatus(false)
			if err := cmd.Run(); err != nil {
				fmt.Printf("Playback error: %v\n", err)
			}
		}()
	}
	return nil
}

func (p *DefaultPlayer) setPlayStatus(playing bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.isPlaying = playing
}

func (p *DefaultPlayer) createCommand(audioPath string) (*exec.Cmd, error) {
	if audioPath == "" {
		return nil, ErrInvalidAudioPath
	}

	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		return nil, ErrInvalidAudioPath
	}
	switch runtime.GOOS {
	case "linux":
		return exec.Command("aplay", "-q", audioPath), nil
	case "windows":
		return exec.Command("ffplay", "-nodisp", "-autoexit", audioPath), nil
	default:
		return nil, ErrUnsupportedOS
	}
}
