package main

import (
	"context"
	"fmt"
	"github.com/TheConcierge/doomsday-algo-trainer/core"
	"github.com/TheConcierge/doomsday-algo-trainer/injections/sessionmanager/dummy"
	"github.com/google/uuid"
	"github.com/manifoldco/promptui"
	"os"
	"text/tabwriter"
	"time"
)

const (
	DateFormat = "Jan 2, 2006"
)

func main() {
	sess := dummy.NewDummySessionManager()

	c := core.NewDoomsdayCore(sess)

	ctx := context.Background()
	ctx = context.WithValue(ctx, core.UserSessionCtxKey, uuid.New().String())

	startDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2099, 12, 12, 23, 59, 59, 0, time.UTC)

	for true {
		fmt.Println("=====================================================================================")
		date, err := c.GetRandomDate(ctx, startDate, endDate)
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

		correct, actualDay, err := c.GuessDay(ctx, date.GuessID, guess)
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
	report, err := c.GetReport(ctx)
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

func generateCorrectPercent(correct int, incorrect int) string {
	correctF := float64(correct)
	incorrectF := float64(incorrect)

	percent := correctF / (correctF + incorrectF) * 100

	return fmt.Sprintf("%.2f", percent)
}

var daysOfWeek = map[string]time.Weekday{
	"Sunday":    time.Sunday,
	"Monday":    time.Monday,
	"Tuesday":   time.Tuesday,
	"Wednesday": time.Wednesday,
	"Thursday":  time.Thursday,
	"Friday":    time.Friday,
	"Saturday":  time.Saturday,
}

func parseWeekday(v string) (time.Weekday, error) {
	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}
