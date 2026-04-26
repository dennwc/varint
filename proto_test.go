package varint

import (
	"encoding/binary"
	"math"
	"strconv"
	"testing"
)

func protoTag(buf []byte) (num int64, typ byte, n int) {
	var tag uint64
	tag, n = Uvarint(buf)
	if n == 0 || tag == 0 {
		return 0, 0, 0
	}
	typ, num = byte(tag&0x07), int64(tag>>3)
	if num <= 0 || typ > 5 {
		return 0, 0, n
	}
	return
}

var casesProtoTag = []int64{
	int64(MaxVal1 >> 3),
	int64(MaxVal2 >> 3),
	int64(MaxVal3 >> 3),
	int64(MaxVal4 >> 3),
	int64(MaxVal5 >> 3),
	int64(MaxVal6 >> 3),
	int64(MaxVal7 >> 3),
	int64(MaxVal8 >> 3),
	int64(MaxVal9 >> 3),
}

func putProtoTag(b []byte, num int64) int {
	// wire type 1, for test only
	return binary.PutUvarint(b[:], uint64(num)<<3|1)
}

func TestProtoTag(t *testing.T) {
	const wireType = 1
	for i, v := range casesProtoTag {
		sz := i + 1
		t.Run(strconv.Itoa(sz), func(t *testing.T) {
			var b [MaxLen64]byte
			n := putProtoTag(b[:], v)
			p := b[:n]
			if num, typ, n2 := protoTag(p); num != v || typ != wireType || n2 != n {
				t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
			}
			if num, typ, n2 := ProtoTag(p); num != v || typ != wireType || n2 != n {
				t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
			}
			v--
			n = putProtoTag(b[:], v)
			p = b[:n]
			if num, typ, n2 := protoTag(p); num != v || typ != wireType || n2 != n {
				t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
			}
			if num, typ, n2 := ProtoTag(p); num != v || typ != wireType || n2 != n {
				t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
			}
			t.Run("short", func(t *testing.T) {
				n := putProtoTag(b[:], v)
				p := b[:n-1]
				if _, _, n2 := protoTag(p); n2 != 0 {
					t.Fatalf("unexpected error: %d", n2)
				}
				if _, _, n2 := ProtoTag(p); n2 != 0 {
					t.Fatalf("unexpected error: %d", n2)
				}
			})
			t.Run("large", func(t *testing.T) {
				var b [16]byte
				n = putProtoTag(b[:], v)
				p := b[:]
				if num, typ, n2 := protoTag(p); num != v || typ != wireType || n2 != n {
					t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
				}
				if num, typ, n2 := ProtoTag(p); num != v || typ != wireType || n2 != n {
					t.Fatalf("%d,%d vs %d,%d", num, typ, v, wireType)
				}
			})
		})
	}
	t.Run("overflow", func(t *testing.T) {
		var b [MaxLen64 + 2]byte
		b[0] = 0xff
		b[1] = 0xff
		n := putProtoTag(b[2:], math.MaxInt64)
		p := b[:n+2]
		if _, _, n2 := protoTag(p); n2 != 0 {
			t.Fatalf("unexpected error: %d", n2)
		}
		if _, _, n2 := ProtoTag(p); n2 != 0 {
			t.Fatalf("unexpected error: %d", n2)
		}
	})
	t.Run("overflow_short", func(t *testing.T) {
		var b [MaxLen64 + 2]byte
		b[0] = 0xff
		b[1] = 0xff
		n := putProtoTag(b[2:], math.MaxInt64)
		p := b[:n+2]
		p[len(p)-1] = 0xff
		if _, _, n2 := protoTag(p); n2 != 0 {
			t.Fatalf("unexpected error: %d", n2)
		}
		if _, _, n2 := ProtoTag(p); n2 != 0 {
			t.Fatalf("unexpected error: %d", n2)
		}
	})
	t.Run("empty", func(t *testing.T) {
		if _, _, n := ProtoTag(nil); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("invalid_wire", func(t *testing.T) {
		if _, _, n := ProtoTag([]byte{0x06}); n != 0 { // type 6
			t.Fatalf("expected 0, got %d", n)
		}
		if _, _, n := ProtoTag([]byte{0x86, 0x01}); n != 0 { // type 6, multi-byte
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("zero_field", func(t *testing.T) {
		if _, _, n := ProtoTag([]byte{0x00}); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
		if _, _, n := ProtoTag([]byte{0x01}); n != 0 { // field 0, type 1
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("invalid_wire_all", func(t *testing.T) {
		for typ := byte(6); typ <= 7; typ++ {
			if _, _, n := ProtoTag([]byte{typ}); n != 0 {
				t.Fatalf("expected 0 for type %d, got %d", typ, n)
			}
			if _, _, n := ProtoTag([]byte{0x80 | typ, 0x01}); n != 0 {
				t.Fatalf("expected 0 for type %d (multi), got %d", typ, n)
			}
		}
	})
	t.Run("sz9_exactly_short", func(t *testing.T) {
		buf := []byte{0x88, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
		if _, _, n := ProtoTag(buf); n != 0 {
			t.Fatalf("expected 0, got %d", n)
		}
	})
	t.Run("sz9_all", func(t *testing.T) {
		for i := 1; i <= 9; i++ {
			var buf [9]byte
			if i == 1 {
				buf[0] = 0x08 // field 1, wire 0
			} else {
				buf[0] = 0x88 // field 1, wire 0, more bytes
				for j := 1; j < i-1; j++ {
					buf[j] = 0x80
				}
				buf[i-1] = 0x00 // terminator
			}

			got, typ, n := ProtoTag(buf[:])
			if n != i {
				t.Fatalf("expected size %d, got %d for step %d", i, n, i)
			}
			if got != 1 || typ != 0 {
				t.Fatalf("expected field 1, typ 0; got %d, %d for step %d", got, typ, i)
			}
		}
	})
}

func BenchmarkProtoTagSimple(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesProtoTag[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var buf [MaxLen64]byte
			_ = putProtoTag(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				num int64
				typ byte
				n   int
			)
			for i := 0; i < b.N; i++ {
				num, typ, n = protoTag(p)
			}
			if num != v || typ != 1 {
				b.Fatal(num)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
		b.Run(strconv.Itoa(sz)+"_large", func(b *testing.B) {
			var buf [16]byte
			_ = putProtoTag(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				num int64
				typ byte
				n   int
			)
			for i := 0; i < b.N; i++ {
				num, typ, n = protoTag(p)
			}
			if num != v || typ != 1 {
				b.Fatal(num)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
	}
	b.Run("zeros", func(b *testing.B) {
		var buf [16]byte
		p := buf[:]

		b.ResetTimer()
		var (
			num int64
			typ byte
			n   int
		)
		for i := 0; i < b.N; i++ {
			num, typ, n = protoTag(p)
		}
		_, _ = num, typ
		if n != 0 {
			b.Fatal(n)
		}
	})
}

func BenchmarkProtoTag(b *testing.B) {
	b.StopTimer()
	for sz := 1; sz <= 9; sz++ {
		v := casesProtoTag[sz-1]
		b.Run(strconv.Itoa(sz), func(b *testing.B) {
			var buf [MaxLen64]byte
			_ = putProtoTag(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				num int64
				typ byte
				n   int
			)
			for i := 0; i < b.N; i++ {
				num, typ, n = ProtoTag(p)
			}
			if num != v || typ != 1 {
				b.Fatal(num)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
		b.Run(strconv.Itoa(sz)+"_large", func(b *testing.B) {
			var buf [16]byte
			_ = putProtoTag(buf[:], v)
			p := buf[:sz]

			b.ResetTimer()
			var (
				num int64
				typ byte
				n   int
			)
			for i := 0; i < b.N; i++ {
				num, typ, n = ProtoTag(p)
			}
			if num != v || typ != 1 {
				b.Fatal(num)
			}
			if n != sz {
				b.Fatal(n)
			}
		})
	}
	b.Run("zeros", func(b *testing.B) {
		var buf [16]byte
		p := buf[:]

		b.ResetTimer()
		var (
			num int64
			typ byte
			n   int
		)
		for i := 0; i < b.N; i++ {
			num, typ, n = ProtoTag(p)
		}
		_, _ = num, typ
		if n != 0 {
			b.Fatal(n)
		}
	})
}
