// Package config holds process-wide game settings (window, sim, UI).
package config

// C is the active configuration. Mutate before opening the window, or replace wholesale.
var C = Default()

// Config is the top-level game configuration.
type Config struct {
	Window Window
	Sim    Sim
	UI     UI
}

// Window controls the desktop window.
type Window struct {
	Width     int32
	Height    int32
	Title     string
	TargetFPS int32
}

// Sim controls simulation timing.
type Sim struct {
	// MaxDT clamps a frame hitch so the sim does not jump (seconds).
	MaxDT float32
}

// UI controls interface scale and typography.
type UI struct {
	// Scale multiplies all UI font sizes (and later other UI metrics). 1 = default.
	Scale float32
	Font  Font
}

// Font describes the UI typeface sizing model.
// Pixel size = BasePx * UI.Scale * role multiplier (Title, Body, …).
type Font struct {
	// AtlasSize is the GPU glyph atlas height used at load time.
	AtlasSize int32
	// Spacing is DrawTextEx character spacing in pixels (pre-scale).
	Spacing float32
	// BasePx is the reference size at Scale=1 for a 1.0 relative role (Label).
	BasePx float32

	// Relative role multipliers (unitless).
	Title   float32
	Body    float32
	Label   float32
	Caption float32
}

// Default returns the built-in configuration.
func Default() Config {
	return Config{
		Window: Window{
			Width:     1280,
			Height:    720,
			Title:     "This City",
			TargetFPS: 60,
		},
		Sim: Sim{
			MaxDT: 0.1,
		},
		UI: UI{
			Scale: 1,
			Font: Font{
				AtlasSize: 64,
				Spacing:   1,
				BasePx:    16,
				Title:     1.75,  // 28px at Scale=1
				Body:      1.125, // 18px
				Label:     1.0,   // 16px
				Caption:   0.875, // 14px
			},
		},
	}
}

// px returns BasePx * Scale * rel.
func (u UI) px(rel float32) float32 {
	scale := u.Scale
	if scale <= 0 {
		scale = 1
	}
	base := u.Font.BasePx
	if base <= 0 {
		base = 16
	}
	return base * scale * rel
}

// SizeTitle is the HUD / brand title size in pixels.
func (u UI) SizeTitle() float32 { return u.px(u.Font.Title) }

// SizeBody is the primary UI body text size.
func (u UI) SizeBody() float32 { return u.px(u.Font.Body) }

// SizeLabel is the default control / agent label size.
func (u UI) SizeLabel() float32 { return u.px(u.Font.Label) }

// SizeCaption is the secondary hint / caption size.
func (u UI) SizeCaption() float32 { return u.px(u.Font.Caption) }

// Size returns an arbitrary relative size: BasePx * Scale * rel.
func (u UI) Size(rel float32) float32 { return u.px(rel) }
