package varint

import (
	"encoding/binary"
	"strconv"
	"testing"
)

var casesUvarintSize = []uint64{
	MaxVal1,
	MaxVal2,
	MaxVal3,
	MaxVal4,
	MaxVal5,
	MaxVal6,
	MaxVal7,
	MaxVal8,
	MaxVal9,
	maxUint64,
}

func TestUvarintSize(t *testing.T) {
	for i, v := range casesUvarintSize {
		sz := i + 1
		t.Run(strconv.Itoa(sz), func(t *testing.T) {
			if got := UvarintSize(v); got != sz {
				t.Fatalf("%d vs %d", got, sz)
			}
			if got := UvarintSize(v - 1); got != sz {
				t.Fatalf("%d vs %d", got, sz)
			}
			if got := uvarintSizeLoop(v); got != sz {
				t.Fatalf("%d vs %d", got, sz)
			}
			if got := uvarintSizeFlat(v); got != sz {
				t.Fatalf("%d vs %d", got, sz)
			}
		})
	}
}

func TestUvarint(t *testing.T) {
	for i, v := range casesUvarintSize {
		sz := i + 1
		t.Run(strconv.Itoa(sz), func(t *testing.T) {
			var b [MaxLen64]byte
			n := binary.PutUvarint(b[:], v)
			p := b[:n]
			if got, n2 := binary.Uvarint(p); got != v || n2 != n {
				t.Fatalf("%d vs %d", got, v)
			}
			if got, n2 := Uvarint(p); got != v || n2 != n {
				t.Fatalf("%d vs %d", got, v)
			}
			v--
			n = binary.PutUvarint(b[:], v)
			p = b[:n]
			if got, n2 := binary.Uvarint(p); got != v || n2 != n {
				t.Fatalf("%d vs %d", got, v)
			}
			if got, n2 := Uvarint(p); got != v || n2 != n {
				t.Fatalf("%d vs %d", got, v)
			}
			t.Run("short", func(t *testing.T) {
				n := binary.PutUvarint(b[:], v)
				p := b[:n-1]
				if _, n2 := binary.Uvarint(p); n2 != 0 {
					t.Fatalf("unexpected error: %d", n2)
				}
				if _, n2 := Uvarint(p); n2 != 0 {
					t.Fatalf("unexpected error: %d", n2)
				}
			})
			t.Run("large", func(t *testing.T) {
				var b [16]byte
				n = binary.PutUvarint(b[:], v)
				p := b[:]
				if got, n2 := binary.Uvarint(p); got != v || n2 != n {
					t.Fatalf("%d vs %d", got, v)
				}
				if got, n2 := Uvarint(p); got != v || n2 != n {
					t.Fatalf("%d vs %d", got, v)
				}
			})
		})
	}
	t.Run("overflow", func(t *testing.T) {
		var b [MaxLen64 + 2]byte
		b[0] = 0xff
		b[1] = 0xff
		n := binary.PutUvarint(b[2:], maxUint64)
		p := b[:n+2]
		if _, n2 := binary.Uvarint(p); n2 != -11 {
			t.Fatalf("unexpected error: %d", n2)
		}
		if _, n2 := Uvarint(p); n2 != -11 {
			t.Fatalf("unexpected error: %d", n2)
		}
	})
	t.Run("overflow_short", func(t *testing.T) {
		var b [MaxLen64 + 2]byte
		b[0] = 0xff
		b[1] = 0xff
		n := binary.PutUvarint(b[2:], maxUint64)
		p := b[:n+2]
		p[len(p)-1] = 0xff
		if _, n2 := binary.Uvarint(p); n2 != -11 {
			t.Fatalf("unexpected error: %d", n2)
		}
		if _, n2 := Uvarint(p); n2 != -11 {
			t.Fatalf("unexpected error: %d", n2)
		}
	})
	t.Run("empty", func(t *testing.T) {
		if _, n := Uvarint(nil); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
		if _, n := Uvarint([]byte{}); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("sz10_short", func(t *testing.T) {
		buf := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		if _, n := Uvarint(buf); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("short_all", func(t *testing.T) {
		buf := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		for i := 1; i <= 9; i++ {
			if _, n := Uvarint(buf[:i]); n != 0 {
				t.Fatalf("expected 0 for size %d, got %d", i, n)
			}
		}
	})
	t.Run("sz10_all", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			var buf [10]byte
			for j := 0; j < i-1; j++ {
				buf[j] = 0x80
			}
			buf[i-1] = 0x01
			// This will have sz >= 10 (it is exactly 10)
			// It should cover the unrolled loop in the `if sz >= 10` block
			got, n := Uvarint(buf[:])
			if n != i {
				t.Fatalf("expected size %d, got %d", i, n)
			}
			if i < 10 && got != 1<<(7*uint(i-1)) {
				t.Fatalf("expected value %d, got %d", 1<<(7*uint(i-1)), got)
			}
		}
	})
	t.Run("exhaustive", func(t *testing.T) {
		for i := 1; i <= 10; i++ {
			// value of length i
			var v uint64
			if i == 10 {
				v = maxUint64
			} else {
				v = 1 << (7 * uint(i-1))
			}
			var buf [10]byte
			n := binary.PutUvarint(buf[:], v)
			if n != i {
				t.Fatalf("expected size %d, got %d", i, n)
			}

			for sz := 0; sz <= 10; sz++ {
				got, n2 := Uvarint(buf[:sz])
				if sz < i {
					if n2 != 0 {
						t.Fatalf("v size %d, buf size %d: expected 0, got %d", i, sz, n2)
					}
				} else {
					if n2 != i {
						t.Fatalf("v size %d, buf size %d: expected %d, got %d", i, sz, i, n2)
					}
					if got != v {
						t.Fatalf("v size %d, buf size %d: expected %v, got %v", i, sz, v, got)
					}
				}
			}
		}
	})
	t.Run("sz10_overflow", func(t *testing.T) {
		buf := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x02}
		if _, n := Uvarint(buf); n != -10 {
			t.Fatalf("expected -10, got %d", n)
		}
	})
	t.Run("sz10_exactly_short", func(t *testing.T) {
		buf := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		if _, n := Uvarint(buf); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("sz12_overflow", func(t *testing.T) {
		buf := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		if _, n := Uvarint(buf); n != -11 {
			t.Fatalf("expected -11, got %d", n)
		}
	})
}

func uvarintSizeLoop(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}

func uvarintSizeFlat(x uint64) int {
	if x <= MaxVal1 {
		return 1
	} else if x <= MaxVal2 {
		return 2
	} else if x <= MaxVal3 {
		return 3
	} else if x <= MaxVal4 {
		return 4
	} else if x <= MaxVal5 {
		return 5
	} else if x <= MaxVal6 {
		return 6
	} else if x <= MaxVal7 {
		return 7
	} else if x <= MaxVal8 {
		return 8
	} else if x <= MaxVal9 {
		return 9
	}
	return 10
}

func BenchmarkUvarintSizeLoop(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesUvarintSize[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var n int
			for i := 0; i < b.N; i++ {
				n = uvarintSizeLoop(v)
			}
			_ = n
		})
	}
}

func BenchmarkUvarintSizeFlat(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesUvarintSize[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var n int
			for i := 0; i < b.N; i++ {
				n = uvarintSizeFlat(v)
			}
			_ = n
		})
	}
}

