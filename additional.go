package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)


// write files name in slice from path from json file
type myJsonConfig struct {
	FileConfig struct {
		FileDerectory string `json:"fileDerectory"`
	} `json:"fileConfig"`
}
// function reads files in a directory and return a slice with named \ path files
func getFilesFromDir() (filesPath []string , err error) {
	myCfg :=myJsonConfig{}
	f,err:=ioutil.ReadFile("config.json")
	checkErr(err)
	if json.Valid(f){
		err=json.Unmarshal(f,&myCfg)
	}
	fdir,err:=ioutil.ReadDir(myCfg.FileConfig.FileDerectory)
	checkErr(err)
	// show all files
	for _,file:= range fdir{
		if file.IsDir()==false{
			filesPath= append(filesPath,myCfg.FileConfig.FileDerectory+"\\"+file.Name())
		}
	}

	return filesPath,err
}

// take 2 slices and return new slice with only difference
 func diffInSlices (slice1 []string,slice2 []string)(diffS []string){
 	var biggerSlice,smallerSlice []string

 	if len(slice1)>len(slice2){
 		biggerSlice=slice1
 		smallerSlice=slice2
	}else if len(slice1)<len(slice2){
		biggerSlice=slice2
		smallerSlice=slice1
	}else{
		return
	}
	 Map :=make(map[string] bool)
 	for _,elS1:= range smallerSlice{
 		Map[elS1] = true
	}

 	for _,elS2:=range biggerSlice {
 		if _,ok:= Map[elS2];ok!=true{
			diffS= append(diffS,elS2)
		}
		}
		return
	}


// 	preparation data to write in database
func retStructForDatabase (filesPath []string) (toDB []dataHMI, err error) {
	// show all files
	for _, fileName := range filesPath {
		//fmt.Println(fileName)
		file, err := os.Open(fileName)
		checkErr(err)
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
					break
				}
			}
			//delete sign "\n" in string line
			lines = append(lines, line[0:len(line)-2])
		}
		for _, line := range lines {
			//add lines to structure toDB []dataHMI
			line=strings.ToLower(line)
			slice:=strings.Split(line," ")
			if len(slice)>4{
				sn, err :=strconv.ParseInt(slice[10],10,64)
				checkErr(err)
				sw,err :=strconv.ParseFloat(slice[16],64)
				checkErr(err)
				toDB= append(toDB,dataHMI{DTRecord : slice[0] + "-" + slice[1]+"-" + slice[2] + " " + // date
					slice[3]+":" + slice[4]+":" + slice[5], // time
					SkipSide : slice[6],
					SkipNum :sn ,
					SkipWeight: sw,
				})
			}
		}
		file.Close()
	}
	return
}
// add slashes for writing database
func addSlash(in string)(out string){
	for _,sign:= range in{
		out+=string(sign)
		if sign=='\\'{
			out+="\\"
		}
	}
	return
}


// handlers of error
func checkErr(err error){
	if err!=nil{
		log.Println(err)
	}
}
func checkPanicErr(err error){
	if err != nil {
		panic(err)
	}
}