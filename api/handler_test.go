package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pluveto/ring-api/internal/player"
)

type mockPlayer struct {
	playError error
}

func (m *mockPlayer) Play(string) error {
	return m.playError
}

func TestHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		mockError    error
		expectedCode int
	}{
		{"Success", http.MethodGet, nil, http.StatusAccepted},
		{"Method Not Allowed", http.MethodPost, nil, http.StatusMethodNotAllowed},
		{"Playback In Progress", http.MethodGet, player.ErrPlaybackInProgress, http.StatusTooManyRequests},
		{"Unsupported OS", http.MethodGet, player.ErrUnsupportedOS, http.StatusNotImplemented},
		{"Invalid Path", http.MethodGet, player.ErrInvalidAudioPath, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mp := &mockPlayer{playError: tt.mockError}
			handler := NewHandler(mp, "/valid/path.wav")

			req := httptest.NewRequest(tt.method, "/api/ring", nil)
			w := httptest.NewRecorder()

			handler.HandleRing(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}

func TestHandlerErrorMessages(t *testing.T) {
	mp := &mockPlayer{playError: errors.New("custom error")}
	handler := NewHandler(mp, "/valid/path.wav")

	req := httptest.NewRequest(http.MethodGet, "/api/ring", nil)
	w := httptest.NewRecorder()

	handler.HandleRing(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	expectedBody := "Playback failed: custom error\n"
	if body := w.Body.String(); body != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, body)
	}
}
