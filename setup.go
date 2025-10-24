package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowSetupForm displays a form to create initial configuration
func ShowSetupForm(app *tview.Application) (*Config, error) {
	config := &Config{}
	completed := false
	var formError error

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" gosync - Initial Setup ").SetTitleAlign(tview.AlignCenter)

	// Add form fields
	form.AddInputField("Local Directory", ".", 50, nil, func(text string) {
		config.LocalDir = text
	})

	form.AddInputField("Remote Host", "", 50, nil, func(text string) {
		config.Host = text
	})

	form.AddInputField("Username", "", 30, nil, func(text string) {
		config.User = text
	})

	form.AddInputField("Remote Directory", "", 50, nil, func(text string) {
		config.RemoteDir = text
	})

	// Auth type dropdown
	authTypes := []string{"password", "key"}
	form.AddDropDown("Auth Type", authTypes, 0, func(option string, optionIndex int) {
		config.AuthType = option
	})
	config.AuthType = "password" // Default

	form.AddPasswordField("Password", "", 30, '*', func(text string) {
		config.Password = text
	})

	form.AddInputField("Private Key Path (if using key auth)", "", 50, nil, func(text string) {
		config.PrivateKeyPath = text
	})

	// Status text
	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusText.SetText("[yellow]Fill in the form and press Save to continue[white]")

	// Buttons
	form.AddButton("Save", func() {
		// Clear password or key path based on auth type
		if config.AuthType == "password" {
			config.PrivateKeyPath = ""
		} else {
			config.Password = ""
		}

		// Validate config
		if err := config.Validate(); err != nil {
			statusText.SetText("[red]Error: " + err.Error() + "[white]")
			return
		}

		// Save config
		if err := SaveConfig(config); err != nil {
			formError = err
			app.Stop()
			return
		}

		completed = true
		app.Stop()
	})

	form.AddButton("Quit", func() {
		app.Stop()
	})

	// Create layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(form, 0, 1, true).
		AddItem(statusText, 3, 0, false)

	// Set up key bindings
	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.Stop()
			return nil
		}
		return event
	})

	// Run the form
	if err := app.SetRoot(flex, true).SetFocus(form).Run(); err != nil {
		return nil, err
	}

	if formError != nil {
		return nil, formError
	}

	if !completed {
		return nil, nil // User quit without saving
	}

	return config, nil
}
