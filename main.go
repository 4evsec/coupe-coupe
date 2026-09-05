package main

import (
	"flag"
	"fmt"
	"os"

	"gioui.org/app"

	"coupecoupe/ui"
)

func OpenImageWindow(filePath string) error {
	ie, err := ui.NewImageEditor(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		return err
	}
	fmt.Fprintf(os.Stdout, "Opening file: %s", ie.OriginalFile)
	ie.CreateWindow()
	return nil
}

func main() {
	flag.Parse()
	filePaths := flag.Args()

	ui.SetupTheme()

	for _, filePath := range filePaths {
		OpenImageWindow(filePath)
	}
	app.Main()
}
