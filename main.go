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
	numberWidth := func(number int64) int {
		n := len(fmt.Sprintf("%d", number))
		return n
	}
	fmt.Println("Begin sync")
	beginDatetime := time.Now()
	var worker inoutput.Worker
	row := &inoutput.Row{}
	worker = row
	if err := worker.Init(); err == nil {
		var totalCount, totalInsertCount, totalUpdateCount, totalDeleteCount int64
		for {
			worker.Read()
			insertCount, updateCount, deleteCount, err := worker.Write()
			totalInsertCount += insertCount
			totalUpdateCount += updateCount
			totalDeleteCount += deleteCount
			totalCount = totalInsertCount + totalUpdateCount + totalDeleteCount
			seconds := time.Since(beginDatetime).Seconds()
			strLen := numberWidth(insertCount)
			if n := numberWidth(updateCount); n > strLen {
				strLen = n
			}
			if n := numberWidth(deleteCount); n > strLen {
				strLen = n
			}
			numberFormat := fmt.Sprintf("%% %dd", strLen)
			fmt.Println(numberFormat)
			fmt.Println(fmt.Sprintf(`
============================================
 Current                Summary
--------------------------------------------
 Insert: %s             Insert: %d
 Update: %s             Update: %d
 Delete: %s             Delete: %d
--------------------------------------------
 Total: %d   Seconds: %s   Avg: %s/s
============================================
`, fmt.Sprintf(numberFormat, insertCount), totalInsertCount, fmt.Sprintf(numberFormat, updateCount), totalUpdateCount, fmt.Sprintf(numberFormat, deleteCount), totalDeleteCount, totalCount, fmt.Sprintf("%.2f", seconds), fmt.Sprintf("%.0f", math.Round(float64(totalCount)/seconds))))
			time.Sleep(10 * time.Second)
			if err != nil {
				fmt.Println(" > Write error: ", err)
			}
		}
	}
	fmt.Println("Done sync")

}
