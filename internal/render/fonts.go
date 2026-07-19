package render

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/ridespirals/this-city/assets/fonts"
	"github.com/ridespirals/this-city/internal/config"
)

// UIFonts holds GPU-loaded Space Mono faces for HUD and world labels.
type UIFonts struct {
	Regular rl.Font
	Bold    rl.Font
	ok      bool
}

// Active is the font set used by DrawText helpers. Nil-safe (falls back to default font).
var Active *UIFonts

// LoadFonts uploads Space Mono into GPU memory. Call after InitWindow.
//
// Raylib's default atlas is ASCII-only (32–126). Passing nil codepoints therefore
// drops glyphs like "·" (U+00B7), which then draw as "?".
//
// Note: Space Mono does not include "⌘" (U+2318). Browsers/Google Fonts fake it
// via fallback fonts; raylib will still show "?" for missing glyphs.
func LoadFonts() *UIFonts {
	atlas := config.C.UI.Font.AtlasSize
	if atlas <= 0 {
		atlas = 64
	}
	codepoints := uiCodepoints()
	f := &UIFonts{
		Regular: rl.LoadFontFromMemory(".ttf", fonts.SpaceMonoRegular, atlas, codepoints),
		Bold:    rl.LoadFontFromMemory(".ttf", fonts.SpaceMonoBold, atlas, codepoints),
		ok:      true,
	}
	rl.SetTextureFilter(f.Regular.Texture, rl.FilterBilinear)
	rl.SetTextureFilter(f.Bold.Texture, rl.FilterBilinear)
	Active = f
	return f
}

// uiCodepoints is ASCII + Latin-1 Supplement, plus a few punctuation glyphs
// that Space Mono actually contains (verified via font cmap).
func uiCodepoints() []rune {
	// 32–255 covers basic Latin and Latin-1 (includes · U+00B7).
	cps := make([]rune, 0, 256-32+8)
	for r := rune(32); r <= 255; r++ {
		cps = append(cps, r)
	}
	extras := []rune{
		'—', // em dash
		'–', // en dash
		'…', // ellipsis
		'•', // bullet
	}
	cps = append(cps, extras...)
	return cps
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

func spacing() float32 {
	s := config.C.UI.Font.Spacing
	if s < 0 {
		return 0
	}
	return s * config.C.UI.Scale
}

// Text draws with Space Mono Regular at an explicit pixel size.
func Text(x, y int32, size float32, text string, color rl.Color) {
	rl.DrawTextEx(regularFont(), text, rl.NewVector2(float32(x), float32(y)), size, spacing(), color)
}

// TextBold draws with Space Mono Bold at an explicit pixel size.
func TextBold(x, y int32, size float32, text string, color rl.Color) {
	rl.DrawTextEx(boldFont(), text, rl.NewVector2(float32(x), float32(y)), size, spacing(), color)
}

// TextTitle draws brand / HUD title text.
func TextTitle(x, y int32, text string, color rl.Color) {
	TextBold(x, y, config.C.UI.SizeTitle(), text, color)
}

// TextBody draws primary UI body text.
func TextBody(x, y int32, text string, color rl.Color) {
	Text(x, y, config.C.UI.SizeBody(), text, color)
}

// TextLabel draws control and world labels.
func TextLabel(x, y int32, text string, color rl.Color) {
	Text(x, y, config.C.UI.SizeLabel(), text, color)
}

// TextCaption draws secondary hints and captions.
func TextCaption(x, y int32, text string, color rl.Color) {
	Text(x, y, config.C.UI.SizeCaption(), text, color)
}
