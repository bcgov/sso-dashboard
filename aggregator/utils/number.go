package utils

func NegateInt(num int64) int64 {
	return int64(^uint64(int64(num) - 1))
}
