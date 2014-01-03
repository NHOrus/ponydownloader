package main

import ("fmt"
//	"net"
	"os"
//	"errors"
	"log"
//	"io"
	"github.com/vaughan0/go-ini"
	)	


func main(){
	fmt.Println("Check one")
	config, err := ini.LoadFile("config.ini")
	if os.IsNotExist(err) { panic("config.ini does not exist, create it")}
	if err != nil { log.Fatal(err) }

	key, ok := config.Get("main", "key")
	if !ok {
		panic("'key' variable missing from 'main' section")
		}

	fmt.Println(key)
//	config.Close()
}
