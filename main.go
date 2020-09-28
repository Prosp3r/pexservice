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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

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
func fibonacciBack() {
	//backshift
	var newPrevious uint64
	var newNext uint64
	var newCurrent uint64
	var newPosition uint64
	if gfibonacci.Previous >= 1 {
		newPrevious = gfibonacci.Current - gfibonacci.Previous
		newCurrent = gfibonacci.Previous
		newNext = gfibonacci.Current
		newPosition = gfibonacci.Position - 1
	} else {
		newPrevious = 0 //gfibonacci.Previous
		newPosition = 1
		newCurrent = newPrevious
		newNext = newCurrent + 1
	}

	newfibonacci := Fibonacci{
		Hitcount: gfibonacci.Hitcount + 1,
		Position: newPosition,
		Previous: newPrevious, //we moving back now
		Current:  newCurrent,
		Next:     newNext,
	}
	gfibonacci = newfibonacci
	//return current - previous

	go func() {
		mutex.Lock()
		defer mutex.Unlock()
		FiboStore["Hitcount"]++
		FiboStore["Position"] = newPosition
		FiboStore["Current"] = newCurrent
		FiboStore["Previous"] = newPrevious
		FiboStore["Next"] = newNext

	}()
}

//fibonacciGo sets the next set of numbers in the sequence //0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 54
func fibonacciGo() {
	//backshift
	var newPrevious uint64
	var newNext uint64
	var newCurrent uint64
	var newPosition uint64
	if gfibonacci.Current <= 0 {
		newPrevious = 0 //gfibonacci.Previous
		newPosition = 1
		newCurrent = newPrevious + 1
		newNext = newCurrent + newPrevious
	} else {

		newPrevious = gfibonacci.Current
		newCurrent = gfibonacci.Next
		newNext = newCurrent + newPrevious
		newPosition = gfibonacci.Position + 1
	}

	/*newfibonacci := Fibonacci{
		Hitcount: gfibonacci.Hitcount + 1,
		Position: newPosition,
		Previous: newPrevious, //we moving back now
		Current:  newCurrent,
		Next:     newNext,
	}*/

	go func() {

		mutex.Lock()
		defer mutex.Unlock()
		FiboStore["Hitcount"]++
		FiboStore["Position"] = newPosition
		FiboStore["Current"] = newCurrent
		FiboStore["Previous"] = newPrevious
		FiboStore["Next"] = newNext

	}()

	//gfibonacci = newfibonacci
	//return current - previous
}

func previous(w http.ResponseWriter, r *http.Request) {
	fibonacciBack()
	//previous := gfibonacci.Previous
	mutex.Lock()
	defer mutex.Unlock()
	fib, _ := json.Marshal(FiboStore)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fib)
	//json.NewEncoder(w).Encode(previous)
}

func current(w http.ResponseWriter, r *http.Request) {
	//current := gfibonacci.Current
	//current := FiboStore["Current"]
	mutex.Lock()
	defer mutex.Unlock()
	fib, _ := json.Marshal(FiboStore)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(fib)
	json.NewEncoder(w).Encode(fib)
}

func next(w http.ResponseWriter, r *http.Request) {
	fibonacciGo()
	//next := gfibonacci.Next
	mutex.Lock()
	defer mutex.Unlock()
	fib, _ := json.Marshal(FiboStore)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fib)
	//json.NewEncoder(w).Encode(next)
}

//reset restoes all information to zero
func reset(w http.ResponseWriter, r *http.Request) {
	/*newfibonacci := Fibonacci{
		Hitcount: gfibonacci.Hitcount + 1,
		Position: 1,
		Previous: 0, //we moving back now
		Current:  0,
		Next:     1,
	}
	gfibonacci = newfibonacci
	fib := gfibonacci*/

	go func() {

		mutex.Lock()
		defer mutex.Unlock()
		FiboStore["Hitcount"]++
		FiboStore["Position"] = 1
		FiboStore["Previous"] = 0
		FiboStore["Current"] = 0
		FiboStore["Next"] = 1

	}()

	mutex.Lock()
	defer mutex.Unlock()
	fib, _ := json.Marshal(FiboStore)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fib)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, " You reached my fibonacci /homepage \n\n The hits parameter is a traffic hit counter for some easy hit load tracking \n\n The position parameter is your current position on the fibonacci scale \n\n Try the following end points \n /previous to view previous number \n /current to view the current number \n /next to view the next number \n /reset to reset the fibonacci to zero \n ")
}

