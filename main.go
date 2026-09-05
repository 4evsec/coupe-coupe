package main

import (
	"coupecoupe/ui/editor"
	"coupecoupe/ui/theme"
	"flag"
	"fmt"
	"os"

	"gioui.org/app"
)

func OpenImageWindow(filePath string) error {
	ie, err := editor.NewImageEditor(filePath)
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

	theme.SetupTheme()

	for _, filePath := range filePaths {
		OpenImageWindow(filePath)
	}
	app.Main()
}
