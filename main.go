package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	_ "time"
)


// variables for work with daatabase
// dataHMI - data from HMI
type dataHMI struct {
	DTRecord string 	`db:"dt_record"`
	SkipSide	string 	`db:"skip_side"`
	SkipNum	int64		`db:"skip_num"`
	SkipWeight float64 	`db:"skip_weight"`
}
type filesHMI struct {
	id uint64 `db: " id"`
	file_name string `db:"file_name"`
}

var toDB []dataHMI
var filesPathCurr,filesPath [] string
func main() {
	// read settings for connect to Database from json file
	r,err:=readDBSettings()
	// open database / create database if not exist
	createDBifNotExist(r.DBSettings.DBname)
	// open table / create table if not exist
	db,err:=createTables()
	checkErr(err)

	defer db.Close()
	// write names of files to slice which need to put in database
	filesPathCurr,err := getFilesFromDir()
	checkErr(err)
	// print all files
	//db.Close()
	toDB,err = retStructForDatabase(filesPathCurr)
	checkErr(err)
	err=writeDataToHMItable(toDB)
	checkErr(err)
	for i:=0; i<len(filesPathCurr); i++ {
		filesPathCurr[i]=addSlash(filesPathCurr[i])
	}
	for _, value:=range filesPathCurr{
		_,err=db.Exec(fmt.Sprintf("INSERT INTO `FILES` (`file`) VALUES('%s')",value))
	}


	var i int
	ticker := time.NewTicker(time.Minute)
	for{
		t:=<-ticker.C
		// start process once a minute
		if t.Minute()%1 == 0 {
			go func() {
				newFilesInDir,err := getFilesFromDir()
				filesInDB,err:=readFilesFromDB()
				checkErr(err)
				if len(newFilesInDir)!=len(filesInDB){
					fmt.Print("you have a new files in directory\n")
					s:=diffInSlices(newFilesInDir,filesInDB)
					fmt.Println("added file",s)
					// add data in database
					toDB,err = retStructForDatabase(s)
					checkErr(err)
					err=writeDataToHMItable(toDB)
					checkErr(err)
					// write new files in database
					for i:=0; i<len(s); i++ {
						s[i]=addSlash(s[i])
					}
					for _, value:=range s{
						_,err=db.Exec(fmt.Sprintf("INSERT INTO `FILES` (`file`) VALUES('%s')",value))
					}
				}

				i+=1
				fmt.Println(i,"-", time.Now())
			}()
			}
		}


}


