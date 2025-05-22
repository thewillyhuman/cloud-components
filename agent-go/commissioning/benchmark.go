package commissioning

import (
	"agent-go/agent"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"runtime"
	"time"
)

type benchmark struct {
	name string
	fn   func() float64
	runs int
}

func RunBenchmarks() map[string]float64 {
	agent.SetStatus(agent.TestingStatus)
	runtime.GOMAXPROCS(runtime.NumCPU())

	benchmarks := []benchmark{
		{"cpu_prime", benchPrime, 10000},
		{"cpu_fib", benchFib, 100},
		{"mem_alloc", benchMemAlloc, 1000},
		{"mem_fill", benchMemFill, 10},
		{"json_encode", benchJSONEncode, 1},
		{"json_decode", benchJSONDecode, 1},
		{"sha256_hash", benchSHA256, 10},
		{"gzip_compress", benchGzipCompress, 1},
		{"gzip_decompress", benchGzipDecompress, 1},
	}

	log.Println("Starting system benchmarks...")
	results := make(map[string]float64)

	for _, bm := range benchmarks {
		log.Printf("Running %-16s (%4d runs)...", bm.name, bm.runs)
		start := time.Now()
		score := runBenchmark(bm.fn, bm.runs)
		elapsed := time.Since(start)
		results[bm.name] = score
		log.Printf("Finished %-16s → Score: %10.2f | Time: %s", bm.name, score, elapsed.Truncate(time.Millisecond))
	}

	return results
}

func runBenchmark(fn func() float64, runs int) float64 {
	var total float64
	for i := 0; i < runs; i++ {
		total += fn()
	}
	return total / float64(runs)
}

// CPU: Prime Numbers
func benchPrime() float64 {
	start := time.Now()
	count := 0
	for n := 2; n < 100_000; n++ {
		if isPrime(n) {
			count++
		}
	}
	return float64(count) / time.Since(start).Seconds()
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// CPU: Fibonacci
func benchFib() float64 {
	start := time.Now()
	count := 0
	for i := 25; i < 35; i++ {
		_ = fib(i)
		count++
	}
	return float64(count) / time.Since(start).Seconds()
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

// Memory: Allocation
func benchMemAlloc() float64 {
	start := time.Now()
	for i := 0; i < 10_000_000; i++ {
		_ = make([]byte, 64)
	}
	return 10_000_000 / time.Since(start).Seconds()
}

// Memory: Fill
func benchMemFill() float64 {
	const sizeMB = 512
	buf := make([]byte, sizeMB*1024*1024)
	start := time.Now()
	for i := range buf {
		buf[i] = byte(i % 256)
	}
	return float64(sizeMB) / time.Since(start).Seconds()
}

// JSON Encoding
func benchJSONEncode() float64 {
	data := make([]map[string]interface{}, 10000)
	for i := range data {
		data[i] = map[string]interface{}{
			"id":   i,
			"name": fmt.Sprintf("node-%d", i),
			"tags": []string{"test", "benchmark"},
		}
	}
	start := time.Now()
	for i := 0; i < 1000; i++ {
		_, _ = json.Marshal(data)
	}
	return 1000 / time.Since(start).Seconds()
}

// JSON Decoding
func benchJSONDecode() float64 {
	sample := `{"id":123,"name":"benchmark","tags":["go","test"]}`
	start := time.Now()
	for i := 0; i < 1_000_000; i++ {
		var result map[string]interface{}
		_ = json.Unmarshal([]byte(sample), &result)
	}
	return 1_000_000 / time.Since(start).Seconds()
}

// Hashing: SHA256
func benchSHA256() float64 {
	const sizeMB = 256
	data := make([]byte, sizeMB*1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}
	start := time.Now()
	_ = sha256.Sum256(data)
	return float64(sizeMB) / time.Since(start).Seconds()
}

// Compression: GZIP
func benchGzipCompress() float64 {
	const sizeMB = 128
	data := make([]byte, sizeMB*1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	var buf bytes.Buffer
	start := time.Now()
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write(data)
	zw.Close()
	return float64(sizeMB) / time.Since(start).Seconds()
}

// Decompression: GZIP
func benchGzipDecompress() float64 {
	const sizeMB = 128
	data := make([]byte, sizeMB*1024*1024)
	for i := range data {
		data[i] = byte(i % 256)
	}

	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, _ = zw.Write(data)
	zw.Close()

	start := time.Now()
	zr, _ := gzip.NewReader(&buf)
	io.Copy(io.Discard, zr)
	zr.Close()
	return float64(sizeMB) / time.Since(start).Seconds()
}
