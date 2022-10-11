package core

import (
	"context"
	"math/rand"
	"time"
)

const (
	UserSessionCtxKey = "user-session"
)

type doomsdayCoreImpl struct {
	session SessionManager
}

func NewDoomsdayCore(session SessionManager) DoomsdayCore {
	return &doomsdayCoreImpl{
		session: session,
	}
}

func (dc *doomsdayCoreImpl) GetRandomDate(ctx context.Context, lowBound time.Time, upBound time.Time) (RandomDateGenResponse, error) {
	date := newRandomDate(lowBound, upBound)

	id, err := dc.session.LogNewDate(ctx, date)
	if err != nil {
		// TODO: port over errors package from calico-cut-pants. This will need to be cleaned up.
		return RandomDateGenResponse{}, err
	}

	return RandomDateGenResponse{
		date,
		id,
	}, nil
}

func (dc *doomsdayCoreImpl) GuessDay(ctx context.Context, guessId string, guess time.Weekday) (bool, time.Weekday, error) {
	answer, err := dc.session.GetGuess(ctx, guessId)
	if err != nil {
		return false, -1, err
	}

	answer.Guess = &guess

	err = dc.session.UpdateGuess(ctx, guessId, answer)
	if err != nil {
		return false, -1, err
	}

	return guess == answer.Date.Weekday(), answer.Date.Weekday(), nil
}

func (dc *doomsdayCoreImpl) GetReport(ctx context.Context) (*Report, error) {
	guesses, err := dc.session.GetAllGuesses(ctx)
	if err != nil {
		return &Report{}, err
	}

	correctGuesses := 0
	incorrectGuesses := 0
	for _, guess := range guesses {
		// Ignoring non-answered questions for now.
		// TODO: is this what we want?
		if guess.Guess == nil {
			continue
		}

		if guess.Date.Weekday() == *guess.Guess {
			correctGuesses += 1
			continue
		}

		incorrectGuesses += 1
	}

	report := &Report{
		Questions:        guesses,
		CorrectGuesses:   correctGuesses,
		IncorrectGuesses: incorrectGuesses,
	}

	return report, nil
}

func newRandomDate(lowBound time.Time, upBound time.Time) time.Time {
	delta := upBound.Unix() - lowBound.Unix()
	rand.Seed(time.Now().UnixNano())
	sec := rand.Int63n(delta) + lowBound.Unix()
	return time.Unix(sec, 0)
}
