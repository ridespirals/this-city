package sim

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseSVGPathData converts an SVG path `d` attribute into cubic Bézier curves.
// Supported: M/m L/l H/h V/v C/c S/s Q/q T/t Z/z (absolute and relative).
func ParseSVGPathData(d string) ([]CubicBezier, error) {
	toks, err := tokenizePath(d)
	if err != nil {
		return nil, err
	}
	var (
		out                    []CubicBezier
		i                      int
		cmd                    byte
		x, y                   float32 // current point
		sx, sy                 float32 // subpath start
		cx, cy                 float32 // last control (for S/T)
		haveCtrl               bool
		lastWasCubic, lastWasQ bool
	)
	nextNums := func(n int) ([]float32, error) {
		if i+n > len(toks) {
			return nil, fmt.Errorf("path: need %d numbers for %c", n, cmd)
		}
		vals := make([]float32, n)
		for k := 0; k < n; k++ {
			v, err := strconv.ParseFloat(toks[i], 64)
			if err != nil {
				return nil, fmt.Errorf("path number %q: %w", toks[i], err)
			}
			vals[k] = float32(v)
			i++
		}
		return vals, nil
	}
	isCmd := func(s string) bool {
		if len(s) != 1 {
			return false
		}
		c := s[0]
		return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
	}

	for i < len(toks) {
		if isCmd(toks[i]) {
			cmd = toks[i][0]
			i++
		} else if cmd == 0 {
			return nil, fmt.Errorf("path: number before command")
		}
		// Implicit repetition uses the same command, except M/m → L/l.
		rep := cmd
		switch cmd {
		case 'M', 'm':
			vals, err := nextNums(2)
			if err != nil {
				return nil, err
			}
			if cmd == 'm' {
				x += vals[0]
				y += vals[1]
			} else {
				x, y = vals[0], vals[1]
			}
			sx, sy = x, y
			haveCtrl = false
			lastWasCubic, lastWasQ = false, false
			cmd = 'L'
			if rep == 'm' {
				cmd = 'l'
			}
		case 'Z', 'z':
			if x != sx || y != sy {
				out = append(out, lineCubic(Vec2{X: x, Y: y}, Vec2{X: sx, Y: sy}))
			}
			x, y = sx, sy
			haveCtrl = false
			lastWasCubic, lastWasQ = false, false
			cmd = 0
		case 'L', 'l':
			vals, err := nextNums(2)
			if err != nil {
				return nil, err
			}
			var nx, ny float32
			if cmd == 'l' {
				nx, ny = x+vals[0], y+vals[1]
			} else {
				nx, ny = vals[0], vals[1]
			}
			out = append(out, lineCubic(Vec2{X: x, Y: y}, Vec2{X: nx, Y: ny}))
			x, y = nx, ny
			haveCtrl = false
			lastWasCubic, lastWasQ = false, false
		case 'H', 'h':
			vals, err := nextNums(1)
			if err != nil {
				return nil, err
			}
			nx := vals[0]
			if cmd == 'h' {
				nx = x + vals[0]
			}
			out = append(out, lineCubic(Vec2{X: x, Y: y}, Vec2{X: nx, Y: y}))
			x = nx
			haveCtrl = false
			lastWasCubic, lastWasQ = false, false
		case 'V', 'v':
			vals, err := nextNums(1)
			if err != nil {
				return nil, err
			}
			ny := vals[0]
			if cmd == 'v' {
				ny = y + vals[0]
			}
			out = append(out, lineCubic(Vec2{X: x, Y: y}, Vec2{X: x, Y: ny}))
			y = ny
			haveCtrl = false
			lastWasCubic, lastWasQ = false, false
		case 'C', 'c':
			vals, err := nextNums(6)
			if err != nil {
				return nil, err
			}
			var c0x, c0y, c1x, c1y, nx, ny float32
			if cmd == 'c' {
				c0x, c0y = x+vals[0], y+vals[1]
				c1x, c1y = x+vals[2], y+vals[3]
				nx, ny = x+vals[4], y+vals[5]
			} else {
				c0x, c0y = vals[0], vals[1]
				c1x, c1y = vals[2], vals[3]
				nx, ny = vals[4], vals[5]
			}
			seg := CubicBezier{
				P0: Vec2{X: x, Y: y},
				C0: Vec2{X: c0x, Y: c0y},
				C1: Vec2{X: c1x, Y: c1y},
				P1: Vec2{X: nx, Y: ny},
			}
			if !segNearZero(seg) {
				out = append(out, seg)
			}
			cx, cy = c1x, c1y
			x, y = nx, ny
			haveCtrl = true
			lastWasCubic, lastWasQ = true, false
		case 'S', 's':
			vals, err := nextNums(4)
			if err != nil {
				return nil, err
			}
			var c0x, c0y float32
			if lastWasCubic && haveCtrl {
				c0x, c0y = 2*x-cx, 2*y-cy
			} else {
				c0x, c0y = x, y
			}
			var c1x, c1y, nx, ny float32
			if cmd == 's' {
				c1x, c1y = x+vals[0], y+vals[1]
				nx, ny = x+vals[2], y+vals[3]
			} else {
				c1x, c1y = vals[0], vals[1]
				nx, ny = vals[2], vals[3]
			}
			seg := CubicBezier{
				P0: Vec2{X: x, Y: y},
				C0: Vec2{X: c0x, Y: c0y},
				C1: Vec2{X: c1x, Y: c1y},
				P1: Vec2{X: nx, Y: ny},
			}
			if !segNearZero(seg) {
				out = append(out, seg)
			}
			cx, cy = c1x, c1y
			x, y = nx, ny
			haveCtrl = true
			lastWasCubic, lastWasQ = true, false
		case 'Q', 'q':
			vals, err := nextNums(4)
			if err != nil {
				return nil, err
			}
			var qx, qy, nx, ny float32
			if cmd == 'q' {
				qx, qy = x+vals[0], y+vals[1]
				nx, ny = x+vals[2], y+vals[3]
			} else {
				qx, qy = vals[0], vals[1]
				nx, ny = vals[2], vals[3]
			}
			seg := quadToCubic(Vec2{X: x, Y: y}, Vec2{X: qx, Y: qy}, Vec2{X: nx, Y: ny})
			if !segNearZero(seg) {
				out = append(out, seg)
			}
			cx, cy = qx, qy
			x, y = nx, ny
			haveCtrl = true
			lastWasCubic, lastWasQ = false, true
		case 'T', 't':
			vals, err := nextNums(2)
			if err != nil {
				return nil, err
			}
			var qx, qy float32
			if lastWasQ && haveCtrl {
				qx, qy = 2*x-cx, 2*y-cy
			} else {
				qx, qy = x, y
			}
			var nx, ny float32
			if cmd == 't' {
				nx, ny = x+vals[0], y+vals[1]
			} else {
				nx, ny = vals[0], vals[1]
			}
			seg := quadToCubic(Vec2{X: x, Y: y}, Vec2{X: qx, Y: qy}, Vec2{X: nx, Y: ny})
			if !segNearZero(seg) {
				out = append(out, seg)
			}
			cx, cy = qx, qy
			x, y = nx, ny
			haveCtrl = true
			lastWasCubic, lastWasQ = false, true
		default:
			return nil, fmt.Errorf("path: unsupported command %c", cmd)
		}
	}
	return out, nil
}

