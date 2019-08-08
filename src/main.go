package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"inoutput"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	fmt.Println("Begin sync")
	beginTimestamp := time.Now().Unix()
	var worker inoutput.Worker
	row := &inoutput.Row{}
	worker = row
	if err := worker.Init(); err == nil {
		worker.Read()
		insertCount, updateCount, deleteCount, err := worker.Write()
		fmt.Println(fmt.Sprintf(" > Insert: %d, Update: %d, Delete: %d, cost %d seconds.", insertCount, updateCount, deleteCount, time.Now().Unix()-beginTimestamp))
		if err != nil {
			fmt.Println(" > Write error: ", err)
		}
	}
	fmt.Println("Done sync")

}
