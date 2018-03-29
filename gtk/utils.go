package main

import (
	"github.com/gotk3/gotk3/gtk"
	"log"
)

func GetListStore(b *gtk.Builder, id string) (listStore *gtk.ListStore) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("List store error: %s", e)
		return nil
	}

	listStore, _ = obj.(*gtk.ListStore)
	return
}

func GetButton(b *gtk.Builder, id string) (btn *gtk.Button) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Button error: %s", e)
		return nil
	}

	btn, _ = obj.(*gtk.Button)
	return
}

func GetTreeView(b *gtk.Builder, id string) (treeView *gtk.TreeView) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Tree view error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.TreeView)
	return
}

func GetLabel(b *gtk.Builder, id string) (treeView *gtk.Label) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Label error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Label)
	return
}

func GetScrolledWindow(b *gtk.Builder, id string) (treeView *gtk.ScrolledWindow) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Scrolled window error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.ScrolledWindow)
	return
}

func GetSpinner(b *gtk.Builder, id string) (treeView *gtk.Spinner) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Spinner error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Spinner)
	return
}

func GetEntry(b *gtk.Builder, id string) (treeView *gtk.Entry) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("Entry error: %s", e)
		return nil
	}

	treeView, _ = obj.(*gtk.Entry)
	return
}

func GetComboBox(b *gtk.Builder, id string) (combobox *gtk.ComboBox) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("ComboBox error: %s", e)
		return nil
	}

	combobox, _ = obj.(*gtk.ComboBox)
	return
}

func GetCheckButton(b *gtk.Builder, id string) (el *gtk.CheckButton) {
	obj, e := b.GetObject(id)
	if e != nil {
		log.Printf("CheckButton error: %s", e)
		return nil
	}

	el, _ = obj.(*gtk.CheckButton)
	return
}
