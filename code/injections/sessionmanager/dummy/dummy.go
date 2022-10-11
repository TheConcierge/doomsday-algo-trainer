package dummy

import (
	"context"
	"fmt"
	"github.com/TheConcierge/doomsday-algo-trainer/core"
	"github.com/google/uuid"
	"time"
)

type dummySessManagerImpl struct {
	storage dummyStorage
}

type dummyStorage struct {
	sessions map[string]Session
}

type Session struct {
	guesses map[string]core.DateGuess
}

func NewDummySessionManager() core.SessionManager {

	return dummySessManagerImpl{
		storage: dummyStorage{
			sessions: map[string]Session{},
		},
	}
}

func (d dummySessManagerImpl) LogNewDate(ctx context.Context, date time.Time) (id string, err error) {
	sessionKey, ok := ctx.Value(core.UserSessionCtxKey).(string)
	if !ok {
		return "", fmt.Errorf("could not get session from context")
	}

	// Create new session if one doesn't already exist
	if _, ok := d.storage.sessions[sessionKey]; !ok {
		d.storage.sessions[sessionKey] = Session{
			guesses: map[string]core.DateGuess{},
		}
	}

	id = uuid.New().String()

	d.storage.sessions[sessionKey].guesses[id] = core.DateGuess{
		Date: &date,
	}

	return id, nil
}

func (d dummySessManagerImpl) GetGuess(ctx context.Context, guessId string) (core.DateGuess, error) {
	sessionKey, ok := ctx.Value(core.UserSessionCtxKey).(string)
	if !ok {
		return core.DateGuess{}, fmt.Errorf("could not get session from context")
	}

	if _, ok := d.storage.sessions[sessionKey]; !ok {
		return core.DateGuess{}, fmt.Errorf("session %s does not exist", sessionKey)
	}

	guess, ok := d.storage.sessions[sessionKey].guesses[guessId]
	if !ok {
		return core.DateGuess{}, fmt.Errorf("guess %s does not exist for session %s", guessId, sessionKey)
	}

	return guess, nil
}

func (d dummySessManagerImpl) UpdateGuess(ctx context.Context, id string, new core.DateGuess) error {
	sessionKey, ok := ctx.Value(core.UserSessionCtxKey).(string)
	if !ok {
		return fmt.Errorf("could not get session from context")
	}

	if _, ok := d.storage.sessions[sessionKey]; !ok {
		return fmt.Errorf("session %s does not exist", sessionKey)
	}

	if _, ok := d.storage.sessions[sessionKey].guesses[id]; !ok {
		return fmt.Errorf("guess %s does not exist for session %s", id, sessionKey)
	}

	d.storage.sessions[sessionKey].guesses[id] = new

	return nil
}

// TODO: it might make more sense to return a map of ID to value. I'm worried about implementation details
//       leaking into core so I'm leaving it for now.
func (d dummySessManagerImpl) GetAllGuesses(ctx context.Context) ([]core.DateGuess, error) {
	sessionKey, ok := ctx.Value(core.UserSessionCtxKey).(string)
	if !ok {
		return []core.DateGuess{}, fmt.Errorf("could not get session from context")
	}

	if _, ok := d.storage.sessions[sessionKey]; !ok {
		return []core.DateGuess{}, fmt.Errorf("session %s does not exist", sessionKey)
	}

	guessMap := d.storage.sessions[sessionKey].guesses

	if guessMap == nil {
		return []core.DateGuess{}, fmt.Errorf("no guesses exist")
	}

	guesses := make([]core.DateGuess, 0, len(guessMap))

	for _, guess := range guessMap {
		guesses = append(guesses, guess)
	}

	return guesses, nil
}
