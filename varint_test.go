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
