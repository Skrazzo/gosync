package main

import "github.com/rivo/tview"

func ShowView() {
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

	// Render queue
	uploadQueue := tview.NewList().
		AddItem("file1.txt", "", 0, nil).
		AddItem("file2.jpg", "", 0, nil).
		AddItem("file3.pdf", "", 0, nil)

	renameQueue := tview.NewList().
		AddItem("file1.txt", "", 0, nil).
		AddItem("file2.jpg", "", 0, nil).
		AddItem("file3.pdf", "", 0, nil)

	titlePrimitive := func(title string) tview.Primitive {
		return tview.NewTextView().SetText(title).SetTextAlign(tview.AlignCenter)
	}

	queueContainer := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(titlePrimitive("Uploads"), 1, 0, false).
		AddItem(uploadQueue, 0, 1, false).
		AddItem(titlePrimitive("Rename"), 1, 0, false).
		AddItem(renameQueue, 0, 1, false)

	container.AddItem(queueContainer, 1, 0, 1, 1, 0, 0, false)

	// After adding queueContainer
	rightPlaceholder := tview.NewBox().SetBorder(true).SetTitle("Delete confirmation (coming soon)")
	container.AddItem(rightPlaceholder, 1, 1, 1, 1, 0, 0, false)

	// Footer
	footer := tview.NewBox().
		SetBorder(true).
		SetTitle("Footer")

	container.AddItem(footer, 2, 0, 1, 2, 0, 0, false)

	if err := app.SetRoot(container, true).Run(); err != nil {
		panic(err)
	}
}
