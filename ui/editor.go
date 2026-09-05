package ui

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	chaparwidgets "github.com/chapar-rest/chapar/ui/widgets"

	"coupecoupe/filesystem"
)

type Selection struct {
	dragging bool
	start    f32.Point
	end      f32.Point
}

type ImageEditorConfig struct {
	verticalCrop widget.Bool
}

type ImageEditor struct {
	ImageName    string
	OriginalFile string
	image        image.Image
	config       ImageEditorConfig
	selection    Selection
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
	return &ImageEditor{
		ImageName:    filesystem.FileName(imageFilePath),
		OriginalFile: imageFilePath,
		image:        image,
	}, nil
}

func (ie *ImageEditor) toolbar(gtx layout.Context) layout.Dimensions {
	return layout.Flex{
		Axis: layout.Horizontal,
		Gap:  4,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return chaparwidgets.CheckBox(
				Theme,
				&ie.config.verticalCrop,
				"Vertical crop",
			).Layout(gtx)
		}),
	)
}

func (ie *ImageEditor) drawLayout(gtx layout.Context) {
	imageWidget := widget.Image{
		Src:      paint.NewImageOp(ie.image),
		Fit:      widget.Contain,
		Position: layout.Center,
	}

	layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(ie.toolbar),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return imageWidget.Layout(gtx)
		}))
}

// Draws the selection indicator.
// Depending on the crop configuration (vertical or horizontal), draws a full
// width or a full height rectangle.
func (ie *ImageEditor) drawSelectionZone(gtx layout.Context) {
	selection := &ie.selection
	isVertical := ie.config.verticalCrop.Value
	var Min, Max image.Point
	if isVertical {
		Min = image.Pt(int(min(selection.start.X, selection.end.X)), 0)
		Max = image.Pt(int(max(selection.start.X, selection.end.X)), gtx.Constraints.Max.Y)
	} else {
		Min = image.Pt(0, int(min(selection.start.Y, selection.end.Y)))
		Max = image.Pt(gtx.Constraints.Max.X, int(max(selection.start.Y, selection.end.Y)))
	}
	paint.FillShape(
		gtx.Ops,
		color.NRGBA{R: 120, G: 0, B: 0, A: 100},
		clip.Rect{Min: Min, Max: Max}.Op(),
	)
}

// Handles pointer events and updates the selection position and state.
func (ie *ImageEditor) handlePointerEvents(gtx layout.Context) {
	selection := &ie.selection
	for {
		e, ok := gtx.Event(
			pointer.Filter{
				Target: selection,
				Kinds:  pointer.Press | pointer.Move | pointer.Drag | pointer.Release,
			},
		)
		if !ok {
			break
		}
		ev := e.(pointer.Event)

		switch ev.Kind {
		case pointer.Press:
			selection.dragging = true
			selection.start = ev.Position
			selection.end = ev.Position
		case pointer.Move, pointer.Drag:
			if selection.dragging {
				selection.end = ev.Position
			}
		case pointer.Release:
			selection.dragging = false
		}

	}
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

				paint.Fill(gtx.Ops, Theme.Bg)
				// Setup background color

				event.Op(gtx.Ops, &ie.selection)

				ie.handlePointerEvents(gtx)
				ie.drawLayout(gtx)

				if ie.selection.dragging {
					ie.drawSelectionZone(gtx)
				}
				e.Frame(gtx.Ops)
			}
		}
	}(ie)
}
