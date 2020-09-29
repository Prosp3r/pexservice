/*
1. READ CSV INTO MEMORY **
2. USE STRUCT TO STORE STATE IN MEMORY FOR READS AND WRITES **
3. WRITE TO FILE AFTER EVERY UPDATE USING INDEPENTENT ROUTINE QUEUE
4. TRY TERMINATE GRACEFULLY ON CRASH AFTER WRITING TO FILE
*/

package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var port = "8081"

//Fibonacci -
type Fibonacci struct {
	Hitcount uint64 `json:"hits"`
	Position uint64 `json:"postion"`
	Previous uint64 `json:"previous"`
	Current  uint64 `json:"current"`
	Next     uint64 `json:"next"`
}

//FiboStore -- Hold fibonacci data
var FiboStore = make(map[string]uint64)
var mutex = &sync.Mutex{}
var gfibonacci Fibonacci //reachable and writable by all

//failOnError is  single place to handle errors to reduces number of keystrokes per each error handling call.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s %s", msg, err)
	}
}

//fibonacci returns the next number in the sequence
//0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 54
func fibonacciBack() (uint64, error) {
	//backshift
	var newPrevious uint64
	var newNext uint64
	var newCurrent uint64
	var newPosition uint64
	mutex.Lock()
	//if gfibonacci.Previous >= 1 {
	if FiboStore["Previous"] >= 1 {
		newPrevious = FiboStore["Current"] - FiboStore["Previous"]
		newCurrent = FiboStore["Previous"]
		newNext = FiboStore["Current"]
		newPosition = FiboStore["Position"] - 1
	} else {
		newPrevious = 0 //gfibonacci.Previous
		newPosition = 1
		newCurrent = newPrevious
		newNext = newCurrent + 1
	}

	FiboStore["Hitcount"]++
	FiboStore["Position"] = newPosition
	FiboStore["Current"] = newCurrent
	FiboStore["Previous"] = newPrevious
	FiboStore["Next"] = newNext

	mutex.Unlock()
	return newPrevious, nil
}

//fibonacciGo sets the next set of numbers in the sequence //0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 54
func fibonacciGo() (uint64, error) {
	//backshift
	var newPrevious uint64
	var newNext uint64
	var newCurrent uint64
	var newPosition uint64
	mutex.Lock()
	//if gfibonacci.Current <= 0 {
	if FiboStore["Current"] <= 0 {
		newPrevious = 0 //gfibonacci.Previous
		newPosition = 1
		newCurrent = newPrevious + 1
		newNext = newCurrent + newPrevious
	} else {

		newPrevious = FiboStore["Current"]
		newCurrent = FiboStore["Next"]
		newNext = newCurrent + newPrevious
		newPosition = FiboStore["Position"] + 1
	}

	//for {

	//go func() {

	FiboStore["Hitcount"]++
	FiboStore["Position"] = newPosition
	FiboStore["Current"] = newCurrent
	FiboStore["Previous"] = newPrevious
	FiboStore["Next"] = newNext
	//go func(){
	//saveFibo(FiboStore)
	//}()
	mutex.Unlock()

	//}()
	return newCurrent, nil
	//}
}

//next - reads and relays the next number of the fibonacci
func next(w http.ResponseWriter, r *http.Request) {

	next, err := fibonacciGo()
	failOnError(err, "Failed to increment the fibonacci")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(next)
}

func previous(w http.ResponseWriter, r *http.Request) {
	previous, err := fibonacciBack()
	failOnError(err, "Failed to decrement the fibonacci")

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(previous)
}

//current reads and relays the current state of the fibonacci
func current(w http.ResponseWriter, r *http.Request) {

	mutex.Lock()
	fib, _ := json.Marshal(FiboStore["Current"])
	mutex.Unlock()
	xfib, _ := strconv.ParseUint(string(fib), 0, 64)
	//fmt.Println(FiboStore)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(fib)
	json.NewEncoder(w).Encode(xfib)
}

//reset restoes all information to zero
func reset(w http.ResponseWriter, r *http.Request) {

	//go func() {

	mutex.Lock()

	FiboStore["Hitcount"]++
	FiboStore["Position"] = 1
	FiboStore["Previous"] = 0
	FiboStore["Current"] = 0
	FiboStore["Next"] = 1

	//mutex.Unlock()

	//}()
	//saveFibo()

	//mutex.Lock()
	fib, _ := json.Marshal(FiboStore)
	mutex.Unlock()
	//xfib, _ := strconv.ParseUint(string(fib), 0, 64)
	//fmt.Printf("%v", FiboStore["Position"])
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(string(fib))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, " You reached my fibonacci /homepage \n\n The hits parameter is a traffic hit counter for some easy hit load tracking \n\n The position parameter is your current position on the fibonacci scale \n\n Try the following end points \n /previous to view previous number \n /current to view the current number \n /next to view the next number \n /reset to reset the fibonacci to zero \n ")
}

