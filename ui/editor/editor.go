package editor

import (
	"fmt"
	"image"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget"

	chaparwidgets "github.com/chapar-rest/chapar/ui/widgets"

	"coupecoupe/filesystem"
	patchedwidget "coupecoupe/gioui/widget"
	"coupecoupe/ui/theme"
	coupecoupewidgets "coupecoupe/ui/widgets"
)

type ImageEditorConfig struct {
	verticalCrop widget.Bool
}

type ImageEditor struct {
	ImageName    string
	OriginalFile string
	imageWidget  coupecoupewidgets.SelectableImage
	config       ImageEditorConfig
}

func NewImageEditor(filePath string) (*ImageEditor, error) {
	imageFilePath, err := filesystem.ValidateImageFile(filePath)
	if err != nil {
		return nil, err
	}
	imageFD, err := os.Open(imageFilePath)
	if err != nil {
		return nil, err
	}
	defer imageFD.Close()
	image, _, err := image.Decode(imageFD)
	if err != nil {
		return nil, err
	}

	imageWidget := coupecoupewidgets.SelectableImage{Image: patchedwidget.Image{
		Src:      paint.NewImageOp(image),
		Fit:      patchedwidget.Contain,
		Position: layout.Center,
	}}

	return &ImageEditor{
		ImageName:    filesystem.FileName(imageFilePath),
		OriginalFile: imageFilePath,
		imageWidget:  imageWidget,
	}, nil
}

func (ie *ImageEditor) toolbar(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
		Gap:  4,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return chaparwidgets.CheckBox(
				theme.Theme,
				&ie.config.verticalCrop,
				"Vertical crop",
			).Layout(gtx)
		}),
	)
}

func (ie *ImageEditor) drawLayout(gtx layout.Context) {

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(ie.toolbar),
		layout.Flexed(1, ie.imageWidget.Layout),
	)
}

// Creates the editor window.
func (ie *ImageEditor) CreateWindow() {
	go func(ie *ImageEditor) error {
		window := new(app.Window)
		window.Option(
			app.Title(fmt.Sprintf("%s - coupe|coupe", ie.ImageName)),
			app.TopMost(true),
		)

		var ops op.Ops
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				return e.Err

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				paint.Fill(gtx.Ops, theme.Theme.Bg)
				// Setup background color

				ie.drawLayout(gtx)

				e.Frame(gtx.Ops)
			}
		}
	}(ie)
}
