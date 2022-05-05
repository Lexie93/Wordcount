package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"net/rpc"
	"wordcount_service"
	"strconv"
	"math"
)

const work_max=10

// Check file in directories and add them to the list
func check_dirs(files []string) []string {
	var f []string
	for i:=0; i<len(files); i++ {
		fi, err := os.Lstat(files[i])
		if err != nil {
			log.Fatal(err)
		}
		if fi.Mode().IsRegular() {
			f=append(f, fi.Name())
		} else if fi.Mode().IsDir() {
			f_dir, err:=ioutil.ReadDir(fi.Name())
			if err != nil {
				log.Fatal(err)
			}
			for _, f_app:= range f_dir {
				if f_app.Mode().IsRegular() {
					f=append(f, fi.Name() + "/" + f_app.Name())
				}
			}
		}
	}
	return f
}

// Check file existence, delete wrong inputs
func check_files(files []string) []string {
	del:=0
	for i:=0; i<len(files)-del; i++ {
		if _, err := os.Lstat(files[i]); os.IsNotExist(err){
			fmt.Printf("File %s does not exist, it will be omitted during word count\n", files[i])
			files[i]=files[len(files)-1-del]
			del++
			i--
		}
	}
	return check_dirs(files[:len(files)-del])
}

// Read text file
func read_file(file string) string {
	f_byte, err:=ioutil.ReadFile(file)
	if err!=nil{
			log.Fatalf("Cannot open file: %s, error: %v\n", file, err)
		}
	return string(f_byte)
}

// Send some files to a worker
func assign_map(file []string, client *rpc.Client, ch chan wordcount_service.Couple) {
	var c wordcount_service.Couple
	var text string
	for i:=0; i<len(file); i++ {
		text+=read_file(file[i])
	}
	err:= client.Call("Counter.Map", text, &c)
	if err != nil {
		log.Fatal("Error in Counter.Map: ", err)
	}
	ch<-c
}

// Hash a string to decide the worker to send it to for the reduce fase
func hash(str string, w_num int) int {
	sum:=0
	for _, c:= range str {
		sum+= int(c)
	}
	return sum % w_num
}

// Divide work for workers in reduce fase
func partitioner(w wordcount_service.Couple, w_num int) []wordcount_service.Couple {
	var s [work_max]wordcount_service.Couple
	for i:=0; i<len(w); i++ {
		h:=hash(w[i].Word, w_num)
		s[h]=append(s[h], w[i])
	}
	return s[0:w_num]
}

// Determine how many files send to each worker for the map fase
func equality(l int, work_num int) []int {
	var bound []int
	if l>work_num {
		l_app:= l%work_num
		for i:=0; i<l_app; i++ {
			bound=append(bound, int(math.Ceil(float64(l)/float64(work_num))))
		}
		for i:=l_app; i<work_num; i++ {
			bound=append(bound, int(math.Floor(float64(l)/float64(work_num))))
		}
	} else {
		for i:=0; i<l; i++ {
			bound=append(bound, 1)
		}
	}
	return bound
}

func main() {
	var client [work_max]*rpc.Client
	var call [work_max]*rpc.Call
	var r [work_max]wordcount_service.Couple
	var mid, result wordcount_service.Couple
	var err error

	// Check input argument number
	if len(os.Args[1:])<2 {
		log.Fatalf("Usage: %s <#worker> <file_1> (... <file_n>)\n", os.Args[0])
	}

	// Check if number of workers is expressed
	work_num, err:=strconv.Atoi(os.Args[1])
	if err!=nil || work_num<1 {
		log.Fatalf("<#worker> must be an int >0... Usage: %s <#worker> <file_1> (... <file_n>)\n", os.Args[0])
	}
	if work_num>work_max {
		work_num=work_max
	}

	// Check if files exist and ignore them if not
	// Check if some directories are expressed and add files inside to list in case
	files:=check_files(os.Args[2:])
	if len(files)<1 {
		log.Fatalf("Usage: %s <#worker> <file_1> (... <file_n>)\n", os.Args[0])
	}

	// Determine work division
	eq:=equality(len(files), work_num)

	// Connect to workers following incremental port values, starting from :1234
	for i:=0; i<len(eq); i++ {
		client[i], err = rpc.Dial("tcp", "localhost:" + strconv.Itoa(1234+i))
		if err != nil {
			log.Fatal("Error in dialing: ", err)
		}
		defer client[i].Close()
	}

	// Start map fase
	c:=make(chan wordcount_service.Couple)
	start:=0
	for i:=0; i<len(eq); i++ {
		go assign_map(files[start:start+eq[i]], client[i], c)
		start+=eq[i]
	}

	// Wait for map fase to end and
	// synchronize before reduce fase starts
	for i:=0; i<len(eq); i++ {
		mid=append(mid, <-c...)
	}
	part:=partitioner(mid, len(eq))

	// Start reduce fase
	for i:=0; i<len(eq); i++ {
		call[i]= client[i].Go("Counter.Reduce", part[i], &r[i], nil)
	}

	// Wait for reduce fase to end and print final result
	for i:=0; i<len(eq); i++ {
		call[i]=<-call[i].Done
		result=append(result, r[i]...)
	}
	fmt.Println(result)
}
