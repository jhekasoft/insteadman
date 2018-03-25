package main

import (
	"github.com/gotk3/gotk3/gtk"
	// "log"
	"fmt"
	"log"
	//"os"
)

func main() {
	gtk.Init(nil)

	//_, e := gtk.BuilderNewFromFile("./resources/gtk/about.glade")
	b, e := gtk.BuilderNew()
	if e != nil {
		fmt.Printf("E: %v\n", e)
	}
	e = b.AddFromFile("./resources/gtk/main.glade")
	if e != nil {
		fmt.Printf("E: %v\n", e)
	}

	obj, err := b.GetObject("window_main")
	if err != nil {
		//os.Exit(1)
	}

	window, ok := obj.(*gtk.Window)
	if !ok {
		//os.Exit(1)
	}

	if err != nil {
		log.Fatal(err.Error())
		//os.Exit(1)
	}

	window.SetTitle("InsteadMan 3")
	window.SetDefaultSize(770, 500)
	window.Connect("destroy", destroy)
	window.ShowAll()

	gtk.Main()
}

func destroy() {
	gtk.MainQuit()
}

