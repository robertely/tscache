package main

import (
	"fmt"
	"time"

	"github.com/robertely/tscache/tscache"
)

func main() {
	x := tscache.NewCollection(50000)

	for i := int64(10); i <= 22; i++ {
		x.Write(i, time.Unix(i, 0))
	}
	fmt.Println("XXXXXXXXXX")
	x.PrintAll()
	// fmt.Println(x)
	// x.Write(20.800000000, time.Unix(20, 800000000))
	// fmt.Println(x.search(time.Unix(12, 0), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(14, 0)))
	// fmt.Println(x.search(time.Unix(12, 800000000), time.Unix(100, 0)))
	resStart, resEnd, length := x.Search(time.Unix(12, 0), time.Unix(20, 0))
	fmt.Println(resStart, resEnd, length)
	fmt.Println("XXXXXXXXXX")
	for _, i := range x.Read(x.Search(time.Unix(12, 0), time.Unix(20, 0))) {
		fmt.Println(i)
	}
}
