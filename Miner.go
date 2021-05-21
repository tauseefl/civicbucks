package main

import (
	"errors"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	// "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"time"
)

// ComputeResult stores the result on palindrome calculation and the time it took
type ComputeResult struct {
	Number int
	Binary string
	Time   int
}

var (
	Palindrome_temp *prometheus.HistogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "Palindrome",
		Help:    "help",
		Buckets: prometheus.DefBuckets,
	}, []string{"base"})
)

// MinerSingle executes the palindrome computations one at a time
func MinerSingle(number int) (ComputeResult, error) {
	// output := make(map[int]string)
	var output ComputeResult

	if isPalindrome(number) {
		if isBinaryPalindrome(number) {
			output.Number = number
			output.Binary = convertToBinary(number)
			return output, nil
		}
	}
	return output, errors.New("Not a palindrome")
}

func isPalindrome(number int) bool {
	// converts int into string and then check if the string is a palindrome

	defer func(begin time.Time) {

		s := time.Since(begin).Seconds()
		ms := s * 1e3
		Palindrome_temp.WithLabelValues("Ten").Observe(ms)

	}(time.Now())

	forwardString := strconv.Itoa(number)
	reversedString := reverse(forwardString)

	return forwardString == reversedString
}

func isBinaryPalindrome(number int) bool {
	// converts int into a string of binary and then check if the string is a palindrome

	defer func(begin time.Time) {

		s := time.Since(begin).Seconds()
		ms := s * 1e3
		Palindrome_temp.WithLabelValues("Two").Observe(ms)

	}(time.Now())

	forwardBinaryString := convertToBinary(number)
	reversedBinaryString := reverse(forwardBinaryString)

	return forwardBinaryString == reversedBinaryString
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func convertToBinary(i int) string {
	i64 := int64(i)
	return strconv.FormatInt(i64, 2) // base 2 for binary
}
