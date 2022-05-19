package domain

//go:generate stringer -type=EventError
type EventError uint

func (i EventError) Error() string {
	return i.String()
}

const (
	ErrPinTooShort EventError = iota
	ErrNoFace
	ErrTooManyFaces
	ErrFacePhotoNotTaken
	ErrInputDuringPhotoSession
)
