package core

import (
	"context"
	"time"
)

type SessionManager interface {
	// LogNewDate is expected to be implemented in a way where the "id" can be used to find the guess later. This will
	// be important when recording user guesses.
	LogNewDate(ctx context.Context, time time.Time) (id string, err error)
	GetGuess(ctx context.Context, guessId string) (DateGuess, error)
	UpdateGuess(ctx context.Context, id string, answer DateGuess) error
	GetAllGuesses(ctx context.Context) ([]DateGuess, error)
}
