# varint

[![](https://godoc.org/github.com/dennwc/varint?status.svg)](http://godoc.org/github.com/dennwc/varint)

This package provides an optimized implementation of protobuf's varint encoding/decoding.
It has no dependencies.

Benchmarks comparing to a `binary.Uvarint`:

```
benchmark                      old ns/op     new ns/op     delta
BenchmarkUvarint/1-8           3.46          2.64          -23.73%
BenchmarkUvarint/1_large-8     3.48          2.31          -33.73%
BenchmarkUvarint/2-8           5.51          3.20          -41.90%
BenchmarkUvarint/2_large-8     5.45          2.32          -57.43%
BenchmarkUvarint/3-8           7.18          3.46          -51.79%
BenchmarkUvarint/3_large-8     7.00          3.04          -56.52%
BenchmarkUvarint/4-8           8.04          3.54          -55.97%
BenchmarkUvarint/4_large-8     8.34          3.19          -61.70%
BenchmarkUvarint/5-8           9.58          4.03          -57.97%
BenchmarkUvarint/5_large-8     9.77          3.75          -61.58%
BenchmarkUvarint/6-8           10.9          4.34          -60.28%
BenchmarkUvarint/6_large-8     10.9          4.31          -60.55%
BenchmarkUvarint/7-8           12.3          4.63          -62.53%
BenchmarkUvarint/7_large-8     12.4          5.05          -59.25%
BenchmarkUvarint/8-8           13.8          5.48          -60.30%
BenchmarkUvarint/8_large-8     14.1          4.89          -65.43%
BenchmarkUvarint/9-8           15.3          5.65          -63.21%
BenchmarkUvarint/9_large-8     15.3          5.47          -64.15%
```

It also provides additional functionality like `UvarintSize` (similar to `sov*` in `gogo/protobuf`):

```
benchmark                    old ns/op     new ns/op     delta
BenchmarkUvarintSize/1-8     1.71          0.43          -74.85%
BenchmarkUvarintSize/2-8     2.56          0.57          -77.73%
BenchmarkUvarintSize/3-8     3.22          0.72          -77.64%
BenchmarkUvarintSize/4-8     3.74          0.72          -80.75%
BenchmarkUvarintSize/5-8     4.29          0.57          -86.71%
BenchmarkUvarintSize/6-8     4.85          0.58          -88.04%
BenchmarkUvarintSize/7-8     5.43          0.71          -86.92%
BenchmarkUvarintSize/8-8     6.01          0.86          -85.69%
BenchmarkUvarintSize/9-8     6.64          1.00          -84.94%
```

# License

MIT