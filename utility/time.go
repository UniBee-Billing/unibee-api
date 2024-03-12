package utility

func MaxInt64(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxUInt64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func MinInt64(a int64, b int64) int64 {
	if a > b {
		return b
	}
	return a
}
