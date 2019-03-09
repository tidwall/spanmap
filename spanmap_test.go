package spanmap

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/tidwall/celltree"
	"github.com/tidwall/lotsa"

	"github.com/google/btree"
)

func init() {
	seed := time.Now().UnixNano()
	println("seed:", seed)
	rand.Seed(seed)
}

func TestMap(t *testing.T) {
	var m Map
	if m.Len() != 0 {
		t.Fatalf("expected 0, got '%v'", m.Len())
	}
	// empty map operations
	if m.Get(100) != nil {
		t.Fatalf("expected nil, got '%v'", m.Get(100))
	}
	if m.Delete(100) != nil {
		t.Fatalf("expected nil, got '%v'", m.Delete(100))
	}
	func() {
		defer func() {
			v := recover().(string)
			if v != "nil item" {
				t.Fatalf("expected panic")
			}
		}()
		m.Set(100, nil)
	}()

	// insert some data in random order
	N := 10000
	for _, i := range rand.Perm(N) {
		value := m.Set(uint64(i), i)
		if value != nil {
			t.Fatalf("expected nil, got '%v'", value)
		}
	}
	if m.Len() != N {
		t.Fatalf("expected '%v', got '%v'", N, m.Len())
	}
	// check the min/max
	if m.Min() != 0 || m.Max() != uint64(N-1) {
		t.Fatalf("expected '%v/%v', got '%v/%v'", 0, N-1, m.Min(), m.Max())
	}
	for i := 0; i < N; i += 2 {
		value := m.Delete(uint64(i)).(int)
		if value != i {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}
	if m.Len() != N/2 {
		t.Fatalf("expected '%v', got '%v'", N/2, m.Len())
	}
	if m.Min() != 1 || m.Max() != uint64(N-1) {
		t.Fatalf("expected '%v/%v', got '%v/%v'", 1, N-1, m.Min(), m.Max())
	}
	for i := 0; i < N; i++ {
		v := m.Get(uint64(i))
		if i%2 == 0 {
			if v != nil {
				t.Fatalf("expected nil, got '%v'", v)
			}
		} else {
			if v.(int) != i {
				t.Fatalf("expected '%v', got '%v'", i, v.(int))
			}
		}
	}

	// change the type
	for i := 1; i < N; i += 2 {
		value := m.Set(uint64(i), uint64(i))
		if value != i {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}
	for i := 1; i < N; i += 2 {
		value := m.Get(uint64(i)).(uint64)
		if value != uint64(i) {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}

	// delete each item one by one
	n := m.Len()
	for i := 0; i < N; i++ {
		v := m.Delete(uint64(i))
		if i%2 == 0 {
			if v != nil {
				t.Fatalf("expected nil, got '%v'", v)
			}
		} else {
			if v.(uint64) != uint64(i) {
				t.Fatalf("expected '%v', got '%v'", i, v.(uint64))
			}
			if i == N-1 {
				if m.Min() != 0 || m.Max() != 0 {
					t.Fatalf("expected '%v/%v', got '%v/%v'",
						0, 0, m.Min(), m.Max())
				}
			} else {
				n--
				if m.Len() != n {
					t.Fatalf("expected '%v', got '%v'", n, m.Len())
				}
				// check the entire list again for correctness
				for i := 0; i < N; i++ {
					v := m.Get(uint64(i))
					if i%2 == 0 || uint64(i) < m.Min() {
						if v != nil {
							t.Fatalf("expected nil, got '%v'", v)
						}
					} else {
						if v.(uint64) != uint64(i) {
							t.Fatalf("expected '%v', got '%v'", i, v.(uint64))
						}
					}
				}

				if m.Min() != uint64(i)+2 || m.Max() != uint64(N-1) {
					t.Fatalf("expected '%v/%v', got '%v/%v'",
						uint64(i)+2, N-1, m.Min(), m.Max())
				}
			}
		}
	}
	if m.Min() != 0 || m.Max() != 0 {
		t.Fatalf("expected '%v/%v', got '%v/%v'",
			0, 0, m.Min(), m.Max())
	}
	if m.Len() != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, m.Len())
	}
	// insert and reverse deleted
	for _, i := range rand.Perm(N) {
		v := m.Set(uint64(i), i)
		if v != nil {
			t.Fatalf("expected nil, got '%v'", v)
		}
	}
	for i := N - 1; i >= 0; i -= 2 {
		value := m.Delete(uint64(i)).(int)
		if value != i {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}
	if m.Min() != 0 || m.Max() != uint64(N-2) {
		t.Fatalf("expected '%v/%v', got '%v/%v'", 0, N-2, m.Min(), m.Max())
	}

	for i := N - 1; i >= 0; i-- {
		v := m.Get(uint64(i))
		if i%2 == 1 {
			if v != nil {
				t.Fatalf("expected nil, got '%v'", v)
			}
		} else {
			if v.(int) != i {
				t.Fatalf("expected '%v', got '%v'", i, v.(int))
			}
		}
	}

	// change the type
	for i := N - 2; i >= 0; i -= 2 {
		value := m.Set(uint64(i), uint64(i))
		if value != i {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}
	for i := N - 2; i >= 0; i -= 2 {
		value := m.Get(uint64(i)).(uint64)
		if value != uint64(i) {
			t.Fatalf("expected '%v', got '%v'", i, value)
		}
	}

	n = m.Len()
	// delete each item one by one
	for i := N - 1; i >= 0; i-- {
		v := m.Delete(uint64(i))
		if i%2 == 1 {
			if v != nil {
				t.Fatalf("expected nil, got '%v'", v)
			}
		} else {
			if v.(uint64) != uint64(i) {
				t.Fatalf("expected '%v', got '%v'", i, v.(uint64))
			}
			if i == 0 {
				if m.Min() != 0 || m.Max() != 0 {
					t.Fatalf("expected '%v/%v', got '%v/%v'",
						0, 0, m.Min(), m.Max())
				}
			} else {
				n--
				if m.Len() != n {
					t.Fatalf("expected '%v', got '%v'", n, m.Len())
				}
				// check the entire list again for correctness
				for i := N - 1; i >= 0; i-- {
					v := m.Get(uint64(i))
					if i%2 == 1 || uint64(i) > m.Max() {
						if v != nil {
							t.Fatalf("expected nil, got '%v'", v)
						}
					} else {
						if v.(uint64) != uint64(i) {
							t.Fatalf("expected '%v', got '%v'", i, v.(uint64))
						}
					}
				}
				if m.Min() != 0 || m.Max() != uint64(i-2) {
					t.Fatalf("expected '%v/%v', got '%v/%v'",
						0, i-3, m.Min(), m.Max())
				}
			}
		}
	}
	if m.Min() != 0 || m.Max() != 0 {
		t.Fatalf("expected '%v/%v', got '%v/%v'",
			0, 0, m.Min(), m.Max())
	}
	if m.Len() != 0 {
		t.Fatalf("expected '%v', got '%v'", 0, m.Len())
	}

}

type btreeItem uint64

func (item btreeItem) Less(other btree.Item) bool {
	return uint64(item) < uint64(other.(btreeItem))
}

func BenchmarkCompare(b *testing.B) {
	defer b.Skip()
	N := 1000000
	ints := make([]uint64, 0, N)
	for _, i := range rand.Perm(N) {
		ints = append(ints, uint64(i))
	}
	println("\n-- spanmap --")
	lotsa.Output = os.Stdout
	{
		var m Map
		print("set/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			m.Set(uint64(i), uint64(i))
		})
		print("get/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if m.Get(uint64(i)).(uint64) != uint64(i) {
				panic("bad news")
			}
		})
		m = Map{}
		print("set/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			m.Set(ints[i], ints[i])
		})
		print("get/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if m.Get(ints[i]).(uint64) != uint64(ints[i]) {
				panic("bad news")
			}
		})
	}
	println("\n-- stdlib map --")
	{
		m := make(map[uint64]interface{})
		print("set/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			m[uint64(i)] = uint64(i)
		})
		print("get/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if m[uint64(i)].(uint64) != uint64(i) {
				panic("bad news")
			}
		})
		m = make(map[uint64]interface{})
		print("set/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			m[ints[i]] = ints[i]
		})
		print("get/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if m[ints[i]].(uint64) != ints[i] {
				panic("bad news")
			}
		})
	}
	println("\n-- btree --")
	{
		tr := btree.New(32)
		print("set/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.ReplaceOrInsert(btreeItem(uint64(i)))
		})
		print("get/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if uint64(tr.Get(btreeItem(uint64(i))).(btreeItem)) != uint64(i) {
				panic("bad news")
			}
		})
		tr = btree.New(32)
		print("set/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.ReplaceOrInsert(btreeItem(ints[i]))
		})
		print("get/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			if uint64(tr.Get(btreeItem(ints[i])).(btreeItem)) != ints[i] {
				panic("bad news")
			}
		})
	}
	println("\n-- celltree --")
	lotsa.Output = os.Stdout
	{
		var tr celltree.Tree
		print("set/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.Insert(uint64(i), uint64(i))
		})
		print("get/sequential ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.Range(uint64(i), func(cell uint64, v interface{}) bool {
				if v.(uint64) != uint64(i) {
					panic("bad news")
				}
				return false
			})
		})
		tr = celltree.Tree{}
		print("set/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.Insert(ints[i], ints[i])
		})
		print("get/random     ")
		lotsa.Ops(N, 1, func(i, _ int) {
			tr.Range(ints[i], func(cell uint64, v interface{}) bool {
				if v.(uint64) != ints[i] {
					panic("bad news")
				}
				return false
			})
		})
	}

}
