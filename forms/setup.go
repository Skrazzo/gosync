package forms

import (
	"gosync/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// SetupConfig displays a form to create or edit configuration using a provided Config instance
func SetupConfig(app *tview.Application, cfg *utils.Config) error {
	completed := false
	var formError error

	form := tview.NewForm()
	form.SetBorder(true).SetTitle(" gosync - Config setup ").SetTitleAlign(tview.AlignCenter)

	// Empty ignore pattern
	cfg.Ignore = []string{}

	// Add input form fields
	form.AddInputField("Remote Host", cfg.Host, 50, nil, func(text string) {
		cfg.Host = text
	})

	form.AddInputField("Username", cfg.User, 30, nil, func(text string) {
		cfg.User = text
	})

	form.AddInputField("Remote Directory", cfg.RemoteDir, 50, nil, func(text string) {
		cfg.RemoteDir = text
	})

	// Auth type dropdown
	authTypes := []string{"password", "key"}
	authTypeIndex := 0
	if cfg.AuthType == "key" {
		authTypeIndex = 1
	}
	form.AddDropDown("Auth Type", authTypes, authTypeIndex, func(option string, optionIndex int) {
		cfg.AuthType = option
	})
	if cfg.AuthType == "" {
		cfg.AuthType = "password" // Default
	}

	form.AddPasswordField("Password", cfg.Password, 30, '*', func(text string) {
		cfg.Password = text
	})

	form.AddInputField("Private Key Path (if using key auth)", cfg.PrivateKeyPath, 50, nil, func(text string) {
		cfg.PrivateKeyPath = text
	})

	// Status text
	statusText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusText.SetText("[yellow]Fill in the form and press Save to continue[white]")

	// Buttons
	form.AddButton("Save", func() {
		// Clear password or key path based on auth type
		if cfg.AuthType == "password" {
			cfg.PrivateKeyPath = ""
		} else {
			cfg.Password = ""
		}

		// Validate config
		if err := cfg.Validate(); err != nil {
			statusText.SetText("[red]Error: " + err.Error() + "[white]")
			return
		}

		// Save config
		if err := cfg.Save(); err != nil {
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
		return err
	}

	if formError != nil {
		return formError
	}

	if !completed {
		return nil // User quit without saving
	}

	return nil
}
