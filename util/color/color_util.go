package color

import (
	"fmt"
	"github.com/diamondburned/arikawa/v3/discord"
)

func ColorViewerURL(color discord.Color) string {
	r, g, b := color.RGB()
	hex := fmt.Sprintf("%02X%02X%02X", r, g, b)

	return "https://www.color-hex.com/color/" + hex
}

func ColorViewerLink(color discord.Color, text string) string {
	return "[" + text + "](" + ColorViewerURL(color) + ")"
}
