package player

type Player interface {
	Play(audioPath string) error
}