func BenchmarkUvarintSize(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesUvarintSize[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var n int
			for i := 0; i < b.N; i++ {
				n = UvarintSize(v)
			}
			_ = n
		})
	}
}

func BenchmarkUvarintStdlib(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesUvarintSize[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var buf [MaxLen64]byte
			_ = binary.PutUvarint(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				val uint64
				n   int
			)
			for i := 0; i < b.N; i++ {
				val, n = binary.Uvarint(p)
			}
			if val != v {
				b.Fatal(val)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
		b.Run(strconv.Itoa(sz)+"_large", func(b *testing.B) {
			var buf [16]byte
			_ = binary.PutUvarint(buf[:], v)
			p := buf[:]

			b.ResetTimer()
			var (
				val uint64
				n   int
			)
			for i := 0; i < b.N; i++ {
				val, n = binary.Uvarint(p)
			}
			if val != v {
				b.Fatal(val)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
	}
}

func BenchmarkUvarint(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesUvarintSize[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var buf [MaxLen64]byte
			_ = binary.PutUvarint(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				val uint64
				n   int
			)
			for i := 0; i < b.N; i++ {
				val, n = Uvarint(p)
			}
			if val != v {
				b.Fatal(val)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
		b.Run(strconv.Itoa(sz)+"_large", func(b *testing.B) {
			var buf [16]byte
			_ = binary.PutUvarint(buf[:], v)
			p := buf[:]

			b.ResetTimer()
			var (
				val uint64
				n   int
			)
			for i := 0; i < b.N; i++ {
				val, n = Uvarint(p)
			}
			if val != v {
				b.Fatal(val)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
	}
}
