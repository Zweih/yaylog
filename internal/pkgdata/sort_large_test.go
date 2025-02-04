package pkgdata

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"
	"time"
)

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func generateDatasetSizes(start, stop, numSteps int) []int {
	sizes := make([]int, numSteps)
	logMin := math.Log10(float64(start))
	logMax := math.Log10(float64(stop))
	logStep := (logMax - logMin) / float64(numSteps-1)

	for i := 0; i < numSteps; i++ {
		sizes[i] = int(math.Pow(10, logMin+float64(i)*logStep))
	}
	return sizes
}

// Generate dataset dynamically with more steps (20-30 sizes)
var (
	datasetSizes   = generateDatasetSizes(50, 3000, 5) // Increased dataset sizes
	cpuCoreOptions = []int{1, 2, runtime.NumCPU()}     // Different CPU core counts
)

// Generate a dataset of given size
func generateDataset(size int) []PackageInfo {
	pkgs := make([]PackageInfo, size)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < size; i++ {
		pkgs[i] = PackageInfo{
			Name:      randomString(10),
			Timestamp: time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour),
			Size:      rand.Int63n(10 * GB),
		}
	}
	return pkgs
}

// Helper function to generate random package names
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Run benchmark for normal and concurrent sorting at different CPU core settings
func benchmarkSortingWithCores(b *testing.B, datasetSize int, sortType string, comparator PackageComparator) {
	runs := 5 // Number of repetitions for each test
	resultsNormal := make([]time.Duration, runs)

	largePkgList := generateDataset(datasetSize)

	// Benchmark Normal Sorting
	for i := 0; i < runs; i++ {
		dataCopy := make([]PackageInfo, len(largePkgList))
		copy(dataCopy, largePkgList)

		start := time.Now()
		_ = sortNormally(dataCopy, comparator, fmt.Sprintf("Sorting %s Normally", sortType), nil)
		resultsNormal[i] = time.Since(start)
	}

	avgNormal, fastestNormal, slowestNormal := computeStats(resultsNormal)

	// Print Normal Sorting Stats
	fmt.Printf("\nSorting Benchmark for %s (%d packages):\n", sortType, datasetSize)
	fmt.Printf("  - [Normal] Avg: %v | Fastest: %v | Slowest: %v\n", avgNormal, fastestNormal, slowestNormal)

	// Benchmark Concurrent Sorting with Different CPU Cores
	for _, cores := range cpuCoreOptions {
		runtime.GOMAXPROCS(cores) // Set the number of CPU cores
		resultsConcurrent := make([]time.Duration, runs)

		for i := 0; i < runs; i++ {
			dataCopy := make([]PackageInfo, len(largePkgList))
			copy(dataCopy, largePkgList)

			start := time.Now()
			_ = sortConcurrently(dataCopy, comparator, fmt.Sprintf("Sorting %s Concurrently (%d cores)", sortType, cores), nil)
			resultsConcurrent[i] = time.Since(start)
		}

		avgConcurrent, fastestConcurrent, slowestConcurrent := computeStats(resultsConcurrent)
		efficiencyRatio := float64(avgNormal) / float64(avgConcurrent)

		// Print Concurrent Sorting Stats
		fmt.Printf("  - [Concurrent (%d cores)] Avg: %v | Fastest: %v | Slowest: %v | Efficiency: %.2fx\n",
			cores, avgConcurrent, fastestConcurrent, slowestConcurrent, efficiencyRatio)
	}
}

// Compute statistics for benchmark results
func computeStats(times []time.Duration) (avg time.Duration, fastest time.Duration, slowest time.Duration) {
	total := time.Duration(0)
	fastest, slowest = times[0], times[0]

	for _, duration := range times {
		total += duration
		if duration < fastest {
			fastest = duration
		}
		if duration > slowest {
			slowest = duration
		}
	}

	avg = total / time.Duration(len(times))
	return avg, fastest, slowest
}

// Run benchmarks for sorting by Name
func BenchmarkSortByName(b *testing.B) {
	for _, size := range datasetSizes {
		benchmarkSortingWithCores(b, size, "Alphabetical", alphabeticalComparator)
	}
}

// Run benchmarks for sorting by Date
func BenchmarkSortByDate(b *testing.B) {
	for _, size := range datasetSizes {
		benchmarkSortingWithCores(b, size, "Date", dateComparator)
	}
}

// Run benchmarks for sorting by Size (Ascending)
func BenchmarkSortBySizeAsc(b *testing.B) {
	for _, size := range datasetSizes {
		benchmarkSortingWithCores(b, size, "Size Ascending", sizeAscComparator)
	}
}

// Run benchmarks for sorting by Size (Descending)
func BenchmarkSortBySizeDesc(b *testing.B) {
	for _, size := range datasetSizes {
		benchmarkSortingWithCores(b, size, "Size Descending", sizeDecComparator)
	}
}
