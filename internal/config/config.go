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

// UI controls interface scale, typography, and chrome layout.
type UI struct {
	// Scale multiplies fonts and layout metrics (toolbar, insets). 1 = default.
	Scale   float32
	Font    Font
	Toolbar Toolbar
}

// Toolbar is the editor tool panel layout at Scale=1 (pixels).
type Toolbar struct {
	X, Y       float32
	BtnW, BtnH float32
	Gap, Pad   float32
	// TextPad is horizontal inset for button labels at Scale=1.
	TextPad float32
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
			Scale: 1.4,
			Font: Font{
				AtlasSize: 64,
				Spacing:   1,
				BasePx:    16,
				Title:     1.75,  // 28px at Scale=1
				Body:      1.125, // 18px
				Label:     1.0,   // 16px
				Caption:   0.875, // 14px
			},
			Toolbar: Toolbar{
				X:       16,
				Y:       16,
				BtnW:    128,
				BtnH:    32,
				Gap:     6,
				Pad:     8,
				TextPad: 8,
			},
		},
	}
}

// Factor returns the effective UI scale (never <= 0).
func (u UI) Factor() float32 {
	if u.Scale <= 0 {
		return 1
	}
	return u.Scale
}

// S scales a layout length by UI.Scale.
func (u UI) S(v float32) float32 { return v * u.Factor() }

// px returns BasePx * Scale * rel.
func (u UI) px(rel float32) float32 {
	base := u.Font.BasePx
	if base <= 0 {
		base = 16
	}
	return base * u.Factor() * rel
}

// ToolbarLayout is toolbar geometry after applying UI.Scale.
type ToolbarLayout struct {
	X, Y       float32
	BtnW, BtnH float32
	Gap, Pad   float32
	TextPad    float32
}

// Layout returns scaled toolbar metrics for hit-testing and drawing.
func (t Toolbar) Layout(scale float32) ToolbarLayout {
	if scale <= 0 {
		scale = 1
	}
	return ToolbarLayout{
		X:       t.X * scale,
		Y:       t.Y * scale,
		BtnW:    t.BtnW * scale,
		BtnH:    t.BtnH * scale,
		Gap:     t.Gap * scale,
		Pad:     t.Pad * scale,
		TextPad: t.TextPad * scale,
	}
}

// ToolbarLayout returns the active scaled toolbar geometry.
func (u UI) ToolbarLayout() ToolbarLayout {
	return u.Toolbar.Layout(u.Factor())
}

// Height is the total panel height for ToolCount tools.
func (l ToolbarLayout) Height(toolCount int) float32 {
	if toolCount <= 0 {
		return l.Pad * 2
	}
	return l.Pad*2 + float32(toolCount)*l.BtnH + float32(toolCount-1)*l.Gap
}

// Width is the total panel width.
func (l ToolbarLayout) Width() float32 {
	return l.Pad*2 + l.BtnW
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
