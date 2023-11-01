package utils

import (
	"dmxongo/objects"
	"time"
)

func FadeColors(start, end objects.Color, duration time.Duration, steps int) []objects.Color {
	var colors []objects.Color

	// Calculate the step size for each color component
	redStep := float64(end.Red-start.Red) / float64(steps)
	greenStep := float64(end.Green-start.Green) / float64(steps)
	blueStep := float64(end.Blue-start.Blue) / float64(steps)
	whiteStep := float64(end.White-start.White) / float64(steps)

	for i := 0; i <= steps; i++ {
		// Calculate the interpolated color at this step
		interpolatedColor := objects.Color{
			Red:   start.Red + int(redStep*float64(i)),
			Green: start.Green + int(greenStep*float64(i)),
			Blue:  start.Blue + int(blueStep*float64(i)),
			White: start.White + int(whiteStep*float64(i)),
		}

		colors = append(colors, interpolatedColor)

		// Sleep for the duration of each step
		time.Sleep(duration / time.Duration(steps))
	}

	return colors
}

func FadeColorArray(colors []objects.Color, duration time.Duration, steps int) []objects.Color {
	fadedColors := []objects.Color{}

	for i := 0; i < len(colors)-1; i++ {
		startColor := colors[i]
		endColor := colors[i+1]
		fadedColors = append(fadedColors, FadeColors(startColor, endColor, duration, steps)...)
	}

	return fadedColors
}

// Frames system

type Frame struct {
	Colors       []objects.Color
	Duration     int // in framecount
	FadeInSteps  int // in framecount
	FadeOutSteps int // in framecount
}

// from the frames, generate an array of colors. Times are in frames, if fade is for example just one frrame, it goes instantly

func FrameToColors(frames []Frame) []objects.Color {
	var colors []objects.Color

	for _, frame := range frames {
		// Add the colors for the frame
		colors = append(colors, frame.Colors...)

		// Add a fade-in transition
		if frame.FadeInSteps > 0 {
			fadeInColors := FadeColors(objects.Color{}, frame.Colors[0], time.Duration(frame.FadeInSteps)*time.Second, frame.FadeInSteps)
			colors = append(colors, fadeInColors...)
		}

		// Add a fade-out transition
		if frame.FadeOutSteps > 0 {
			fadeOutColors := FadeColors(frame.Colors[len(frame.Colors)-1], objects.Color{}, time.Duration(frame.FadeOutSteps)*time.Second, frame.FadeOutSteps)
			colors = append(colors, fadeOutColors...)
		}

		// Add a duration of black frames
		blackFrame := objects.Color{} // You can change this to represent your black color
		for i := 0; i < frame.Duration; i++ {
			colors = append(colors, blackFrame)
		}
	}

	return colors
}
