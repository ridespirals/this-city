package config

import "testing"

func TestUIScaleAffectsSizes(t *testing.T) {
	ui := Default().UI
	if ui.SizeLabel() != 16 {
		t.Fatalf("label=%v want 16", ui.SizeLabel())
	}
	if ui.SizeTitle() != 28 {
		t.Fatalf("title=%v want 28", ui.SizeTitle())
	}
	ui.Scale = 2
	if ui.SizeLabel() != 32 {
		t.Fatalf("scaled label=%v want 32", ui.SizeLabel())
	}
	if ui.SizeCaption() != 28 { // 16 * 2 * 0.875
		t.Fatalf("scaled caption=%v want 28", ui.SizeCaption())
	}
}

func TestUIScaleFallback(t *testing.T) {
	ui := Default().UI
	ui.Scale = 0
	if ui.SizeLabel() != 16 {
		t.Fatalf("zero scale should treat as 1, got %v", ui.SizeLabel())
	}
}
