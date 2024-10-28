package stats

import (
	"encoding/csv"
	"errors"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Statistics struct {
	Mean   float64 `json:"mean"`
	Median float64 `json:"median"`
	StdDev float64 `json:"stddev"`
}

func CalculateStatisticsN(columnName []string, filePath string) (*Statistics, *time.Time, error) {
	startTime := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	columnIndex := -1
	for i, header := range records[0] {
		if strings.EqualFold(header, columnName[0]) {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return nil, nil, errors.New("column not found")
	}

	var values []float64

	// Parse values from the specified column
	for _, record := range records[1:] { // Skip the header
		if len(record) <= columnIndex {
			continue
		}
		value, err := strconv.ParseFloat(record[columnIndex], 64)
		if err == nil {
			values = append(values, value)
		}
	}

	// Calculate statistics
	stats := calculateStats(values)

	return &stats, &startTime, nil
}
func CalculateStatistics(columnName []string, filePath string) (*Statistics, *time.Time, error) {
	startTime := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	columnIndex := -1
	for i, header := range records[0] {
		if strings.EqualFold(header, columnName[0]) {
			columnIndex = i
			break
		}
	}

	if columnIndex == -1 {
		return nil, nil, err
	}

	var values []float64
	var mu sync.Mutex
	var wg sync.WaitGroup
	const numWorkers = 4
	chunkSize := len(records[1:]) / numWorkers

	// Worker function to parse values and calculate partial sums
	processChunk := func(chunk [][]string) {
		defer wg.Done()
		localValues := make([]float64, 0, len(chunk))
		for _, record := range chunk {
			if len(record) <= columnIndex {
				continue
			}
			value, err := strconv.ParseFloat(record[columnIndex], 64)
			if err == nil {
				localValues = append(localValues, value)
			}
		}
		mu.Lock()
		values = append(values, localValues...) // Merge results
		mu.Unlock()
	}

	// Distribute work among workers
	for i := 1; i < len(records); i += chunkSize {
		end := i + chunkSize
		if end > len(records) {
			end = len(records)
		}
		wg.Add(1)
		go processChunk(records[i:end])
	}

	wg.Wait()

	// Calculate statistics
	stats := calculateStats(values)

	return &stats, &startTime, nil
}
func calculateStats(values []float64) Statistics {
	n := len(values)
	if n == 0 {
		return Statistics{Mean: 0, Median: 0, StdDev: 0}
	}

	// Calculate mean
	var sum float64
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(n)

	// Calculate median
	var median float64
	sortedValues := make([]float64, n)
	copy(sortedValues, values)
	sort.Float64s(sortedValues)
	if n%2 == 0 {
		median = (sortedValues[n/2-1] + sortedValues[n/2]) / 2
	} else {
		median = sortedValues[n/2]
	}

	// Calculate standard deviation
	var variance float64
	for _, v := range values {
		variance += (v - mean) * (v - mean)
	}
	stdDev := math.Sqrt(variance / float64(n))

	return Statistics{Mean: mean, Median: median, StdDev: stdDev}
}
