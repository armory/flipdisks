package main

import "fmt"

func main() {
	//height := 28
	//width := 7

	//values := map[int]byte{
	//	0: byte(0),
	//	1: byte(1),
	//}

	//fmt.Println(values)

	//rows := []int{1, 1, 1, 1, 1, 1, 1}

	data := make([]byte, 2)
	//for x := range data {
	//	d := 0
	//	for y := 0; y < width; y++ {
	//		d = d<<1 | int(rows[y])
	//	}
	//	//fmt.Println(d)
	//	data[x] = byte(d)
	//}

	set(&data, 0, 0, 0)
	set(&data, 0, 1, 0)
	set(&data, 0, 2, 0)
	set(&data, 0, 3, 0)
	set(&data, 0, 4, 0)
	set(&data, 0, 5, 0)
	set(&data, 0, 6, 0)
	//fmt.Printf("%08b\n", 1)
	fmt.Println(data)
}

func set(dataPointer *[]byte, x, y, val int) {
	data := *dataPointer

	fmt.Printf("%08b", data[x])

	if val == 1 {
		fmt.Printf(" setting 1 = ")
		if y > 0 {
			data[x] = data[x] | 1<<uint(y)
		} else {
			data[x] = data[x] | 1
		}
	} else {
		fmt.Printf(" setting 0 = ")
		if y > 0 {
			data[x] = data[x] ^ 1<<uint(y)
		} else {
			data[x] = data[x] ^ 1
		}
	}

	fmt.Printf("%08b\n", data[x])
}
