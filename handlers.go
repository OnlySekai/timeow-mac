package main

import (
	"time"

	"github.com/gen2brain/dlgs"
	"github.com/getlantern/systray"
	"github.com/sqweek/dialog"
)

func (a *app) handleIdleItemSelected(mIdleTimes []*systray.MenuItem, index int) {
	prevIndex := getMinutesSliceIndexFromDuration(a.maxAllowedIdleTime, idleTimeOptionsInSettings[:])
	if prevIndex >= 0 && prevIndex < len(mIdleTimes) {
		mIdleTimes[prevIndex].Uncheck()
	}
	mIdleTimes[index].Check()

	a.setMaxAllowedIdleTime(int(idleTimeOptionsInSettings[index]))
}

func (a *app) handleKeepTimeLogsForOptionSelected(mKeepTimeLogsForOptions []*systray.MenuItem, index int) {
	prevIndex := getMinutesSliceIndexFromDuration(a.keepTimeLogsFor, keepTimeLogsForOptionsInSettings[:])
	if prevIndex >= 0 && prevIndex < len(mKeepTimeLogsForOptions) {
		mKeepTimeLogsForOptions[prevIndex].Uncheck()
	}
	mKeepTimeLogsForOptions[index].Check()

	a.setKeepTimeLogsFor(int(keepTimeLogsForOptionsInSettings[index]))

	a.checkAndCleanExpiredTimeLogs()
}

func (a *app) handleOpenAtLoginClicked(item *systray.MenuItem) {
	if a.startup.RunningAtStartup() {
		a.startup.RemoveStartupItem()
		item.Uncheck()
	} else {
		a.startup.AddStartupItem()
		item.Check()
	}
}

func (a *app) handleQuitClicked() {
	a.addActivePeriodEntry(a.lastIdleTime, time.Now())

	systray.Quit()
}

func (a *app) handleChangeWebHookAddActivePeriodClicked() {
	url, _, _ := dlgs.Entry("Enter Webhook URL", "Enter the URL of the webhook", a.webhookAddActivePeriod)
	a.setWebhookAddActivePeriod(url)
}

func (a *app) handleForceSyncClicked() {
	syncError, _ := a.savePeriodsToStorage(activePeriodsKey, a.activePeriods)
	if syncError != nil {
		dialog.Message(syncError.Error()).Title("Error").Error()
		return
	}
}

func (a *app) handleAboutClicked() {
	openURL(aboutURL)
}
