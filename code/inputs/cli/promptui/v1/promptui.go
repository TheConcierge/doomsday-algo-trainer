package v1

import (
	"context"
	"fmt"
	"github.com/TheConcierge/doomsday-algo-trainer/core"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"os"
	"text/tabwriter"
	"time"
)

const (
	DateFormat = "Jan 2, 2006"
)

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

type promptUiImpl struct {
	dc core.DoomsdayCore
}

type PromptUi interface {
	Start()
}

func NewPromptUi(dc core.DoomsdayCore) PromptUi {
	return &promptUiImpl{dc: dc}
}

func (pui *promptUiImpl) Start() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, core.UserSessionCtxKey, uuid.New().String())

	startDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2099, 12, 12, 23, 59, 59, 0, time.UTC)

	for true {
		fmt.Println("=====================================================================================")
		date, err := pui.dc.GetRandomDate(ctx, startDate, endDate)
		if err != nil {
			fmt.Printf("could not generate date %s", err.Error())
			break
		}

		dateString := date.Date.Format(DateFormat)
		fmt.Println(dateString)

		prompt := promptui.Select{
			Label: "Select Day",
			Items: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
				"Saturday", "Quit"},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("prompt failed %v\n", err)
			continue
		}

		// TODO: don't hardcode this when you refactor to an input package
		if result == "Quit" {
			break
		}

		guess, err := parseWeekday(result)
		if err != nil {
			fmt.Printf("failed to parse result: %v", err)
			continue
		}

		correct, actualDay, err := pui.dc.GuessDay(ctx, date.GuessID, guess)
		if err != nil {
			fmt.Printf("failed to submit guess: %s", err.Error())
			continue
		}

		if correct {
			fmt.Printf("Correct! %s IS a %s\n", dateString, actualDay.String())
			continue
		}

		fmt.Printf("Wrong! You guessed %s but \"%s\" is a %s\n", guess.String(), dateString, actualDay.String())
	}

	// We are finished guessing, print out report
	report, err := pui.dc.GetReport(ctx)
	if err != nil {
		fmt.Println("Could not gather report. Sorry!")
		return
	}

	correctPercent := generateCorrectPercent(report.CorrectGuesses, report.IncorrectGuesses)

	w := tabwriter.NewWriter(os.Stdout, 10, 1, 1, ' ', tabwriter.Debug)
	fmt.Fprintf(w, "Date\tGuess\tActual\t\n")
	for _, guess := range report.Questions {
		// TODO: this is maybe a sign that we just shouldn't include "incomplete" guesses in the report.
		// TODO: Possibly log non-guesses as "SKIPPED"
		if guess.Date != nil && guess.Guess != nil {
			fmt.Fprintf(w, "%s\t%s\t%s\t\n", guess.Date.Format(DateFormat), guess.Guess.String(), guess.Date.Weekday().String())
		}
	}
	w.Flush()

	fmt.Println()
	fmt.Printf("You got %s%% of the questions correct!", correctPercent)
}
