package main

import (
	"log"
	"zhibo/cmd"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	cmd.Execute()
}
