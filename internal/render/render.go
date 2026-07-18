// Package render adapts raylib for drawing and raw input.
// It must not run agent AI or own path math beyond display.
package render

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	DefaultWidth  = 1280
	DefaultHeight = 720
	DefaultTitle  = "This City"
	MaxDT         = 0.1 // seconds; hitch clamp for the frame clock
)

// Config controls window creation.
type Config struct {
	Width     int32
	Height    int32
	Title     string
	TargetFPS int32
}

// DefaultConfig returns sensible window defaults.
func DefaultConfig() Config {
	return Config{
		Width:     DefaultWidth,
		Height:    DefaultHeight,
		Title:     DefaultTitle,
		TargetFPS: 60,
	}
}

// Window wraps raylib window lifecycle.
type Window struct {
	cfg Config
}

// Open creates the raylib window and sets the target frame rate.
func Open(cfg Config) *Window {
	if cfg.Width <= 0 {
		cfg.Width = DefaultWidth
	}
	if cfg.Height <= 0 {
		cfg.Height = DefaultHeight
	}
	if cfg.Title == "" {
		cfg.Title = DefaultTitle
	}
	if cfg.TargetFPS <= 0 {
		cfg.TargetFPS = 60
	}
	rl.InitWindow(cfg.Width, cfg.Height, cfg.Title)
	rl.SetTargetFPS(cfg.TargetFPS)
	return &Window{cfg: cfg}
}

// Close shuts down the raylib window.
func (w *Window) Close() {
	if w == nil {
		return
	}
	rl.CloseWindow()
}

// ShouldClose reports whether the user requested quit.
func (w *Window) ShouldClose() bool {
	return rl.WindowShouldClose()
}

// FrameDT returns the last frame time in seconds, clamped for sim stability.
func (w *Window) FrameDT() float32 {
	dt := rl.GetFrameTime()
	if dt > MaxDT {
		return MaxDT
	}
	if dt < 0 {
		return 0
	}
	return dt
}

// BeginFrame starts a draw frame with a clear color.
func BeginFrame(clear rl.Color) {
	rl.BeginDrawing()
	rl.ClearBackground(clear)
}

// EndFrame finishes a draw frame.
func EndFrame() {
	rl.EndDrawing()
}

// DrawPlaceholder paints a Phase-2 splash so the window is visibly alive.
func DrawPlaceholder(paused bool, simTime float64) {
	rl.DrawText("This City", 40, 40, 40, rl.RayWhite)
	rl.DrawText("Phase 2 — module + loop", 40, 100, 20, rl.LightGray)
	status := "running"
	if paused {
		status = "paused (Space)"
	}
	rl.DrawText(status, 40, 140, 20, rl.LightGray)
	rl.DrawText(fmt.Sprintf("sim time: %.1fs", simTime), 40, 180, 20, rl.LightGray)
	rl.DrawText("Esc or close window to quit", 40, 220, 18, rl.Gray)
}
