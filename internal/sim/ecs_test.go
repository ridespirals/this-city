package sim

import "testing"

func TestEntityCreateDestroyReuse(t *testing.T) {
	w := NewWorld()
	a := w.Create()
	b := w.Create()
	if a == b {
		t.Fatal("expected distinct entities")
	}
	if !w.Alive(a) || !w.Alive(b) {
		t.Fatal("new entities should be alive")
	}
	if w.Len() != 2 {
		t.Fatalf("Len = %d, want 2", w.Len())
	}

	w.Transforms.Set(a, Transform2D{X: 1, Y: 2})
	if !w.Destroy(a) {
		t.Fatal("Destroy should succeed")
	}
	if w.Alive(a) {
		t.Fatal("destroyed entity should be dead")
	}
	if _, ok := w.Transforms.Get(a); ok {
		t.Fatal("components should be cleared on destroy")
	}
	if w.Destroy(a) {
		t.Fatal("second Destroy should fail")
	}

	c := w.Create()
	if c.Index != a.Index {
		t.Fatalf("expected slot reuse, got index %d want %d", c.Index, a.Index)
	}
	if c.Generation == a.Generation {
		t.Fatal("reused entity must bump generation")
	}
	if w.Alive(a) {
		t.Fatal("stale ID must not be alive after reuse")
	}
	if !w.Alive(c) || !w.Alive(b) {
		t.Fatal("live entities missing")
	}
}

func TestComponentStoreForEach(t *testing.T) {
	w := NewWorld()
	e1 := w.Create()
	e2 := w.Create()
	w.Transforms.Set(e1, Transform2D{X: 3})
	w.Transforms.Set(e2, Transform2D{X: 4})
	sum := float32(0)
	n := 0
	w.Transforms.ForEach(func(_ Entity, v Transform2D) {
		sum += v.X
		n++
	})
	if n != 2 || sum != 7 {
		t.Fatalf("ForEach n=%d sum=%v", n, sum)
	}
}
