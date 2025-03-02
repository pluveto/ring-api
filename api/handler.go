package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/pluveto/ring-api/internal/player"
)

type Handler struct {
	player    player.Player
	audioPath string
}

func NewHandler(p player.Player, audioPath string) *Handler {
	return &Handler{
		player:    p,
		audioPath: audioPath,
	}
}

func (h *Handler) HandleRing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.player.Play(h.audioPath); err != nil {
		switch {
		case errors.Is(err, player.ErrPlaybackInProgress):
			http.Error(w, err.Error(), http.StatusTooManyRequests)
		case errors.Is(err, player.ErrUnsupportedOS):
			http.Error(w, err.Error(), http.StatusNotImplemented)
		case errors.Is(err, player.ErrInvalidAudioPath):
			http.Error(w, "Invalid audio configuration", http.StatusInternalServerError)
		default:
			http.Error(w, fmt.Sprintf("Playback failed: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Playback started"))
}
