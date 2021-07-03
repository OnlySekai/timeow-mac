package main

import (
	"fmt"
	"time"

	"github.com/hako/durafmt"
)

type period struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (p *period) duration() time.Duration {
	return calculateDuration(p.Start, p.End)
}

func (p *period) string() string {
	var format string

	if p.Start.Day() == time.Now().Day() {
		format = "15:04"
	} else {
		format = "2 Jan 15:04"
	}
	duration := p.duration()
	limit := 1
	if duration > time.Hour {
		limit = 2
	}

	return fmt.Sprintf(
		"%s - %s (%s)",
		p.Start.Format(format),
		p.End.Format(format),
		durafmt.Parse(duration).LimitFirstN(limit).String(),
	)
}

func (a *app) addBreakEntry(start, end time.Time) {
	entry := period{
		Start: start,
		End:   end,
	}

	// Avoid duplicate breaks.
	length := len(a.breaks)
	if length > 0 {
		lastBreak := a.breaks[length-1]
		if entry.Start.Truncate(time.Minute) == lastBreak.Start.Truncate(time.Minute) {
			if entry.End.After(lastBreak.End) {
				a.breaks[length-1] = entry
			}
		} else {
			a.breaks = append(a.breaks, entry)
		}
	} else {
		a.breaks = append(a.breaks, entry)
	}

	fmt.Println("----------")
	fmt.Printf("New break added: start = %v, end =%v , duration = %v.\n", entry.Start, entry.End, entry.duration())
	fmt.Printf("string: %v\n", entry.string())
	fmt.Println("current time:", time.Now())
	fmt.Println("----------")

	a.updateBreakMenuItems()
	a.saveBreaksToStorage()
}

func (a *app) updateBreakMenuItems() {
	totalNewMenuItems := len(a.breaks) - len(a.breakMenuItems)
	for i := 0; i < totalNewMenuItems; i++ {
		item := a.mBreaks.AddSubMenuItem("", "")
		item.Disable()
		a.breakMenuItems = append(a.breakMenuItems, item)
	}

	length := len(a.breaks)
	for index, entry := range a.breaks {
		a.breakMenuItems[length-index-1].SetTitle(entry.string())
	}
}

func (a *app) readBreaksFromStorage() ([]period, error) {
	var breaks []period
	err := a.defaults.Unmarshal(breaksKey, &breaks)

	return breaks, err
}

func (a *app) saveBreaksToStorage() error {
	return a.defaults.Marshal(breaksKey, a.breaks)
}