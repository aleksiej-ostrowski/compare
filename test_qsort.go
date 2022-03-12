/*

#--------------------------------#
#                                #
#  version 0.0.1                 #
#                                #
#  Aleksiej Ostrowski, 2022      #
#                                #
#  https://aleksiej.com          #
#                                #
#--------------------------------#

*/

package main

import (
    "fmt"
    "math"
    "math/rand"
    "time"
    "sync"
    "runtime"
    "sort"
    // "strings"
    "encoding/json"
    q "./sorts"
    "./priority_queue"
    "github.com/shawnsmithdev/zermelo"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/mem"
    // "github.com/pkg/profile"
    "io/ioutil"
)

type MyOut struct {
    Data [][]time.Duration
    Labels []string
    X []int
    Xlabel string
    Xfilter string
    Ylabel string 
    Title string 
}

type Node struct {
    priority int
    value    int
}

func (this *Node) Less(other interface{}) bool {
    return this.priority < other.(*Node).priority
}

/*

func Merge(arrs...[]int) []int {

    q := priority_queue.New()

    for _, a := range arrs {
        for _, v := range a {
            q.Push(&Node{priority: v, value: v})
        }
    }

    res := make([]int, 0)

    for q.Len() > 0 {
        x := q.Top().(*Node)
        res = append(res, x.value)
        q.Pop()
    }

    return res

}

*/


func simplest_check_arr(arr []int) int64 {
    var hash int64 = 0
    for i, v := range arr {
        hash += int64(i) * int64(v)
    }
    return hash
}



