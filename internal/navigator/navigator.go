/*
author: deadc0de6 (https://github.com/deadc0de6)
Copyright (c) 2024, deadc0de6
*/

package navigator

import (
	"gocatcli/internal/log"
	"gocatcli/internal/node"
	"gocatcli/internal/stringer"
	"gocatcli/internal/utils"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// CallbackFunc callback to get list of entries
type CallbackFunc func(string, bool, bool) (bool, []*stringer.Entry)

// Navigator base struct
type Navigator struct {
	app            *tview.Application
	list           *tview.List
	textarea       *tview.TextView
	callBack       CallbackFunc
	path           string
	exitFlag       bool
	goBackFlag     bool
	selectedFlag   bool
	showHiddenFlag bool
	longMode       bool
	helpFlag       bool
	reloadFlag     bool
	hasDotDot      bool
}

var (
	help = `
		j: down
		k: up
		h: go to parent directory
		q: exit
		esc: exit
		L: toggle long mode
		H: toggle hidden files
		enter: open file/directory
		?: show this help.
	`
)

// handle user input
func (a *Navigator) eventHandler(eventKey *tcell.EventKey) *tcell.EventKey {
	if eventKey.Rune() == 'q' {
		// exit
		a.exitFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Key() == tcell.KeyEscape {
		// exit
		a.exitFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Rune() == '?' {
		// help
		a.helpFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Key() == tcell.KeyEnter {
		// open
		a.selectedFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Rune() == 'j' {
		// down
		idx := (a.list.GetCurrentItem() + 1) % a.list.GetItemCount()
		a.list.SetCurrentItem(idx)
		return nil
	} else if eventKey.Rune() == 'k' {
		// up
		idx := a.list.GetCurrentItem() - 1
		a.list.SetCurrentItem(idx)
		return nil
	} else if eventKey.Rune() == 'l' || eventKey.Key() == tcell.KeyRight {
		// open
		a.selectedFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Rune() == 'h' || eventKey.Key() == tcell.KeyLeft {
		// open parent directory
		a.goBackFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Rune() == 'H' {
		// toggle hidden files
		a.showHiddenFlag = !a.showHiddenFlag
		a.reloadFlag = true
		a.app.Stop()
		return nil
	} else if eventKey.Rune() == 'L' {
		a.longMode = !a.longMode
		a.reloadFlag = true
		a.app.Stop()
		return nil
	}
	return eventKey
}

// file list with file infos
// -1 inserts at the end of the list
func (a *Navigator) fillList() []*stringer.Entry {
	// insert entries
	top, entries := a.callBack(a.path, a.showHiddenFlag, a.longMode)
	a.hasDotDot = false
	if !top {
		// insert ".."
		a.hasDotDot = true
		a.list.InsertItem(0, "..", "", 0, nil)
	}
	for _, entry := range entries {
		line := stringer.ColorLineByType(entry.Line, entry.Node, true)
		a.list.InsertItem(-1, line, "", 0, nil)
	}
	return entries
}

// create the view
func (a *Navigator) createList() {
	// create list
	a.list = tview.NewList()
	a.list.ShowSecondaryText(false)
	a.list.SetWrapAround(true)
	a.fillList()

	// current working directory
	a.textarea = tview.NewTextView()
	a.textarea.SetText(a.path)
	a.textarea.SetTextColor(tcell.ColorSlateGray)

	// create layout
	content := tview.NewGrid()
	content.SetRows(1, 0)
	content.SetBorders(false)
	content.AddItem(a.textarea, 0, 0, 1, 1, 0, 0, false)
	content.AddItem(a.list, 1, 0, 1, 1, 0, 0, true)

	// add to app
	a.app.SetRoot(content, true)
	a.app.SetFocus(a.list)
}

// update the list with new content
func (a *Navigator) updateList() []*stringer.Entry {
	a.list.Clear()
	entries := a.fillList()
	a.textarea.SetText(a.path)
	return entries
}

// show modal help
func (a *Navigator) showHelp() {
	old := a.app.GetFocus()
	modal := tview.NewModal()
	modal.SetText(help)
	modal.AddButtons([]string{"Quit"})
	modal.SetDoneFunc(func(_ int, label string) {
		if label == "Quit" {
			a.app.Stop()
		}
	})

	// show modal
	a.app.SetRoot(modal, false)
	a.app.SetFocus(modal)
	err := a.app.Run()
	if err != nil {
		log.Error(err)
	}

	// reset to old primitive
	a.app.SetRoot(old, true)
	a.app.SetFocus(old)
}

// run app
func (a *Navigator) runApp(_ string) {
	var err error

	// user inputs
	a.app.SetInputCapture(a.eventHandler)
	a.createList()

	// navigate
	var entries []*stringer.Entry
	reload := true
	for {
		if reload {
			// reload content
			entries = a.updateList()
			reload = false
			log.Debug("reload")
		}

		// run the app and display list
		err = a.app.Run()
		if err != nil {
			break
		}
		log.Debugf("app exited to recheck...")
		log.Debugf("exitflag: %v", a.exitFlag)
		log.Debugf("helpFlag: %v", a.helpFlag)
		log.Debugf("selectedFlag: %v", a.selectedFlag)
		log.Debugf("goBackFlag: %v", a.goBackFlag)
		log.Debugf("reloadFlag: %v", a.reloadFlag)

		// exit
		if a.exitFlag {
			break
		}

		// show help
		if a.helpFlag {
			a.helpFlag = false
			a.showHelp()
			continue
		}

		// enter directory
		if a.selectedFlag {
			a.selectedFlag = false
			idx := a.list.GetCurrentItem()
			if idx == 0 && a.hasDotDot {
				// ..
				a.goBackFlag = true
			} else {
				// item selected
				if a.hasDotDot {
					idx--
				}

				if idx < 0 || idx >= len(entries) {
					// entry out of range
					break
				}
				entry := entries[idx]
				reload = false
				if entry.Node.GetType() == node.FileTypeArchive {
					a.path = filepath.Join(a.path, entry.Name)
					reload = true
				} else if entry.Node.GetType() == node.FileTypeStorage {
					a.path = filepath.Join(a.path, entry.Name)
					reload = true
				} else if entry.Node.GetType() == node.FileTypeDir {
					a.path = filepath.Join(a.path, entry.Name)
					reload = true
				}
				continue
			}
		}

		// parent directory
		if a.goBackFlag {
			a.goBackFlag = false
			a.goBack()
			reload = true
			continue
		}

		// reload flag for options
		if a.reloadFlag {
			reload = true
			a.reloadFlag = false
		}
	}
}

func (a *Navigator) goBack() {
	if len(a.path) < 1 {
		return
	}

	fields := utils.SplitPath(a.path)
	if len(fields) < 2 {
		a.path = ""
		return
	}
	base := filepath.Dir(a.path)
	a.path = base
}

// Start start the navigator
func (a *Navigator) Start(path string) {
	a.path = path
	a.app = tview.NewApplication()
	a.runApp(path)
}

// NewNavigator creates a new navigator
func NewNavigator(callback CallbackFunc) *Navigator {
	n := Navigator{
		callBack:       callback,
		showHiddenFlag: true,
		longMode:       true,
	}
	return &n
}
