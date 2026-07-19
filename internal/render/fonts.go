package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/assets/fonts"
)

// UIFonts holds GPU-loaded Space Mono faces for HUD and world labels.
type UIFonts struct {
	Regular rl.Font
	Bold    rl.Font
	ok      bool
}

// Active is the font set used by DrawText helpers. Nil-safe (falls back to default font).
var Active *UIFonts

const defaultSpacing = float32(1)

// LoadFonts uploads Space Mono into GPU memory. Call after InitWindow.
func LoadFonts() *UIFonts {
	// Base size; DrawTextEx scales from this atlas.
	const baseSize int32 = 64
	f := &UIFonts{
		Regular: rl.LoadFontFromMemory(".ttf", fonts.SpaceMonoRegular, baseSize, nil),
		Bold:    rl.LoadFontFromMemory(".ttf", fonts.SpaceMonoBold, baseSize, nil),
		ok:      true,
	}
	rl.SetTextureFilter(f.Regular.Texture, rl.FilterBilinear)
	rl.SetTextureFilter(f.Bold.Texture, rl.FilterBilinear)
	Active = f
	return f
}

// Unload releases GPU font resources.
func (f *UIFonts) Unload() {
	if f == nil || !f.ok {
		return
	}
	rl.UnloadFont(f.Regular)
	rl.UnloadFont(f.Bold)
	f.ok = false
	if Active == f {
		Active = nil
	}
}

func regularFont() rl.Font {
	if Active != nil && Active.ok {
		return Active.Regular
	}
	return rl.GetFontDefault()
}

func boldFont() rl.Font {
	if Active != nil && Active.ok {
		return Active.Bold
	}
	return regularFont()
}

// Text draws a string with Space Mono Regular.
func Text(x, y int32, size float32, text string, color rl.Color) {
	rl.DrawTextEx(regularFont(), text, rl.NewVector2(float32(x), float32(y)), size, defaultSpacing, color)
}

// TextBold draws a string with Space Mono Bold.
func TextBold(x, y int32, size float32, text string, color rl.Color) {
	rl.DrawTextEx(boldFont(), text, rl.NewVector2(float32(x), float32(y)), size, defaultSpacing, color)
}