//basicRateHandler - could be implemented as a rate limiter with parameters set on server startup or pre-set default.
func basicRateHandler() { /*TODO*/ }

//save sequence to file
func saveFibo() {
	//fmt.Println("Asked to save")
	filename := "fibo.csv"
	fmt.Printf("Saving fibonacci state to %v \n", filename)
	for {

		f := FiboStore

		//recordString := strconv.FormatUint(gfibonacci.Hitcount, 10) + "," + strconv.FormatUint(gfibonacci.Position, 10) + "," + strconv.FormatUint(gfibonacci.Previous, 10) + "," + strconv.FormatUint(gfibonacci.Current, 10) + "," + strconv.FormatUint(gfibonacci.Next, 10)

		xrecordString := strconv.FormatUint(f["Hitcount"], 10) + "," + strconv.FormatUint(f["Position"], 10) + "," + strconv.FormatUint(f["Previous"], 10) + "," + strconv.FormatUint(f["Current"], 10) + "," + strconv.FormatUint(f["Next"], 10)

		output := []byte(xrecordString)
		mutex.Lock()
		openFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		failOnError(err, "An error was encountered saving fibo to file")

		err = ioutil.WriteFile(filename, output, 0644)
		openFile.Close()
		mutex.Unlock()
		//fmt.Printf("File saved to %v \n", filename)
		time.Sleep(time.Second)
	}
}

//read json file
//func readFibo() (map[string]uint64, error) {
func readFibo() {
	dataFile, err := os.Open("fibo.csv")
	if err != nil {
		fmt.Println("Could not find fibo.csv creating...")
		saveFibo()
		time.After(time.Second / 10)
		//readFibo()
	}
	//failOnError(err, "Could not read fibo.csv")

	reader := csv.NewReader(bufio.NewReader(dataFile))
	line, err := reader.Read()
	//failOnError(err, "Cannot read file")
	//check for end of file
	if err == io.EOF {
		fmt.Println("Empty structureless fibo.csv found. Re-creating...")
		saveFibo()
		time.After(time.Second / 10)
		readFibo()
		//break
	}

	//if err != nil {
	//	fmt.Printf("Cannot read file: %v following error occured: %v \n", fileSource, err.Error())
	//}
	//line0, err := strconv.Atoi(line[0])
	line0, err := strconv.ParseUint(string(line[0]), 0, 64)
	failOnError(err, "Could not read or convert to int")

	line1, err := strconv.ParseUint(string(line[1]), 0, 64)
	failOnError(err, "Could not read or convert to int")

	//line2, err := strconv.Atoi(line[2])
	line2, err := strconv.ParseUint(string(line[2]), 0, 64)
	failOnError(err, "Could not read or convert to int")

	//line3, err := strconv.Atoi(line[3])
	line3, err := strconv.ParseUint(string(line[3]), 0, 64)
	failOnError(err, "Could not read or convert to int")

	//line4, err := strconv.Atoi(line[4])
	line4, err := strconv.ParseUint(string(line[4]), 0, 64)
	failOnError(err, "Could not read or convert to int")

	go func() {

		mutex.Lock()

		FiboStore["Hitcount"] = uint64(line0)
		FiboStore["Position"] = uint64(line1)
		FiboStore["Previous"] = uint64(line2)
		FiboStore["Current"] = uint64(line3)
		FiboStore["Next"] = uint64(line4)

		mutex.Unlock()

	}()

	//return FiboStore, nil
}

func main() {
	//Read data (csv)store into memory
	go readFibo()
	//failOnError(err, "Could not read csv file store")
	//start the fibonacci saving process
	go saveFibo()

	//Start Gorrila mux
	router := mux.NewRouter()
	//All crud handlers that will be needed by front ends will be defined and handled here.
	//User activity
	//http.HandleFunc("/", homepage)
	router.HandleFunc("/", homepage).Methods("GET")
	router.HandleFunc("/previous", previous).Methods("GET")
	router.HandleFunc("/current", current).Methods("GET")
	router.HandleFunc("/next", next).Methods("GET")
	router.HandleFunc("/reset", reset).Methods("GET")
	//router.Use(loggingMiddleware)
	//start serving handles

	if len(os.Args) > 0 {
		port = ":" + string(os.Args[1]) //strconv.Atoi(os.Args[1])
	}

	fmt.Printf("Starting server at port %s \n", port)
	log.Fatal(http.ListenAndServe(port, router))

}

func appCleanup() {
	log.Println("Cleaning up before exit")
	saveFibo()
	os.Exit(1)
}
