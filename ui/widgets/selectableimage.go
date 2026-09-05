package widgets

import (
	"coupecoupe/gioui/widget"
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"

	"coupecoupe/math"
)

type SelectableImage struct {
	Image widget.Image

	startPosition f32.Point
	endPosition   f32.Point
	dragging      bool
}

const minimalSelectionDistanceScreenPx = 20

func getSelectionRectangle(gtx layout.Context, start, end f32.Point) *image.Rectangle {
	distX, distY := math.Distance(start, end)
	if max(distX, distY) < minimalSelectionDistanceScreenPx {
		return nil
	}
	isHorizontal := distY >= distX

	var Min, Max image.Point
	// Depending on the crop configuration (vertical or horizontal), get a full
	// width or a full height rectangle.
	if isHorizontal {
		Min = image.Pt(0, int(min(start.Y, end.Y)))
		Max = image.Pt(gtx.Constraints.Max.X, int(max(start.Y, end.Y)))
	} else {
		Min = image.Pt(int(min(start.X, end.X)), 0)
		Max = image.Pt(int(max(start.X, end.X)), gtx.Constraints.Max.Y)
	}
	return &image.Rectangle{Min: Min, Max: Max}
}

// Handles pointer events and updates the selection position and state.
func (s *SelectableImage) handlePointerEvents(gtx layout.Context) {
	for {
		e, ok := gtx.Event(
			pointer.Filter{
				Target: s,
				Kinds:  pointer.Press | pointer.Move | pointer.Drag | pointer.Release,
			},
		)
		if !ok {
			break
		}
		ev := e.(pointer.Event)

		switch ev.Kind {
		case pointer.Press:
			s.dragging = true
			s.startPosition = ev.Position
			s.endPosition = ev.Position
		case pointer.Move, pointer.Drag:
			if s.dragging {
				s.endPosition = ev.Position
			}
		case pointer.Release:
			s.dragging = false
		}

	}
}

// Draws the selection indicator.
func (s *SelectableImage) drawSelectionZone(gtx layout.Context) {
	selection := getSelectionRectangle(gtx, s.startPosition, s.endPosition)
	if selection != nil {
		paint.FillShape(
			gtx.Ops,
			color.NRGBA{R: 120, G: 0, B: 0, A: 100},
			clip.Rect(*selection).Op(),
		)
	}
}

func (s *SelectableImage) Layout(gtx layout.Context) layout.Dimensions {
	dims := s.Image.Layout(gtx)

	defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()

	event.Op(gtx.Ops, s)

	s.handlePointerEvents(gtx)

	if s.dragging {
		s.drawSelectionZone(gtx)
	}

	return dims
}
