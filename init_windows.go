package main

//Windows-only check to prevent... accidents with incorrect use and improve user experience
//by denying ability to start from explorer by double-click and giving correct warnings
import (
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
)

func init() {
	if mousetrap.StartedByExplorer() {
		fmt.Println("Don't double-click ponydownloader")
		fmt.Println("You need to open cmd.exe and run it from the command line!")
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}
