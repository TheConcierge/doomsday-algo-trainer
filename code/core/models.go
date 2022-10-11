package core

import "time"

type DateGuess struct {
	Date  *time.Time
	Guess *time.Weekday
}

type Report struct {
	Questions        []DateGuess
	CorrectGuesses   int
	IncorrectGuesses int
}