func main() {
	//Read data (csv)store into memory
	_, err := readFibo()
	failOnError(err, "Could not read csv file store")
	//Start Gorrila mux
	router := mux.NewRouter()
	//All crud handlers that will be needed by front ends will be defined and handled here.
	//User activity
	router.HandleFunc("/", homepage).Methods("GET")
	router.HandleFunc("/previous", previous).Methods("GET")
	router.HandleFunc("/current", current).Methods("GET")
	router.HandleFunc("/next", next).Methods("GET")
	router.HandleFunc("/reset", reset).Methods("GET")

	//start serving handles

	if len(os.Args) > 0 {
		port = ":" + string(os.Args[1]) //strconv.Atoi(os.Args[1])
	}

	fmt.Printf("Starting server at port %s \n", port)
	log.Fatal(http.ListenAndServe(port, router))

}

//save sequence to file
func saveFibo() bool {

	//recordString := strconv.FormatUint(gfibonacci.Hitcount, 10) + "," + strconv.FormatUint(gfibonacci.Position, 10) + "," + strconv.FormatUint(gfibonacci.Previous, 10) + "," + strconv.FormatUint(gfibonacci.Current, 10) + "," + strconv.FormatUint(gfibonacci.Next, 10)
	mutex.Lock()
	defer mutex.Unlock()
	xrecordString := strconv.FormatUint(FiboStore["Hitcount"], 10) + "," + strconv.FormatUint(FiboStore["Position"], 10) + "," + strconv.FormatUint(FiboStore["Previous"], 10) + "," + strconv.FormatUint(FiboStore["Current"], 10) + "," + strconv.FormatUint(FiboStore["Next"], 10)
	//toSave := []string{recordString} //no need for the slice
	//finalOutput, err := endFormats.Convert(f)
	filename := "fibo.csv"
	//output := []byte(recordString)
	output := []byte(xrecordString)
	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	failOnError(err, "An error was encountered saving fibo to file")

	err = ioutil.WriteFile(filename, output, 0644)
	fmt.Printf("File saved to %v \n", filename)
	return true
}

//read json file
func readFibo() (Fibonacci, error) {
	dataFile, err := os.Open("fibo.csv")
	failOnError(err, "Could not read fibo.csv")

	reader := csv.NewReader(bufio.NewReader(dataFile))

	line, err := reader.Read()
	failOnError(err, "Cannot read file")
	//check for end of file
	//if err == io.EOF {
	//	break
	//}
	//if err != nil {
	//	fmt.Printf("Cannot read file: %v following error occured: %v \n", fileSource, err.Error())
	//}
	line0, err := strconv.Atoi(line[0])
	failOnError(err, "Could not read or convert to int")

	line1, err := strconv.Atoi(line[1])
	failOnError(err, "Could not read or convert to int")

	line2, err := strconv.Atoi(line[2])
	failOnError(err, "Could not read or convert to int")

	line3, err := strconv.Atoi(line[3])
	failOnError(err, "Could not read or convert to int")

	line4, err := strconv.Atoi(line[4])
	failOnError(err, "Could not read or convert to int")

	fibonacci := Fibonacci{
		Hitcount: uint64(line0),
		Position: uint64(line1),
		Previous: uint64(line2),
		Current:  uint64(line3),
		Next:     uint64(line4),
	}
	//set global fibonacci
	gfibonacci = fibonacci

	go func() {

		mutex.Lock()
		defer mutex.Unlock()
		FiboStore["Hitcount"]++
		FiboStore["Position"] = uint64(line1)
		FiboStore["Previous"] = uint64(line2)
		FiboStore["Current"] = uint64(line3)
		FiboStore["Next"] = uint64(line4)

	}()

	return fibonacci, nil
}
