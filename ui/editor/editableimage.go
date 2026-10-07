package editor

import (
	"coupecoupe/helpers/math"
	"coupecoupe/pkg/gioui/widget"
	"coupecoupe/services/cutout"
	"coupecoupe/ui/theme"
	"fmt"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	patchedWidget "coupecoupe/pkg/gioui/widget"
)

const SelectionMinimalDistancePx = 10

var SelectionFillColor = color.NRGBA{R: 120, G: 0, B: 0, A: 100}

type EditableImage struct {
	Image       image.Image
	imageWidget *widget.Image

	selectionStart f32.Point
	selectionEnd   f32.Point
	dragging       bool
}

func NewEditableImage(inputImage image.Image) *EditableImage {
	e := EditableImage{Image: inputImage}
	e.SetupWidget()
	return &e
}

func isSelectionHorizontal(start, end f32.Point) (bool, error) {
	distX, distY := math.Distance(start, end)
	if max(distX, distY) < SelectionMinimalDistancePx {
		return true, fmt.Errorf("Selection distance is too short.")
	}
	isHorizontal := distY >= distX
	return isHorizontal, nil
}

func (s *EditableImage) getRealCoordinates(p f32.Point) f32.Point {
	return s.imageWidget.Transform.Invert().Transform(p)
}

func (s *EditableImage) crop() error {
	scaledStart := &s.selectionStart
	scaledEnd := &s.selectionEnd
	isHorizontal, err := isSelectionHorizontal(*scaledStart, *scaledEnd)
	if err != nil {
		return err
	}
	startCoordinate := s.getRealCoordinates(*scaledStart)
	endCoordinate := s.getRealCoordinates(*scaledEnd)

	outputImage, err := cutout.Cutout(
		startCoordinate.Round(),
		endCoordinate.Round(),
		s.Image,
		isHorizontal,
	)
	if err != nil {
		return err
	}
	s.Image = outputImage
	return nil
}

func (s *EditableImage) SetupWidget() {
	s.imageWidget = &patchedWidget.Image{
		Fit:      patchedWidget.Contain,
		Position: layout.Center,
		Src:      paint.NewImageOp(s.Image),
	}
}

// Handles "crop" pointer drag/mousedown gestures.
func (s *EditableImage) handlePointerEvents(ev pointer.Event) bool {
	switch ev.Kind {
	case pointer.Press:
		s.dragging = true
		s.selectionStart = ev.Position
		s.selectionEnd = ev.Position
	case pointer.Move, pointer.Drag:
		if s.dragging {
			s.selectionEnd = ev.Position
		}
	case pointer.Release:
		if s.dragging {
			s.dragging = false
			s.crop()
			return true
		}
	}
	return false
}

// Handles key presses.
func (s *EditableImage) handleKeyEvents(ev key.Event) bool {
	switch ev.Name {
	case key.NameEscape:
		s.dragging = false
	}
	return false
}

// Root event handling logic.
func (s *EditableImage) handleEvents(gtx layout.Context) bool {
	for {
		e, ok := gtx.Event(
			pointer.Filter{
				Target: s,
				Kinds:  pointer.Press | pointer.Move | pointer.Drag | pointer.Release,
			},
			key.Filter{Name: key.NameEscape},
		)
		if !ok {
			break
		}
		switch ev := e.(type) {
		case pointer.Event:
			return s.handlePointerEvents(ev)

		case key.Event:
			return s.handleKeyEvents(ev)
		}
	}
	return false
}

// Draws the selection indicator.
func (s *EditableImage) drawSelectionZone(gtx layout.Context) error {
	isHorizontal, err := isSelectionHorizontal(s.selectionStart, s.selectionEnd)
	if err != nil {
		return err
	}
	selection := cutout.GetSelectionRectangle(
		s.selectionStart.Round(),
		s.selectionEnd.Round(),
		gtx.Constraints.Max,
		isHorizontal,
	)
	paint.FillShape(
		gtx.Ops,
		SelectionFillColor,
		clip.Rect(selection).Op(),
	)
	return nil
}

func (s *EditableImage) Layout(gtx layout.Context) layout.Dimensions {
	var (
		imageStack clip.Stack
		dims       layout.Dimensions
	)

	for {
		paint.FillShape(
			gtx.Ops,
			theme.Theme.Bg,
			clip.Rect(image.Rect(0, 0, dims.Size.X, dims.Size.Y)).Op(),
		)

		s.imageWidget.Src = paint.NewImageOp(s.Image)
		dims = s.imageWidget.Layout(gtx)

		imageStack = clip.Rect{Max: dims.Size}.Push(gtx.Ops)

		event.Op(gtx.Ops, s)
		reupdateRequired := s.handleEvents(gtx)

		if !reupdateRequired {
			break
		}
	}

	if s.dragging {
		s.drawSelectionZone(gtx)
	}

	imageStack.Pop()
	return dims
}
