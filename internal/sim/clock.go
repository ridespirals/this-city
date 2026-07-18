package sim

// ClampDT limits a frame delta so a hitch does not blow up the sim step.
func ClampDT(dt, max float32) float32 {
	if dt < 0 {
		return 0
	}
	if dt > max {
		return max
	}
	return dt
}