func lineCubic(a, b Vec2) CubicBezier {
	d := b.Sub(a)
	return CubicBezier{
		P0: a,
		C0: a.Add(d.Scale(1.0 / 3)),
		C1: a.Add(d.Scale(2.0 / 3)),
		P1: b,
	}
}

func quadToCubic(p0, q, p1 Vec2) CubicBezier {
	// C0 = P0 + 2/3 (Q-P0); C1 = P1 + 2/3 (Q-P1)
	return CubicBezier{
		P0: p0,
		C0: p0.Add(q.Sub(p0).Scale(2.0 / 3)),
		C1: p1.Add(q.Sub(p1).Scale(2.0 / 3)),
		P1: p1,
	}
}

func segNearZero(b CubicBezier) bool {
	const eps = 1e-4
	d := b.P1.Sub(b.P0)
	if d.Dot(d) > eps*eps {
		return false
	}
	d0 := b.C0.Sub(b.P0)
	d1 := b.C1.Sub(b.P0)
	return d0.Dot(d0) <= eps*eps && d1.Dot(d1) <= eps*eps
}

func tokenizePath(d string) ([]string, error) {
	d = strings.TrimSpace(d)
	if d == "" {
		return nil, fmt.Errorf("empty path data")
	}
	var toks []string
	i := 0
	for i < len(d) {
		for i < len(d) && (d[i] == ' ' || d[i] == '\t' || d[i] == '\n' || d[i] == '\r' || d[i] == ',') {
			i++
		}
		if i >= len(d) {
			break
		}
		c := d[i]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			toks = append(toks, string(c))
			i++
			continue
		}
		start := i
		if d[i] == '+' || d[i] == '-' {
			i++
		}
		sawDot := false
		sawExp := false
		for i < len(d) {
			ch := d[i]
			if ch >= '0' && ch <= '9' {
				i++
				continue
			}
			if ch == '.' && !sawDot && !sawExp {
				sawDot = true
				i++
				continue
			}
			if (ch == 'e' || ch == 'E') && !sawExp {
				sawExp = true
				i++
				if i < len(d) && (d[i] == '+' || d[i] == '-') {
					i++
				}
				continue
			}
			break
		}
		if start == i || (i == start+1 && (d[start] == '+' || d[start] == '-')) {
			return nil, fmt.Errorf("path: bad number near %q", d[start:])
		}
		toks = append(toks, d[start:i])
		// SVG allows "1-2" without separator.
		if i < len(d) && (d[i] == '-' || (d[i] == '.' && i+1 < len(d) && unicode.IsDigit(rune(d[i+1])))) {
			continue
		}
	}
	return toks, nil
}
