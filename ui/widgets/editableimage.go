package widgets

import (
	"coupecoupe/cutout"
	"coupecoupe/helpers/math"
	"coupecoupe/pkg/gioui/widget"
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

	patchedwidget "coupecoupe/pkg/gioui/widget"
)

const minimalSelectionDistanceScreenPx = 10

type EditableImage struct {
	Image       image.Image
	imageWidget *widget.Image

	selectionStart f32.Point
	selectionEnd   f32.Point
	dragging       bool
	updated        bool
}

func NewEditableImage(inputImage image.Image) *EditableImage {
	e := EditableImage{Image: inputImage}
	e.SetupWidget()
	return &e
}

func isSelectionHorizontal(start, end f32.Point) (bool, error) {
	distX, distY := math.Distance(start, end)
	if max(distX, distY) < minimalSelectionDistanceScreenPx {
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
	s.updated = true
	return nil
}

func (s *EditableImage) SetupWidget() {
	s.imageWidget = &patchedwidget.Image{
		Fit:      patchedwidget.Contain,
		Position: layout.Center,
		Src:      paint.NewImageOp(s.Image),
	}
}

// Handles "crop" pointer drag/mousedown gestures.
func (s *EditableImage) handlePointerEvents(ev pointer.Event) {
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
			err := s.crop()
			if err != nil {
			}
		}
	}
}

// Handles key presses.
func (s *EditableImage) handleKeyEvents(ev key.Event) {
	switch ev.Name {
	case key.NameEscape:
		s.dragging = false
	}
}

// Root event handling logic.
func (s *EditableImage) handleEvents(gtx layout.Context) {
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
			s.handlePointerEvents(ev)
		case key.Event:
			s.handleKeyEvents(ev)
		}
	}
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
	paint.FillShape(gtx.Ops, color.NRGBA{R: 120, G: 0, B: 0, A: 100}, clip.Rect(selection).Op())
	return nil
}

func (s *EditableImage) Layout(gtx layout.Context) layout.Dimensions {
	dims := s.imageWidget.Layout(gtx)
	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	event.Op(gtx.Ops, s)
	s.handleEvents(gtx)

	if s.updated {
		s.updated = false
		s.imageWidget.Src = paint.NewImageOp(s.Image)
		dims = s.imageWidget.Layout(gtx)
	}

	if s.dragging {
		s.drawSelectionZone(gtx)
	}
	return dims
}
