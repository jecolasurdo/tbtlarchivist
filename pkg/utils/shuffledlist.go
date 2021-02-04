package utils

import (
	"math/rand"
	"time"
)

// GetShuffledIntList returns a list of contiguous sudo-randomly shuffled
// integers n, where n[i] is in [1, size].
func GetShuffledIntList(size int) []int {
	nn := make([]int, size)
	for i := 0; i < size; i++ {
		nn[i] = i + 1
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(size, func(i, j int) {
		nn[i], nn[j] = nn[j], nn[i]
	})
	return nn
}
