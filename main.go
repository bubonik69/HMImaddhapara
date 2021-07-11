package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)


// variables for work with daatabase
// dataHMI - data from HMI
type dataHMI struct {
	DTRecord string `db:"dt_record"`
	SkipSide	string `db:"skip_side"`
	SkipNum	int64		`db:"skip_num"`
	SkipWeight float64 `db:"skip_weight"`
}
var toDB []dataHMI
var filesPathCurr,filesPath [] string
func main() {
	// write names of files in slice
	createAndOpenDB("batcher_Scale")
	db1,err:=createAndOpenTable("batcher_Scale","hmi")
	_,err=db1.Query("INSERT INTO `HMI` (`dt_record`,`skip_side`,`skip_num`,`skip_weight`) VALUES('2021-07-11 18:05:24','l',1,1)")
	checkErr(err)
	db2,err:=createAndOpenTable("batcher_Scale","files")
	_,err=db2.Query("INSERT INTO `FILES` (`file`) VALUES('test1.txt')")

	checkErr(err)

	checkErr(err)
	var i int
	var newFilesInDir []string
	ticker := time.NewTicker(time.Minute)
	for{
		t:=<-ticker.C
		// start process once a minute
		if t.Minute()%1 == 0 {
			go func() {
				filesPathCurr,err := configFilePath()
				fmt.Println(filesPathCurr)
				fmt.Println(filesPath)
				checkErr(err)
				if len(filesPathCurr)!=len(filesPath){
					newFilesInDir=diffInSlices(filesPathCurr,filesPath)
					}
				i+=1
				fmt.Println(i,"-", time.Now())
				fmt.Println(newFilesInDir)
			}()
			}
		}


	}


