package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"
	"time"
)

/*
This is the entry point to the CivicBucks Mining challenge.
Params:
start	- start of mining block
end 	- end of mining block
timeout	- predetermined timeout
threads	- number of thread allowed
*/

func main() {
	startPtr := flag.Int("start", 1, "starting block")
	endPtr := flag.Int("end", 100, "ending block")
	timeoutPtr := flag.Int("timeout", 30, "timeout to stop the mining")
	threadsPtr := flag.Int("threads", 8, "number of threads to use in the mining process")
	overridePtr := flag.Bool("override", false, "override(bool) timeout to run indefinitely")
	flag.Parse()

	start := *startPtr
	end := *endPtr
	timeout := *timeoutPtr
	threads := *threadsPtr
	override := *overridePtr
	defaultThreads := 8

	if start < 0 || end < 0 || start > end {
		fmt.Println("'start' and 'end' should be positive integers and 'start' should be less than 'end'")
		os.Exit(0)
	}
	if timeout < 1 {
		fmt.Println("'timeout' cannot be less than 1")
		os.Exit(0)
	}
	if threads < 1 {
		fmt.Println("'threads' cannot be less than 1")
		os.Exit(0)
	}

	// Launch Civic Server
	civicServer()

	// Change thread limit if it is different from default
	if threads != runtime.GOMAXPROCS(-1) {
		runtime.GOMAXPROCS(threads)

		// Set back to default value when main() is done
		defer runtime.GOMAXPROCS(defaultThreads)
	}

	startTime := time.Now()
	result, _ := MinerExecutor(start, end, timeout, threads, override)
	elapsedTime := time.Since(startTime)

	fmt.Println("\nPalindromes:")
	printSortedList(result)
	fmt.Printf("\nPalindromes computed: %v\n", len(result))
	printResultPerformance(result)
	fmt.Printf("\nTotal duration: %v\n", elapsedTime)

}

func printSortedList(result []ComputeResult) {
	sort.Slice(result[:], func(i, j int) bool {
		return result[i].number < result[j].number
	})

	for i := 0; i < len(result); i++ {
		fmt.Printf("%v \t binary: %v \t\n", result[i].number, result[i].binary)
	}
}

func printResultPerformance(result []ComputeResult) {
	total := 0
	max := 0
	if len(result) == 0 {
		fmt.Printf("\nNo palindromes were calculated\n")
		return
	}

	for i := range result {
		total = total + result[i].time
		if result[i].time > max {
			max = result[i].time
		}
	}

	average := total / len(result)

	fmt.Printf("\nPerformance (nanoseconds): max: %v, mean: %v\n", max, average)
}
