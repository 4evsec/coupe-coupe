package cutout

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

func imageSub(img image.Image, rect image.Rectangle) image.Image {
	if sub, ok := img.(interface {
		SubImage(image.Rectangle) image.Image
	}); ok {
		return sub.SubImage(rect)
	}
	return nil
}

func GetSelectionRectangle(
	start, end image.Point,
	imgMaxBound image.Point,
	isHorizontal bool,
) image.Rectangle {
	var Min, Max image.Point
	if isHorizontal {
		Min = image.Pt(0, min(start.Y, end.Y))
		Max = image.Pt(imgMaxBound.X, max(start.Y, end.Y))
	} else {
		Min = image.Pt(min(start.X, end.X), 0)
		Max = image.Pt(max(start.X, end.X), imgMaxBound.Y)
	}
	return image.Rectangle{Min: Min, Max: Max}.Canon()
}

func Cutout(
	start, end image.Point,
	inputImage image.Image,
	isHorizontal bool,
) (image.Image, error) {
	selectionRectagle := GetSelectionRectangle(start, end, inputImage.Bounds().Max, isHorizontal)
	bandWidth := 50

	// The following assumptions are stated on a vertical axis:
	cropMax := image.Point{X: selectionRectagle.Bounds().Min.X, Y: selectionRectagle.Bounds().Max.Y}
	cropMin := image.Point{X: selectionRectagle.Bounds().Max.X, Y: selectionRectagle.Bounds().Min.Y}

	if isHorizontal {
		cropMax, cropMin = cropMin, cropMax
	}

	crop1 := imageSub(inputImage, image.Rectangle{Min: inputImage.Bounds().Min, Max: cropMax})
	crop2 := imageSub(inputImage, image.Rectangle{Min: cropMin, Max: inputImage.Bounds().Max})

	var canvasRectangle image.Rectangle
	if isHorizontal {
		canvasRectangle = image.Rect(
			0, 0,
			inputImage.Bounds().Dx(),
			(inputImage.Bounds().Dy() + bandWidth - selectionRectagle.Dy()),
		)
	} else {
		canvasRectangle = image.Rect(
			0, 0,
			(inputImage.Bounds().Dx() + bandWidth - selectionRectagle.Dx()),
			inputImage.Bounds().Dy(),
		)
	}

	canvas := image.NewRGBA(canvasRectangle)
	if canvas == nil {
		return nil, fmt.Errorf("An error occured while creating the new image canvas.")
	}
	draw.Draw(
		canvas,
		canvas.Bounds(),
		&image.Uniform{color.RGBA{0, 0, 0, 100}},
		image.Point{},
		draw.Src,
	)

	draw.Draw(canvas, crop1.Bounds(), crop1, image.Point{}, draw.Src)
	draw.Draw(
		canvas,
		image.Rectangle{
			Min: canvas.Bounds().Max.Sub(crop2.Bounds().Size()),
			Max: cropMin.Add(crop2.Bounds().Size()),
		},
		crop2,
		cropMin,
		draw.Src,
	)

	return canvas, nil
}
