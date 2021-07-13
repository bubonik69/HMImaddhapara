package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
)
// struct for Data login
type DBSettingsJson struct{
 	DBSettings struct{
	User string 		`json:"user"`
	Pass string 		`json:"pass"`
	DriverName string 	`json:"driverName"`
	DBname string		`json:"DBname"`
} 						`json:"DatabaseSettings"`
}
// struct for Database


// read settings database
//set in variable
func readDBSettings()(DBsttngs DBSettingsJson, err error){
	f,err:=ioutil.ReadFile("loginData.json")
	checkErr(err)
	if json.Valid(f){
		err=json.Unmarshal(f,&DBsttngs)
		checkErr(err)
	}
	return
}

// open or create database if not exist
func createDBifNotExist(name string) (err error) {
	// read login Data sql server
	r,err:=readDBSettings()
	db, err := sql.Open(r.DBSettings.DriverName, fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/",r.DBSettings.User,r.DBSettings.Pass))
		checkPanicErr(err)
	_,err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS " + name))
		checkPanicErr(err)
	db, err = sql.Open(r.DBSettings.DriverName,
						fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",r.DBSettings.User,r.DBSettings.Pass,name))
		checkPanicErr(err)
	db.Close()
	return err
}


func createTables() (db *sql.DB, err error) {
	// read login Data sql server
	r,err:=readDBSettings()
	db, err = sql.Open(r.DBSettings.DriverName,
				fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",r.DBSettings.User,r.DBSettings.Pass,r.DBSettings.DBname))
	checkPanicErr(err)
	//switch nameTable{
	//case "hmi":
		_,err = db.Exec(
			"CREATE TABLE IF NOT EXISTS HMI (dt_record  timestamp PRIMARY KEY," +
				"skip_side VARCHAR(2) NOT NULL, " +
				"skip_num int UNSIGNED NOT NULL," +
				"skip_weight double)")
		checkPanicErr(err)
	//case "files":
		_,err = db.Exec(
			"CREATE TABLE IF NOT EXISTS FILES (id int(11) PRIMARY KEY AUTO_INCREMENT," +
				"file text NOT NULL)")
		checkPanicErr(err)
	//default:
	//	panic("Not right table name")

	//}
	//defer db.Close()
	return db, err
}

func writeDataToHMItable (toDB []dataHMI) (err error){
		// read login Data sql server
		r, err := readDBSettings()
	//fmt.Println(toDB)
	db, err := sql.Open(r.DBSettings.DriverName,
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",r.DBSettings.User,r.DBSettings.Pass,r.DBSettings.DBname))
	checkErr(err)
	defer db.Close()
	for _, line := range toDB {
		_, err := db.Exec("INSERT INTO `HMI` (`dt_record`,`skip_side`,`skip_num`,`skip_weight`) VALUES( ?,? ,?,?)",
		line.DTRecord, line.SkipSide, line.SkipNum, float64(line.SkipWeight))
		checkErr(err)
	}
	return
}
func readFilesFromDB() (retFiles []string, err error){
	r, err := readDBSettings()
	db, err := sql.Open(r.DBSettings.DriverName,
		fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s",r.DBSettings.User,r.DBSettings.Pass,r.DBSettings.DBname))
	checkErr(err)
	defer db.Close()
	rows,err:=db.Query("SELECT `file` FROM `files`" )
	checkErr(err)
	var file string
	for rows.Next(){
		err=rows.Scan(&file)
		checkErr(err)
		retFiles=append(retFiles,file)
	}
return
}


