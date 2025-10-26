package forms

import (
	"gosync/utils"
	"time"

	"github.com/rivo/tview"
)

func ShowView(sftp *utils.SFTP) {
	app := tview.NewApplication()

	container := tview.NewGrid().SetRows(2, 0, 1).SetColumns(0, 0).SetBorders(true)

	// Add headerStatus
	headerStatus := tview.NewTextView().SetText("gosync status: disconnected\nskrazzo@skrazzo.xyz:/home/skrazzo/gosync")
	container.AddItem(headerStatus,
		0,     // row
		0,     // column
		1,     // rowSpan
		2,     // colSpan
		0,     // minGridHeight
		0,     // minGridWidth
		false, // focus
	)

	// Title render function
	titlePrimitive := func(title string) tview.Primitive {
		return tview.NewTextView().SetText(title).SetTextAlign(tview.AlignCenter)
	}

	uploadQueue := tview.NewList()
	deleteQueue := tview.NewList()

	render := func() {
		// Render queue
		uploadQueue.Clear()
		for _, fileName := range sftp.Queue.Uploads {
			uploadQueue.AddItem(fileName, "", 0, nil)
		}

		// Render delete confirmation queue
		deleteQueue.Clear()
		for _, fileName := range sftp.Queue.Deletes {
			deleteQueue.AddItem(fileName, "", 0, nil)
		}

	}

	// Render queue container
	uploadContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titlePrimitive("Uploads"), 1, 0, false).
		AddItem(uploadQueue, 0, 1, false)

	container.AddItem(uploadContainer, 1, 0, 1, 1, 0, 0, false)

	// Render delete queue
	deleteContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titlePrimitive("Deletes"), 1, 0, false).
		AddItem(deleteQueue, 0, 1, false)

	container.AddItem(deleteContainer, 1, 1, 1, 1, 0, 0, false)

	// Footer
	footer := tview.NewBox().
		SetBorder(true)

	container.AddItem(footer, 2, 0, 1, 2, 0, 0, false)

	// Initial render of dynamic variables
	render()

	// Run ticker for dynamic variable update
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			// Queue for redraw, will redraw after func()
			app.QueueUpdateDraw(func() {
				render()
			})
		}
	}()

	if err := app.SetRoot(container, true).Run(); err != nil {
		panic(err)
	}
}
