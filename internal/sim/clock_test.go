package sim

import "testing"

func TestClampDT(t *testing.T) {
	tests := []struct {
		name string
		dt   float32
		max  float32
		want float32
	}{
		{name: "passthrough", dt: 1.0 / 60, max: 0.1, want: 1.0 / 60},
		{name: "clamp high", dt: 1, max: 0.1, want: 0.1},
		{name: "negative", dt: -1, max: 0.1, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ClampDT(tt.dt, tt.max); got != tt.want {
				t.Fatalf("ClampDT(%v, %v) = %v, want %v", tt.dt, tt.max, got, tt.want)
			}
		})
	}
}
