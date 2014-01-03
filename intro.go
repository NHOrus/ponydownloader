package main

import ("fmt"
//	"net"
//	"os"
//	"io"
	"github.com/vaughan0/go-ini"
	)	


func main(){
	fmt.Println("Check one")
	config, err := ini.LoadFile("config.ini")

	if err != nil { panic(err) }

	key, ok := config.Get("main", "key")
	if !ok {
		panic("'key' variable missing from 'main' section")
		}

	fmt.Println(key)
//	config.Close()
}
