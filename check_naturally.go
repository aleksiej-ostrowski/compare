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
    "io"
    "os"
    "encoding/binary"
    "bufio"
    "github.com/shawnsmithdev/zermelo"
    "math"
    "math/rand"
    "time"
)

const EPS = 1E-14
const LEN = 100_000_000
const BIG_FILE = "very_big_file.bin"

func IsEqual(a, b float64) bool {
    return math.Abs(a - b) < EPS
}   

func vPearson(v1, v2 []int) float64 {

    mu_1 := .0
    su_1 := .0        
    su_2 := .0
    su_k_1 := .0
    su_k_2 := .0

    sz := len(v1)

    for i:= 0; i < sz; i++ {
        zn1 := float64(v1[i])
        zn2 := float64(v2[i])

        mu_1 += zn1 * zn2

        su_1 += zn1
        su_2 += zn2

        su_k_1 += zn1 * zn1
        su_k_2 += zn2 * zn2
    }

    numerator := float64(sz) * mu_1 - su_1 * su_2
    denominator := math.Sqrt( (float64(sz) * su_k_1 - su_1 * su_1) * (float64(sz) * su_k_2 - su_2 * su_2) )

    if IsEqual(math.Abs(denominator), .0) {
       return .0
    } else {
       return numerator / denominator
    }   

}

func main() {

    if _, e := os.Stat(BIG_FILE); os.IsNotExist(e) {
        errorString := fmt.Sprintf("You must create a very large natural data file %s in the current directory.", BIG_FILE)
		panic(errorString)
    }
       
    f, err := os.Open(BIG_FILE)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer f.Close()

    data := make([]byte, 2)

    f.Seek(100_000, os.SEEK_SET)

    reader := bufio.NewReader(f)    

    ind := 0

    target := make([]int, LEN)

    for {

        n, err := reader.Read(data)    
        if err != nil {    
            if err == io.EOF {    
                break    
            }    
            fmt.Println(err)
        }    
    
        if n == 0 {    
            fmt.Printf("Nothing for read from file %s \n", BIG_FILE)
            break    
        }           
 
        target[ind] = int(binary.BigEndian.Uint16(data))

        ind += 1

        if ind >= LEN {
            break
        }
    }

    b := make([]int, len(target))
    copy(b, target)
    zermelo.Sort(b)

    // fmt.Println(target[:10])
    // fmt.Println(b[len(b)-10:])

    case_1 := math.Abs(vPearson(b, target))
    fmt.Printf("Coefficient for %s : %.10f \n", BIG_FILE, case_1)

    rand.Seed(time.Now().UTC().UnixNano())
    for i, _ := range target {
        target[i] = rand.Intn(math.MaxInt32)
    }
    
    copy(b, target)
    zermelo.Sort(b)

    // fmt.Println(target[:10])
    // fmt.Println(b[:10])

    case_2 := math.Abs(vPearson(b, target))
    fmt.Printf("Coefficient for random data: %.10f \n", case_2)
    fmt.Printf("Coefficient of naturally data: %.10f \n", case_1 / (case_2 + EPS))
}    
