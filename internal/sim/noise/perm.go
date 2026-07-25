package noise

// permTable is a duplicated 256-entry permutation for Perlin/Simplex hashing.
type permTable [512]uint8

func newPerm(seed int64) permTable {
	var p [256]uint8
	for i := 0; i < 256; i++ {
		p[i] = uint8(i)
	}
	// SplitMix64-style shuffle for deterministic, seed-sensitive order.
	s := uint64(seed)
	if s == 0 {
		s = 0x9E3779B97F4A7C15
	}
	for i := 255; i > 0; i-- {
		s = splitMix64(s)
		j := int(s % uint64(i+1))
		p[i], p[j] = p[j], p[i]
	}
	var out permTable
	for i := 0; i < 256; i++ {
		out[i] = p[i]
		out[i+256] = p[i]
	}
	return out
}

func splitMix64(x uint64) uint64 {
	x += 0x9E3779B97F4A7C15
	z := x
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

// hash2/hash3 return 0..255 from lattice coords using the permutation.
func (p permTable) hash2(ix, iy int) int {
	return int(p[uint8(ix)+p[uint8(iy)]])
}

func (p permTable) hash3(ix, iy, iz int) int {
	return int(p[uint8(ix)+p[uint8(iy)+p[uint8(iz)]]])
}

// cellSeed hashes integer cell coords with an independent seed (Worley).
func cellSeed(seed int64, ix, iy, iz int) uint64 {
	x := uint64(seed) ^ (uint64(uint32(ix)) * 0x85EBCA77C2B2AE63)
	x ^= uint64(uint32(iy)) * 0xC2B2AE3D27D4EB4F
	x ^= uint64(uint32(iz)) * 0x165667B19E3779F9
	return splitMix64(x)
}

func rand01(s uint64) (float64, uint64) {
	s = splitMix64(s)
	// top 53 bits → [0,1)
	return float64(s>>11) / (1 << 53), s
}
