package main

import (
	"fmt"

	"github.com/danivideda/satu-apotek-be/scripts/playground/enum"
)

func main() {
	// runScript2()
	enum.RunScript3()
	
	var MyState enum.ServerState = 2
	fmt.Println(MyState)
}