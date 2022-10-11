package core

import (
	"context"
	"time"
)

type RandomDateGenResponse struct {
	Date    time.Time
	GuessID string
}

type DoomsdayCore interface {
	GetRandomDate(ctx context.Context, lowBound time.Time, upBound time.Time) (RandomDateGenResponse, error)
	GuessDay(ctx context.Context, guessId string, guess time.Weekday) (bool, time.Weekday, error)
	GetReport(ctx context.Context) (*Report, error)
}
