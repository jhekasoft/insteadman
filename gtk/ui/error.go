package ui

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
)

func ShowErrorDlgFatal(txt string) {
	showErrorDlg(txt, true)
}

func ShowErrorDlg(txt string) {
	showErrorDlg(txt, false)
}

func showErrorDlg(txt string, fatal bool) {
	log.Printf("Error: %v", txt)

	dlg, _ := gtk.DialogNew()
	dlg.SetTitle("InsteadMan error")
	dlg.AddButton("Close", gtk.RESPONSE_ACCEPT)
	dlgBox, _ := dlg.GetContentArea()
	dlgBox.SetSpacing(6)

	lbl, _ := gtk.LabelNew(txt)
	lbl.SetMarginStart(6)
	lbl.SetMarginEnd(6)
	lbl.SetLineWrap(true)
	dlgBox.Add(lbl)
	lbl.Show()

	dlg.SetModal(true)
	dlg.SetPosition(gtk.WIN_POS_CENTER)
	dlg.SetResizable(false)
	//dlg.SetTransientFor(window)

	dlg.SetKeepAbove(true)
	dlg.Run()
	dlg.Destroy()
	if fatal {
		os.Exit(1)
	}
}
