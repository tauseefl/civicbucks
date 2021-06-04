package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var GlobalResults []ComputeResult

// MinerExecutor starts the mining process and returns the results
func MinerExecutor(startBlock int, endBlock int, timeout int, threads int, timeoutOverrride bool) ([]ComputeResult, error) {
	fmt.Println("\n*********************")
	fmt.Println("* Mining CivicBucks *")
	fmt.Println("*********************")
	fmt.Printf("\nBlocks to mine: %v to %v\n", startBlock, endBlock)
	fmt.Printf("Timeout(seconds): %v, Override: %v\n", timeout, timeoutOverrride)
	fmt.Printf("Threads: %v\n", threads)

	var allResults []ComputeResult
	var wg sync.WaitGroup
	var mu = &sync.Mutex{}
	wg.Add(threads)

	// Send each block one by one to MinerSingle and time for each returns.
	// Append results to allResults before continuing
	// Wrap entire call to wait for the operation to finish or the timer to run out
	routineStart := startBlock
	for i := 1; i <= threads; i++ {
		routineEnd := routineStart + (int((endBlock - startBlock) / threads)) + 1

		if routineEnd > endBlock {
			routineEnd = endBlock
		}

		if routineStart > routineEnd {
			routineStart = routineEnd
		}

		// Print out for debugging threads
		// fmt.Printf("\nThread %v starting at block %v and ending at block %v", i, routineStart, routineEnd)

		// // Use local computations
		// go func(routineStart int, routineEnd int, threadNumber int) {
		// 	defer wg.Done()
		// 	for j := routineStart; j < routineEnd; j++ {
		// 		start := time.Now()
		// 		result, err := MinerSingle(j)
		// 		elapsed := time.Since(start)
		// 		if err == nil {
		// 			result.time = int(elapsed.Nanoseconds())
		// 			mu.Lock()
		// 			allResults = append(allResults, result)
		// 			GlobalResults = allResults
		// 			mu.Unlock()
		// 		}
		// 	}
		// }(routineStart, routineEnd, i)

		// Call server, for each number
		go func(routineStart int, routineEnd int, threadNumber int) {
			defer wg.Done()
			for j := routineStart; j < routineEnd; j++ {
				postBody, _ := json.Marshal(map[string]int{
					"block": j,
				})
				responseBody := bytes.NewBuffer(postBody)
				result, err := http.Post("http://localhost:8082/MinerSingle", "application/json", responseBody)

				if err != nil {
					// do something
				}

				mu.Lock()
				test_num := ComputeResult{}
				err = json.NewDecoder(result.Body).Decode(&test_num)
				mu.Unlock()

				if err == nil {
					mu.Lock()
					allResults = append(allResults, test_num)
					GlobalResults = allResults
					mu.Unlock()
				}
			}
		}(routineStart, routineEnd, i)

		// 	// Call server using blocks of numbers
		// 	go func(routineStart int, routineEnd int, threadNumber int) {
		// 		defer wg.Done()
		// 		postBody, _ := json.Marshal(map[string]int{
		// 			"startBlock": routineStart,
		// 			"endBlock": routineEnd,
		// 		 })
		// 		responseBody := bytes.NewBuffer(postBody)
		// 		result, err := http.Post("http://localhost:8082/MinerBlock", "application/json", responseBody)

		// 		if err != nil {
		// 			// do something
		// 		}

		// 		mu.Lock()
		// 		test_num := []ComputeResult{}
		// 		err = json.NewDecoder(result.Body).Decode(&test_num)
		// 		mu.Unlock()

		// 		if err == nil {
		// 			mu.Lock()
		// 			for _, element := range test_num {
		// 				allResults = append(allResults, element)
		// 			}

		// 			GlobalResults = allResults
		// 			mu.Unlock()
		// 		}
		// 	}(routineStart, routineEnd, i)

		routineStart = routineEnd + 1
	}

	// timeoutOverrride := true

	if waitTimeout(&wg, time.Duration(timeout)*time.Second, timeoutOverrride) {
		fmt.Println("\nReturning partial data...")
	} else {
		fmt.Println("\n\nMining Completed!!")
	}

	// Return either partial or full data
	return allResults, nil

}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
// Returns true if there is a syscall.SIGINT
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration, override bool) bool {
	c := make(chan struct{})
	manualInterrupt := make(chan os.Signal, 1)
	signal.Notify(manualInterrupt, os.Interrupt)

	go func() {
		defer close(c)
		wg.Wait()
	}()
	if !override {
		select {
		case <-c:
			return false // completed normally
		case <-time.After(timeout):
			fmt.Printf("\n\nMining took longer than %v seconds. Timing out :(", timeout)
			return true //timed out
		case <-manualInterrupt:
			fmt.Println("\nOS Interrupt received: Shutting Mining operation down")
			return true // manual OS interrupt
		}
	} else {
		select {
		case <-c:
			return false // completed normally
		case <-manualInterrupt:
			fmt.Println("\nOS Interrupt received: Shutting Mining operation down")
			return true // manual OS interrupt
		}
	}

}

func getMiningResults() []ComputeResult {
	return GlobalResults
}
