package main

import "log"

func init() {
	log.SetPrefix("LOG >>> ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	log.Println("message")
	log.Fatalln("fatal message")
	log.Panicln("panic message")
}
