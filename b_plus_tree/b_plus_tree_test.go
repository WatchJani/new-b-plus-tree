package b_plus_tree

import (
	"fmt"
	"math/rand"
	"testing"
)

// is insert function work well
func TestInsert(t *testing.T) {
	numOfOperation := 20_000

	bPTree := NewBPTree[int, int](numOfOperation, 50)

	for index := 0; index < numOfOperation; index++ {
		bPTree.Insert(rand.Intn(2500), rand.Intn(1024))
	}

	get, want := bPTree.All()
	if get != want {
		t.Errorf("get: %d | want: %d", get, want)
	}

}

func BenchmarkInsertSameValue(b *testing.B) {
	b.StopTimer()
	BPTree := NewBPTree[int, int](40_000, 5)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		BPTree.Insert(5, 0)
	}
}

// 281.2 ns/op
func BenchmarkInsertRandom(b *testing.B) {
	b.StopTimer()

	data := make([]string, b.N)

	for index := range data {
		data[index] = fmt.Sprintf("%d", rand.Intn(500_000))
	}

	BPTree := NewBPTree[string, int](40_000, 50)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		BPTree.Insert(data[i], rand.Intn(500000))
	}
}

func Benchmark(b *testing.B) {
	b.StopTimer()
	data := make([]string, 40_001)

	for index := range data {
		data[index] = fmt.Sprintf("%d", rand.Intn(500_000))
	}

	BPTree := NewBPTree[string, int](40_000_000, 50)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20_000; j++ {
			BPTree.Insert(data[j], 0)
		}
		BPTree.ClearTree()
	}
}
