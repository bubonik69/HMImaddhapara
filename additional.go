package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

func ParseLine (line string){
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

func checkErr(err error){
	if err!=nil{
		log.Println(err)
	}
}
// write files name in slice from path from json file
type myJsonConfig struct {
	FileConfig struct {
		FileDerectory string `json:"fileDerectory"`
	} `json:"fileConfig"`
}
func configFilePath() (filesPath []string , err error) {
	myCfg :=myJsonConfig{}
	f,err:=ioutil.ReadFile("config.json")
	checkErr(err)
	if json.Valid(f){
		err=json.Unmarshal(f,&myCfg)
	}
	fmt.Printf("в дерректории:  %s находятся следующие файлы: \n",myCfg.FileConfig.FileDerectory )
	fdir,err:=ioutil.ReadDir(myCfg.FileConfig.FileDerectory)
	checkErr(err)
	// show all files
	for i,file:= range fdir{
		fmt.Printf("%d - %s \n" ,i,file.Name())
		if file.IsDir()==false{
			filesPath= append(filesPath,myCfg.FileConfig.FileDerectory+"\\"+file.Name())
		}
	}

	return filesPath,err
}

// take 2 slices and return new slice with only difference
 func diffInSlices (slice1 []string,slice2 []string)(diffS []string){
 	var biggerSlice,smallerSlice []string
 	s1_map:=make(map[string] bool)
 	if len(slice1)>len(slice2){
 		biggerSlice=slice1
 		smallerSlice=slice2
	}else{
		biggerSlice=slice2
		smallerSlice=slice1
	}

 	for _,elS1:= range biggerSlice{
 		s1_map[elS1] = true
	}

 	for _,elS2:=range smallerSlice{
 		if _,ok:=s1_map[elS2];ok!=true{
			diffS= append(diffS,elS2)
		}
		}
		return
	}
