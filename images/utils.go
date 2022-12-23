package images

// https://stackoverflow.com/questions/746899/how-to-calculate-an-rgb-colour-by-specifying-an-alpha-blending-amount#:~:text=For%20blending%20purposes%20you%20can,transformations%20to%20do%20it%20right.
func blendAlpha(c, a uint32) uint8 {
	alpha := float64(a/257) / 255.0
	color := c / 257
	value := alpha*float64(color) + (1.0-alpha)*255
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
