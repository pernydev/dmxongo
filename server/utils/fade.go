package utils

import (
	"dmxongo/objects"
)

func Fade(colors []objects.Color, fadeTimeInFrames int) []objects.Color {
	var result []objects.Color
	for i := 0; i < len(colors)-1; i++ {
		result = append(result, colors[i]) // Add the original color

		startColor := colors[i]
		endColor := colors[i+1]

		for j := 1; j <= fadeTimeInFrames; j++ {
			t := float64(j) / float64(fadeTimeInFrames)
			interpolatedColor := objects.Color{
				Red:   int(float64(startColor.Red)*(1-t) + float64(endColor.Red)*t),
				Green: int(float64(startColor.Green)*(1-t) + float64(endColor.Green)*t),
				Blue:  int(float64(startColor.Blue)*(1-t) + float64(endColor.Blue)*t),
				White: int(float64(startColor.White)*(1-t) + float64(endColor.White)*t),
			}
			result = append(result, interpolatedColor) // Add the interpolated color
		}
	}

	result = append(result, colors[len(colors)-1]) // Add the last color

	return result
}
