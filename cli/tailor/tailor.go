package main

import (
	"fmt"
	"os"

	"github.com/parro-it/tailor"
)

func main() {
	f, err := tailor.OpenFile(os.Args[1], 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	buf := make([]byte, 100)
	for {
		n, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Print(string(buf[:n]))
	}
}
