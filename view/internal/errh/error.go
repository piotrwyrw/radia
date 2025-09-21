package errh

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/state"
)

func DisplayError(state *state.State, title string, err error) {
	t := title
	if t == "" {
		t = "Error"
	}
	fyne.Do(func() {
		dialog.ShowCustom(title, "Dismiss", widget.NewLabel(err.Error()), state.MainWindow)
	})
}