func main() {

    // runtime.GOMAXPROCS(2)

    // defer profile.Start().Stop()

    RESULT := "./result.xml"

    N_CHUNKS := runtime.NumCPU()

    X := []int{1_000, 5_000, 10_000, 30_000, 50_000, 100_000, 1_000_000} //, 10_000_000}

    MET := []string{"sort.Ints()", "RadixSort()", "QuickSort()", "QuickSort_parallel()", "Parallel schema #1", "Parallel schema #2"}

    DATA := make([][]time.Duration, len(MET))

    for i, _ := range DATA {
        DATA[i] = make([]time.Duration, len(X))
    }

    rand.Seed(time.Now().UTC().UnixNano())

    ITER := 10

    for Xi, Xv := range X {

        TIMES := make([]time.Duration, len(MET))

        for i, _ := range TIMES {
            TIMES[i] = time.Duration(0)
        }

        // fmt.Println("for ", Xi + 1, " of ", len(X), "run...")

        for j := 0; j < ITER; j++ {

            a := make([]int, Xv)

            for i, _ := range a {
                a[i] = 1 + rand.Intn(math.MaxInt32)
            }

            var ideal_hash int64

            // 0

            {

            b := make([]int, len(a))
            copy(b, a)

            start := time.Now()

            sort.Ints(b)

            duration := time.Since(start)

            TIMES[0] += duration

            ideal_hash = simplest_check_arr(b)

            }

            // 1

            {

            b := make([]int, len(a))
            copy(b, a)

            start := time.Now()

            zermelo.Sort(b)

            duration := time.Since(start)

            TIMES[1] += duration

            new_hash := simplest_check_arr(b)

            if (ideal_hash != new_hash) {
                panic(fmt.Sprintf("%s does not work", MET[1]))
            }

            }

            // 2

            {

            b := make([]int, len(a))
            copy(b, a)

            start := time.Now()

            q.QuickSort(&b)

            duration := time.Since(start)

            TIMES[2] += duration

            new_hash := simplest_check_arr(b)

            if (ideal_hash != new_hash) {
                panic(fmt.Sprintf("%s does not work", MET[2]))
            }

            }

            // 3

            {

            b := make([]int, len(a))
            copy(b, a)

            start := time.Now()

            q.QuickSort_parallel(&b)

            duration := time.Since(start)

            TIMES[3] += duration

            new_hash := simplest_check_arr(b)

            if (ideal_hash != new_hash) {
                panic(fmt.Sprintf("%s does not work", MET[3]))
            }

            }


            // preparing chunks

            divided := make(chan *[]int, N_CHUNKS)

            {

            b := make([]int, len(a))
            copy(b, a)

            chunkSize := (len(b) + N_CHUNKS - 1) / N_CHUNKS

            for i := 0; i < len(b); i += chunkSize {
                end := i + chunkSize

                if end > len(b) {
                    end = len(b)
                }

                chunk := b[i : end]
                divided <- &chunk
            }

            }


            // 4 

            {

            start := time.Now()

            pq := priority_queue.New()

            ch := make(chan *[]int)

            var wg sync.WaitGroup

            for i := 0; i < N_CHUNKS; i++ {
                wg.Add(1)
                go func(x *[]int) {
                    defer wg.Done()
                    q.QuickSort(x)
                    ch <- x
                }(<- divided)
            }

            go func() {
                wg.Wait()
                close(ch)
            }()

            for vv := range ch {
                for _, v := range *vv {
                    pq.Push(&Node{priority: v, value: v})
                }
            }

            res := make([]int, 0)

            for pq.Len() > 0 {
                x := pq.Top().(*Node)
                res = append(res, x.value)
                pq.Pop()
            }

            duration := time.Since(start)

            TIMES[4] += duration

            new_hash := simplest_check_arr(res)

            if (ideal_hash != new_hash) {
                panic(fmt.Sprintf("%s does not work", MET[4]))
            }

            }


            // preparing chunks

            divided = make(chan *[]int, N_CHUNKS)

            {

            b := make([]int, len(a))
            copy(b, a)

            chunkSize := (len(b) + N_CHUNKS - 1) / N_CHUNKS

            for i := 0; i < len(b); i += chunkSize {
                end := i + chunkSize

                if end > len(b) {
                    end = len(b)
                }

                chunk := b[i : end]
                divided <- &chunk
            }

            }


            // 5

            {

            start := time.Now()

            pq := priority_queue.New()

            ch := make(chan *[]int)

            /*

            for _, v := range b {
                go func(x []int, r chan<- *[]int) {
                    q.QuickSort(&x)
                    r <- &x
                }(v, ch)
            }

            */

            /*

            for i := 0; i < len_divided; i++ {
                go func(x *[]int) {
                    q.QuickSort(x)
                    ch <- x
                }(&b[i])
            }

            */

            for i := 0; i < N_CHUNKS; i++ {
                go func(x *[]int) {
                    q.QuickSort(x)
                    ch <- x
                }(<- divided)
            }

            for i := 0; i < N_CHUNKS; i++ {
                for _, v := range *<-ch {
                    pq.Push(&Node{priority: v, value: v})
                }
            }

            res := make([]int, 0)

            for pq.Len() > 0 {
                x := pq.Top().(*Node)
                res = append(res, x.value)
                pq.Pop()
            }

            duration := time.Since(start)

            TIMES[5] += duration

            new_hash := simplest_check_arr(res)

            if (ideal_hash != new_hash) {
                panic(fmt.Sprintf("%s does not work", MET[5]))
            }

            }
        }

        for i, _ := range TIMES {
            TIMES[i] /= time.Duration(ITER)
        }

        for i, v := range TIMES {
            DATA[i][Xi] = v
        }
    }

    cpuStat, _ := cpu.Info()
    vmStat, _ := mem.VirtualMemory()

    cpu := cpuStat[0].ModelName
    ram := vmStat.Total / 1024 / 1024 / 1024

    info := fmt.Sprintf("%s, logical CPUs: %d, %d GB RAM", cpu, runtime.NumCPU(), ram)

    mydata := &MyOut{
        Data: DATA,
        Labels: MET,
        X: X,
        Xlabel: "N",
        Xfilter: "(x > 10_000) and not((x < 60_000) and (y < 5.))",
        Ylabel: "Time",
        Title:  "Algorithm's compare, " + info,

    }

    // fmt.Println(mydata)

    b, err := json.Marshal(mydata)

    if err != nil {
        fmt.Println(err)
    }
    
    _ = ioutil.WriteFile(RESULT, b, 0644)

    /*
{        
"Data":  [        
        [200, 1000, 10, 2100, 2501, 3501],        
        [505, 5010, 50, 5200, 5502, 5502],        
        [210, 1020, 30, 2300, 2503, 3503],        
        [220, 1030, 40, 2400, 2504, 3504],        
        [225, 1040, 50, 2500, 2505, 3505]        
        ],        
"Labels": ["1", "2", "3", "4", "5"],        
"X": [10000, 100000, 500000, 1000000, 5000000, 10000000],        
"Xlabel": "N",        
"Xfilter": "x > 10_000",
"Ylabel": "Time, sec.",        
"Title": "Algorithm's compare"        
} 

    */

    n := runtime.NumGoroutine()

    if n != 1 {
        fmt.Println("Warning! There are working routines, n = ", n)
    }

}
