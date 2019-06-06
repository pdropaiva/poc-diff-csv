package util

import (
	"fmt"
	"runtime"
	"time"

	"github.com/pdropaiva/poc-diff-csv/domain"
)

// SplitDiff will split diff map on two arrays one with added users and other with removed users
func SplitDiff(diff map[string]*domain.ExportDiff) (add []domain.UserAudience, remove []domain.UserAudience) {
	for _, u := range diff {
		if !u.IsOld && u.IsNew {
			add = append(add, u.Data)
		}

		if u.IsOld && !u.IsNew {
			remove = append(remove, u.Data)
		}
	}
	return add, remove
}

// PrintDiff will print count and content of add and remove arrays
func PrintDiff(add, remove []domain.UserAudience) {
	fmt.Println("************* Count add *************")
	fmt.Println(len(add))
	fmt.Println("************* Array add *************")
	fmt.Println(add)
	fmt.Println("************ Count remove ***********")
	fmt.Println(len(remove))
	fmt.Println("************ Array remove ***********")
	fmt.Println(remove)
}

// Benchmark  ...
func Benchmark(ctxName string, callback domain.Fn) {
	start := time.Now()

	callback()

	fmt.Printf("\n********************************\n")
	fmt.Printf("%s Benchmark Stats\n", ctxName)
	fmt.Printf("********************************\n")

	fmt.Printf("Tempo de execução:\n\t%v\n\n", time.Since(start))
	fmt.Printf("Consumo de memória:\n")
	printMemUsage()
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("\tAlloc = %v MiB\n", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB\n", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB\n", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
