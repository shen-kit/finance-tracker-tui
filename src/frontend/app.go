package frontend

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	app   *tview.Application
	pages *tview.Pages
	flex  *tview.Flex
)

func CreateTUI() {
	setTheme()
	app = tview.NewApplication()
	pages = tview.NewPages()

	catTv := createCategoriesTable()
	cf := createCategoryForm()
	setCategoryTableKeybinds(catTv, cf)

	invTv := createInvestmentsTable()
	inf := createInvestmentForm()
	setInvestmentTableKeybinds(invTv, inf)

	recTv := createRecordsTable()
	rf := createRecordForm()
	setRecordTableKeybinds(recTv, rf)

	createHomepage(recTv, catTv, invTv, rf)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ctrl+D to exit, or ctrl+C/q when on homepage
		if event.Key() == tcell.KeyCtrlD ||
			(flex.GetItemCount() == 1 && (event.Key() == tcell.KeyCtrlC || event.Rune() == 'q')) {
			app.Stop()
			return nil
		} else if event.Key() == tcell.KeyCtrlC { // disable default behaviour (exit app)
			return tcell.NewEventKey(tcell.KeyCtrlC, 0, tcell.ModNone)
		}
		return event
	})

	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}

func setTheme() {
	// frappe -> https://catppuccin.com/palette
	tview.Styles = tview.Theme{
		PrimitiveBackgroundColor:    tcell.NewRGBColor(48, 52, 70),    // Main background color for primitives.
		ContrastBackgroundColor:     tcell.NewRGBColor(129, 200, 190), // Background color for contrasting elements.
		MoreContrastBackgroundColor: tcell.NewRGBColor(239, 159, 118), // Background color for even more contrasting elements.
		BorderColor:                 tcell.NewRGBColor(148, 156, 187), // Box borders.
		TitleColor:                  tcell.NewRGBColor(198, 208, 245), // Box titles.
		PrimaryTextColor:            tcell.NewRGBColor(198, 208, 245), // Primary text.
		SecondaryTextColor:          tcell.NewRGBColor(181, 191, 226), // Secondary text (e.g. labels).
		TertiaryTextColor:           tcell.NewRGBColor(165, 173, 206), // Tertiary text (e.g. subtitles, notes).
		InverseTextColor:            tcell.NewRGBColor(48, 52, 70),    // Text on primary-colored backgrounds.
		ContrastSecondaryTextColor:  tcell.NewRGBColor(65, 69, 89),    // Secondary text on ContrastBackgroundColor-colored backgrounds.
	}
}

func createHomepage(recTv, catTv, invTv *tableView, rf recordForm) {
	flex = tview.NewFlex()

	lv := tview.NewList().
		ShowSecondaryText(false).
		SetSelectedBackgroundColor(tview.Styles.ContrastBackgroundColor).
		AddItem("  Add Record            ", "", 0, func() { showRecordForm(flex, rf, -1, "", "", "", "") }).
		AddItem("  View Month Summary    ", "", 0, nil).
		AddItem("  View Year Summary     ", "", 0, nil).
		AddItem("  Records               ", "", 0, func() { showTable(flex, recTv) }).
		AddItem("  Categories            ", "", 0, func() { showTable(flex, catTv) }).
		AddItem("  Investments           ", "", 0, func() { showTable(flex, invTv) }).
		AddItem("  Quit                  ", "", 0, func() { app.Stop() })

	lv.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'j' {
			return tcell.NewEventKey(tcell.KeyDown, 'j', tcell.ModNone)
		} else if event.Rune() == 'k' {
			return tcell.NewEventKey(tcell.KeyUp, 'k', tcell.ModNone)
		} else if event.Rune() == 'l' {
			return tcell.NewEventKey(tcell.KeyEnter, 'l', tcell.ModNone)
		}
		return event
	})

	lv.SetTitle("Options").SetBorder(true).SetBorderPadding(1, 1, 2, 2)

	flex.AddItem(lv, 30, 0, true)

	pages.AddPage("main", flex, true, true)
}
