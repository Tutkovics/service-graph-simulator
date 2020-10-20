package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type arrayPortsFlag struct {
	ports []uint64
}
type endpoint struct {
	name     string //name of the endpoint
	url      string //string to endpoint
	cpuUsage uint   //CPU usage in mCPU
	delay    uint   //delay time (ms) for endpoint
}

var (
	endpoints   []endpoint // array for the endpoints
	nodeName    string     // name for the node/system
	cpuUsage    int        // percent of cpu usage in idle time
	memoryUsage int64      // memory usage in kB in idle time
	processTime int        // process time before first call-out
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz"
	namePrefix  = "node"
)

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (i *arrayPortsFlag) String() string {
	return fmt.Sprint(i.ports)
}

func (i *arrayPortsFlag) Set(value string) error {
	ports := strings.Split(value, ",")
	for _, item := range ports {
		a, _ := strconv.ParseUint(string(item), 10, 64)
		i.ports = append(i.ports, a)
	}
	return nil
}

func (i *arrayPortsFlag) printPortsToString() string {
	var str string
	for j, port := range i.ports {
		if j != 0 {
			str = str + ", " + fmt.Sprint(port)
		} else {
			str = fmt.Sprint(port)
		}
	}
	// fmt.Println(str)
	return str
}

func initFlags() {
	defaultName := namePrefix + "-<HASH>"
	flag.StringVar(&nodeName, "name", defaultName, "Name of the node")

	flag.IntVar(&cpuUsage, "cpu", 10, "Set up the cpu usage for node")
	flag.Int64Var(&memoryUsage, "memory", 1000, "Set up the memory usage in kB")
	flag.IntVar(&processTime, "process", 20, "Set up the process time before the first call-out (ms)")
	flag.Parse()

	if nodeName == namePrefix+"-<HASH>" {
		nodeName = namePrefix + "-" + randStringBytes(16)
	}

	// print the parameters
	fmt.Printf("Arguments\n")
	fmt.Printf("\tName: %s\n", nodeName)
	fmt.Printf("\tCPU: %d %% \n", cpuUsage)
	fmt.Printf("\tMemory: %d (kB)\n", memoryUsage)
	fmt.Printf("\tProcess Time: %d (ms)\n", processTime)
}

func main2() {
	initFlags()
	// allocateMemory()
	time.Sleep(20 * time.Second)
	// allocateMemory()
	openPorts()

	callOut()
}

// func allocateMemory() {
// 	fmt.Println("Allocate memory")
// 	a := make([]int8, 0, 999999999)
// 	overallMemory = append(overallMemory, a)
// 	printMemoryUsage()
// }

func responseHandler(w http.ResponseWriter, r *http.Request) {
	// callOut()
	// time.Sleep()
	w.Write([]byte("Hello from " + nodeName))
	//fmt.Fprintf(w, "Hi from %s", nodeName)

}

// create a simple webserver listen in the given port
// and add to waitGroup to dont escape the application
func createServer(port uint64, group *sync.WaitGroup) {

	// fmt.Println("Try to create new server")
	var addr string = ":" + fmt.Sprint(port)

	// configure server
	s := &http.Server{
		Addr:           addr,
		Handler:        http.HandlerFunc(responseHandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start the server in goroutine
	go func() {
		fmt.Printf("Starting server at port %d\n", port)
		fmt.Println(s.ListenAndServe())
		group.Done()
	}()

}

func openPorts() {
	wg := new(sync.WaitGroup)
	listenPorts := arrayPortsFlag{ports: nil}
	wg.Add(len(listenPorts.ports))

	for _, port := range listenPorts.ports {
		createServer(port, wg)
		printMemoryUsage()
	}

	wg.Wait()
	//time.Sleep(60 * time.Second)
	//fmt.Println("Program ended")
}

func callOut() {
	//fmt.Println("Callout function")
	waitResponseTime()
}

func waitResponseTime() {
	// time.Sleep(time.Duration(responseTime) * time.Millisecond)
}

// return the total Mbytes of memory obtained from the OS.
// details: https://golang.org/pkg/runtime/#MemStats
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("(Alloc) > %d MiB\n", m.Alloc)
	fmt.Printf("(TotalAlloc) > %d MiB\n", m.TotalAlloc)
	fmt.Printf("(Sys) > %d MiB\n", m.Sys)
	fmt.Printf("(Mallocs) > %d MiB\n", m.Mallocs)
	fmt.Printf("(HeapAlloc) > %d MiB\n", m.HeapAlloc)
	fmt.Printf("(HeapSys) > %d MiB\n", m.HeapSys)
	fmt.Printf("(StackSys) > %d MiB\n", m.StackSys)

	return m.StackSys / 1024 / 1024
}

// print out current reserved memory
func printMemoryUsage() {
	fmt.Printf(" () > %d MiB\n", getMemoryUsage())
}
