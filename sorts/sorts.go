// https://qvault.io/golang/quick-sort-golang/

package qsort

import (
    // "math/rand"
    // "sync"
    // "fmt"
    // "time"
	// "runtime"
	// "unsafe"
)

func partition(arr *[]int, low, high int) int {
	pivot := (*arr)[high]
	i := low
	for j := low; j < high; j++ {
		if (*arr)[j] < pivot {
			(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
			i++
		}
	}
	(*arr)[i], (*arr)[high] = (*arr)[high], (*arr)[i]
	return i
}


func QuickSort(arr *[]int) {
    quickSort(arr, 0, len(*arr) - 1)
}


func quickSort(arr *[]int, low, high int) {
	if low < high {
        p := partition(arr, low, high)
        quickSort(arr, low, p - 1)
		quickSort(arr, p + 1, high)
	}
}



type TASwap struct {
	a1 int
	a2 int
}


func change_parallel(arr *[]int, paar chan TASwap) {
    for val := range paar {
        i := val.a1
        j := val.a2
        (*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
    }
}

func partition_parallel(arr *[]int, low int, high int, paar chan TASwap) int {
	pivot := (*arr)[high]
	i := low
	for j := low; j < high; j++ {
		if (*arr)[j] < pivot {
            var pa TASwap
            pa.a1, pa.a2 = i, j
            paar <- pa
			i++
		}
	}

    var pa TASwap
    pa.a1, pa.a2 = i, high
    paar <- pa

	return i
}


func QuickSort_parallel(arr *[]int) {

	c := make(chan struct{})
    go quickSort_parallel(arr, 0, len(*arr) - 1, c)
    <- c
}


func quickSort_parallel(arr *[]int, low int, high int, ch chan struct{}) {
	if low < high {

        paar := make(chan TASwap, (high - low) + 1)
        p := partition_parallel(arr, low, high, paar)

        close(paar)

        change_parallel(arr, paar)

	    c := make(chan struct{})

        go quickSort_parallel(arr, low, p - 1, c)
		go quickSort_parallel(arr, p + 1, high, c)

        <-c
        <-c
	}

	ch <- struct{}{}
}



/*

// https://github.com/farazdagi/bitonic/blob/master/sorter.go

const (
	SORT_ASC  bool = true
	SORT_DESC bool = false
)


func BitonicSort(arr *[]int, dir bool) {
	sentinel := make(chan struct{})
	go bitonicSort(arr, 0, len(*arr), dir, sentinel)
	<-sentinel
}


func bitonicSort(arr *[]int, lo int, n int, dir bool, sentinel chan struct{}) {
	if n > 1 {
		m := n / 2

		c := make(chan struct{})

		go bitonicSort(arr, lo, m, SORT_ASC, c)
		go bitonicSort(arr, lo + m, m, SORT_DESC, c)

		<-c
		<-c

		bitonicMerge(arr, lo, n, dir, sentinel)
	} else {
		sentinel <- struct{}{}
	}
}


func bitonicMerge(arr *[]int, lo int, n int, dir bool, sentinel chan struct{}) {
	if n > 1 {
		m := n / 2

		for i := lo; i < lo+m; i++ {
			compareAndSwap(arr, i, i + m, dir)
		}

		c := make(chan struct{})

        go bitonicMerge(arr, lo, m, dir, c)
		go bitonicMerge(arr, lo + m, m, dir, c)

		<-c
		<-c
	}

	sentinel <- struct{}{}
}


func compareAndSwap(arr *[]int, i int, j int, dir bool) {
	if dir == ((*arr)[i] > (*arr)[j]) {
		(*arr)[i], (*arr)[j] = (*arr)[j], (*arr)[i]
    }
}

*/
