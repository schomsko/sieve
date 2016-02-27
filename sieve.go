package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now().UnixNano()
	const limit = 1000000000
	const concur = 4
	numbers := make([]bool, limit) // false means prime candidate
	numbers[1] = true              // 1 not considered prime
	num := 2
	c := make(chan bool)
	r := 0

	for {
		if r < concur {
			// first allowed optimization:  loop only goes to sqrt(limit)
			num2 := num * num
			if num2 >= limit {
				break
			}

			if num < limit/1000 {
				r = r + 3
				go eliminator(num2, num, numbers, limit/3, c)
				go eliminator(int((limit/3)/num)*num, num, numbers, 2*limit/3, c)
				go eliminator(2*int((limit/3)/num)*num, num, numbers, limit, c)
			} else {
				r++
				go eliminator(num2, num, numbers, limit, c)
			}
		}

		select {
		case <-c:
			r--
		default:
		}

		// scan to get next prime for eliminators
		if r < concur {
			for {
				num++
				if !numbers[num] {
					// fmt.Println("Prime: ", num)
					break
				}
			}
		}
	}

	fmt.Println((time.Now().UnixNano() - t) / 1000000)
	// sieve complete.  now print a representation.
	ergebnis := 0
	for n := 1; n < limit; n++ {
		if !numbers[n] {
			// fmt.Println(n)
			ergebnis++
		}
	}
	fmt.Println(ergebnis)

}

func eliminator(num2 int, num int, numbers []bool, limit int, c chan<- bool) {
	for i := num2; i < limit; i += num { // second allowed optimization:  eliminator starts at sqr(p)
		numbers[i] = true // it's a composite
	}
	c <- true
}
