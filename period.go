package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/getlantern/systray"
	"github.com/hako/durafmt"
)

type period struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	IsSynced bool      `json:"is_synced"`
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

func (a *app) readPeriodsFromStorage(key string) ([]period, error) {
	var periods []period
	err := a.defaults.Unmarshal(key, &periods)

	return periods, err
}

func (a *app) savePeriodsToStorage(key string, periods []period) error {
	// make http request POST to localhost 8000 with body is periods when key is activePeriodsKey
	if key == activePeriodsKey && a.webhookAddActivePeriod != "" {
		a.syncPeriods(periods)
	}
	return a.defaults.Marshal(key, periods)
}

func (a *app) syncPeriods(periods []period) {
	var unsyncedPeriods []period
	for _, p := range periods {
		if !p.IsSynced {
			unsyncedPeriods = append(unsyncedPeriods, p)

		}
	}
	if len(unsyncedPeriods) <= 0 {
		return
	}
	payload, err := json.Marshal(unsyncedPeriods)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	resp, err := http.Post(a.webhookAddActivePeriod, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// set is_synced to true for all periods
	for i := range periods {
		periods[i].IsSynced = true
	}
	defer resp.Body.Close()
}

func updatePeriodMenuItems(
	periods []period,
	periodsMenuItem *systray.MenuItem,
	currentMenuItems []*systray.MenuItem,
) []*systray.MenuItem {
	menuItems := currentMenuItems

	totalNewMenuItems := len(periods) - len(menuItems)
	fmt.Printf("totalNewMenuItems: %v", totalNewMenuItems)

	if totalNewMenuItems < 0 {
		// Hide redundant menu items.
		for i := len(menuItems) - 1; i >= len(menuItems)-(-totalNewMenuItems); i-- {
			menuItems[i].Hide()
		}
	} else {
		// Add missing menu items.
		for i := 0; i < totalNewMenuItems; i++ {
			item := periodsMenuItem.AddSubMenuItem("", "")
			item.Disable()
			item.Show()
			menuItems = append(menuItems, item)
		}
	}

	length := len(periods)
	for index, entry := range periods {
		menuItems[length-index-1].SetTitle(entry.string())
	}

	return menuItems
}
