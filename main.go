package main

import (
	"fmt"
	"time"
)

type NumberStatus string

const (
	Even NumberStatus = "even"
	Odd  NumberStatus = "odd"
)

type Data struct {
	Number int64
	Status NumberStatus
}

func square(number Data) Data {
	number.Number = number.Number * number.Number
	return number
}

func fanOut(in chan Data) chan Data {
	out := make(chan Data)
	go func() {
		for i := range in {
			out <- square(i)
		}
		close(out)
	}()
	return out
}

func fanIn(channels ...chan Data) chan Data {
	out := make(chan Data)
	done := make(chan struct{})

	for _, c := range channels {
		go func(in chan Data) {
			for d := range in {
				out <- d
			}
			done <- struct{}{}
		}(c)
	}

	go func() {
		for i := 0; i < len(channels); i++ {
			<-done // wait for all channels to be closed
		}
		close(out) // close the output channel
	}()

	return out
}

func generateEvenOddNumber(numbers ...int64) chan Data {
	out := make(chan Data)
	go func() {
		for _, num := range numbers {
			if num%2 == 0 {
				out <- Data{
					Number: num,
					Status: Even,
				}
			} else if num%2 != 0 {
				out <- Data{
					Number: num,
					Status: Odd,
				}
			}
		}
		close(out)
	}()
	return out
}

func main() {
	numbers := generateEvenOddNumber(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	w1 := fanOut(numbers) // worker 1
	w2 := fanOut(numbers) // worker 2
	w3 := fanOut(numbers) // worker 3

	result := fanIn(w1, w2, w3)
	for i := range result {
		fmt.Println(i)
	}

	time.Sleep(time.Second)
}
