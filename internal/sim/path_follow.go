package sim

import "math/rand"

// PathFollower advances an entity along a network edge at constant speed.
// Distance is always parameterized From→To (0 at From, Length at To).
type PathFollower struct {
	Edge     EdgeID
	Distance float32
	Speed    float32
	Forward  bool // true: toward To; false: toward From
}

// AdvancePathFollower updates distance, resolves junctions via PathDecision, writes Transform2D.
func AdvancePathFollower(w *World, e Entity, f PathFollower, dt float32) PathFollower {
	if w == nil || !w.Alive(e) || w.Network == nil {
		return f
	}
	if f.Speed < 0 {
		f.Speed = -f.Speed
	}
	remaining := f.Speed * dt
	const maxHops = 16
	for hop := 0; remaining > 0 && hop < maxHops; hop++ {
		edge, ok := w.Network.GetEdge(f.Edge)
		if !ok || edge.Poly.Length <= 0 {
			break
		}
		length := edge.Poly.Length
		if f.Distance < 0 {
			f.Distance = 0
		}
		if f.Distance > length {
			f.Distance = length
		}

		if f.Forward {
			room := length - f.Distance
			if remaining <= room+1e-6 {
				f.Distance += remaining
				if f.Distance > length {
					f.Distance = length
				}
				remaining = 0
				break
			}
			remaining -= room
			f = crossJunction(w, e, f, edge.To, true)
		} else {
			room := f.Distance
			if remaining <= room+1e-6 {
				f.Distance -= remaining
				if f.Distance < 0 {
					f.Distance = 0
				}
				remaining = 0
				break
			}
			remaining -= room
			f = crossJunction(w, e, f, edge.From, false)
		}
	}

	edge, ok := w.Network.GetEdge(f.Edge)
	if !ok {
		return f
	}
	pos, tangent := edge.Poly.SampleAt(f.Distance)
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

func crossJunction(w *World, e Entity, f PathFollower, node NodeID, viaForward bool) PathFollower {
	dec, hasDec := w.Decisions.Get(e)
	if !hasDec {
		dec = DefaultPathDecision()
	}
	arr := Arrival{Node: node, ViaEdge: f.Edge, Forward: viaForward}
	choice, dec := ChooseNext(w.Network, dec, arr, w.RNG)
	w.Decisions.Set(e, dec)

	f.Edge = choice.Edge
	f.Forward = choice.Forward
	edge, ok := w.Network.GetEdge(f.Edge)
	if !ok {
		return f
	}
	if f.Forward {
		f.Distance = 0
	} else {
		f.Distance = edge.Poly.Length
	}
	return f
}

// TickPathFollowers advances every PathFollower.
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

// PlaceOnEdge snaps a follower onto an edge at distance along From→To.
func PlaceOnEdge(w *World, e Entity, edgeID EdgeID, distance float32, forward bool, speed float32) {
	if w == nil || !w.Alive(e) {
		return
	}
	f := PathFollower{Edge: edgeID, Distance: distance, Speed: speed, Forward: forward}
	if _, ok := w.Decisions.Get(e); !ok {
		w.Decisions.Set(e, DefaultPathDecision())
	}
	w.Followers.Set(e, AdvancePathFollower(w, e, f, 0))
}

// EnsureRNG initializes World.RNG if missing.
func EnsureRNG(w *World, seed int64) {
	if w == nil {
		return
	}
	if w.RNG == nil {
		w.RNG = rand.New(rand.NewSource(seed))
	}
}
