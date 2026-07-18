package sim

// AdvancePathFollower updates distance/direction for dt and writes Transform2D.
func AdvancePathFollower(w *World, e Entity, f PathFollower, dt float32) PathFollower {
	if w == nil || !w.Alive(e) {
		return f
	}
	path, ok := w.Paths.Get(f.Path)
	if !ok || path.Poly.Length <= 0 {
		return f
	}
	if f.Speed < 0 {
		f.Speed = -f.Speed
	}
	delta := f.Speed * dt
	if !f.Forward {
		delta = -delta
	}
	f.Distance += delta
	length := path.Poly.Length

	if f.PingPong {
		for f.Distance < 0 || f.Distance > length {
			if f.Distance < 0 {
				f.Distance = -f.Distance
				f.Forward = true
			} else if f.Distance > length {
				f.Distance = 2*length - f.Distance
				f.Forward = false
			}
		}
	} else {
		if f.Distance < 0 {
			f.Distance = 0
			f.Forward = true
		} else if f.Distance > length {
			f.Distance = length
			f.Forward = false
		}
	}

	pos, tangent := path.Poly.SampleAt(f.Distance)
	if !f.Forward {
		tangent = tangent.Scale(-1)
	}
	xf, _ := w.Transforms.Get(e)
	xf.X, xf.Y = pos.X, pos.Y
	xf.Rotation = tangent.Angle()
	if xf.Scale == 0 {
		xf.Scale = 1
	}
	w.Transforms.Set(e, xf)
	return f
}

// TickPathFollowers advances every PathFollower and syncs transforms.
func TickPathFollowers(w *World, dt float32) {
	if w == nil {
		return
	}
	type pair struct {
		e Entity
		f PathFollower
	}
	var batch []pair
	w.Followers.ForEach(func(e Entity, f PathFollower) {
		batch = append(batch, pair{e, f})
	})
	for _, p := range batch {
		w.Followers.Set(p.e, AdvancePathFollower(w, p.e, p.f, dt))
	}
}
