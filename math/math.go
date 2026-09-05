package math

import (
	"math"
)

const MinBandWidthPx = 20

func ComputeBandWidth() int {
	return MinBandWidthPx
}

// Takes two positions on a (same) axis, as well as a "minimmal" cutout band
// width, and returns the computed cutout start and end position on this same
// axis.
func ValidateCutoutPositions(pos1, pos2, bandWidthPx int) (int, int) {
	distance := int(math.Abs(float64(pos2 - pos1)))
	if distance < bandWidthPx {
		mean := (pos2 + pos1) / 2
		halfBandWidth := bandWidthPx / 2
		return mean - halfBandWidth, mean + halfBandWidth
	}
	return min(pos1, pos2), max(pos1, pos2)
}
