package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hiscaler/mysql2es/inoutput"
	"log"
	"math"
	"time"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	fmt.Println("Begin sync")
	beginDatetime := time.Now()
	var worker inoutput.Worker
	row := &inoutput.Row{}
	worker = row
	if err := worker.Init(); err == nil {
		worker.Read()
		var totalCount int64
		insertCount, updateCount, deleteCount, err := worker.Write()
		totalCount = int64(insertCount) + int64(updateCount) + int64(deleteCount)
		seconds := time.Since(beginDatetime).Seconds()
		fmt.Println(fmt.Sprintf(`
====================
     Summary
====================
 Insert: %d
 Update: %d
 Delete: %d
--------------------
         %d
--------------------
Seconds: %s
    Avg: %s/s
====================
`, insertCount, updateCount, deleteCount, totalCount, fmt.Sprintf("%.2f", seconds), fmt.Sprintf("%.0f", math.Round(float64(totalCount)/seconds))))
		if err != nil {
			fmt.Println(" > Write error: ", err)
		}
	}
	fmt.Println("Done sync")

}
