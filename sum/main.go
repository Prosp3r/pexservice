package main

import "fmt"

var num1 uint64
var num2 uint64

func main() {
	num1 = 10957648158266230146
	num2 = 2938395859457206495
	sumx := num1 + num2
	fmt.Printf(" %v \n %v \n = \n %v \n", num1, num2, sumx)
	revSum := sumx - num2
	fmt.Printf(" reverse = \n %v \n", revSum)
}
