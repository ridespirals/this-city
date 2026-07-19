package game

import "github.com/ridespirals/this-city/internal/sim"

// TickEvents ages timed events and despawns them when lifetime elapses.
func TickEvents(w *sim.World, dt float32) {
	if w == nil {
		return
	}
	var doomed []sim.Entity
	w.Events.ForEach(func(e sim.Entity, ev sim.EventSource) {
		if !ev.Active || ev.Lifetime <= 0 {
			return
		}
		ev.Age += dt
		if ev.Age >= ev.Lifetime {
			doomed = append(doomed, e)
			return
		}
		w.Events.Set(e, ev)
	})
	for _, e := range doomed {
		w.Destroy(e)
	}
}
