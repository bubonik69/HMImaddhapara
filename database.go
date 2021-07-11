package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
)
// open or create database if not exist
func createAndOpenDB(name string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/")
		checkPanicErr(err)
	defer db.Close()

	_,err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS "+name))
		checkPanicErr(err)
	db.Close()

	db, err = sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/" + name)

	defer db.Close()
	return db, err
}

func createAndOpenTable(nameDB string, nameTable string) (db *sql.DB, err error) {
	db, err = sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/" + nameDB)
	checkPanicErr(err)
	//defer db.Close()
	switch nameTable{
	case "hmi":
		_,err = db.Exec(
			"CREATE TABLE IF NOT EXISTS HMI (dt_record  timestamp PRIMARY KEY," +
				"skip_side VARCHAR(2) NOT NULL, " +
				"skip_num int UNSIGNED NOT NULL," +
				"skip_weight double)")
		checkPanicErr(err)
	case "files":
		_,err = db.Exec(
			"CREATE TABLE IF NOT EXISTS FILES (id int(11) PRIMARY KEY AUTO_INCREMENT," +
				"file text NOT NULL)")
		checkPanicErr(err)
	default:
		panic("Not right table name")

	}
	//defer db.Close()
	return db, err
}



func checkPanicErr(err error){
	if err != nil {
	panic(err)
}
}

func databaseHandler (){
	// show all files
	for _,fileName:= range filesPath {
		fmt.Println(fileName)
		file, err := os.Open(fileName)
		checkErr(err)
		defer file.Close()
		reader := bufio.NewReader(file)
		var lines []string
		lines = nil
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println(err)
					return
				}
			}
			// убираем со строки \n
			lines = append(lines, line[0:len(line)-2])
		}
		for _, line := range lines {
			//add lines to structure toDB []dataHMI
			ParseLine(line)
		}
	}
	fmt.Println(toDB)
	db, err := sql.Open("mysql", "mysql:mysql@tcp(127.0.0.1:3306)/1")
	checkErr(err)
	defer db.Close()
	for i, line := range toDB {
		insert, err := db.Query(
			fmt.Sprintf("INSERT INTO `HMI` (`dt_record`,`skip_side`,`skip_num`,`skip_weight`) VALUES('%s','%s',%d,%f)",
				line.DTRecord, line.SkipSide, line.SkipNum, float64(line.SkipWeight)))
		checkErr(err)
		insert.Next()
		fmt.Println(i)
		// close connection if last record
		if i==len(toDB)-1{
			insert.Close()
		}
	}
}