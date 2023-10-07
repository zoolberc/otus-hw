package main

import (
	"fmt"
	"strings"
)

func main() {
	//files, err := os.ReadDir(".")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for i := 0; i < len(files); i++ {
	//	fileInfo, _ := files[i].Info()
	//	fmt.Printf("name - %s, isDir - %v\n", fileInfo.Name(), fileInfo.IsDir())
	//
	//}

	//err := os.Unsetenv("11")
	//fmt.Println(err)

	//stringToTrim := "\t\t\n   Go \tis\t Awesome \t\t  "
	//trimResult := strings.TrimSpace(stringToTrim)
	//fmt.Println(trimResult)

	fmt.Print(strings.TrimRight("¡¡¡Hello, Gophers	", "\t "))
}
