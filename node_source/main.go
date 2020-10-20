package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/octago/sflags/gen/gflag"
)

// example: read-args.go -name Frontend -delay 9 -port 9090 -cpu 90 -memory 900 -endpoint-url /read -endpoint-cpu 99 -endpoint-delay 98 -endpoint-url /index -endpoint-cpu 22 -endpoint-delay 202
type config struct {
	Name           string   `flag:"name a" desc:"Server/service name"`
	InitDelay      uint     `flag:"delay" desc:"Delay after start up [ms]"`
	Port           uint     `flag:"port" desc:"Open port to listen"`
	CPUusage       uint     `flag:"cpu" desc:"CPU usage in idle time [mCPU]"`
	MemoryUsage    uint     `flag:"memory" desc:"Memory usage in idle time [kB]"`
	Endpoints      []string `flag:"endpoint-url" desc:"Endpoints to listen"`
	EndpointsCPU   []uint   `flag:"endpoint-cpu" desc:"CPU usage for the endpoints"`
	EndpointsDelay []uint   `flag:"endpoint-delay" desc:"Delay for each endpoint [ms]"`
	EndpointsCall  []string `flag:"endpoint-call" desc:"If the endpoint need to call other service"`
}

func readConfigParameters() *config {
	// Set default parameters
	c := &config{
		Name:        "Service-#ID",
		InitDelay:   0,
		Port:        8080,
		CPUusage:    50,
		MemoryUsage: 64,
		Endpoints: []string{
			"/",
			"/health",
		},
		EndpointsCPU: []uint{
			200,
			10,
		},
		EndpointsDelay: []uint{
			30,
			0,
		},
	}

	err := gflag.ParseToDef(c)
	if err != nil {
		log.Fatalf("[READ_PARAMS]\terr: %v", err)
	}
	flag.Parse()

	// Check given paramters
	fmt.Printf("[READ_PARAMS]\tParameters OK: %t\n", c.check())

	return c
}

func (c *config) check() bool {
	// endpoint-call out?
	//for i, e := range c.EndpointsCall {
	//	fmt.Println(i, " --> ", e)
	//}

	if len(c.Endpoints) == len(c.EndpointsCPU) &&
		len(c.Endpoints) == len(c.EndpointsDelay) &&
		len(c.Endpoints) == len(c.EndpointsCall) {
		return true
	}
	return false
}

func main() {
	var cfg = readConfigParameters()

	fmt.Printf("[MAIN]\t\tConfig values:\t%+v\n", cfg)

	var addr string = ":" + fmt.Sprint(cfg.Port)

	// start webserver
	for _, endpoint := range cfg.Endpoints {
		// fmt.Printf("%d --> %s --> %d\n", i, endpoint, cfg.EndpointsCPU[i])
		http.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			//fmt.Fprintf(w, "Req: %s --> %s\n", r.URL, r.URL.Path)
			for k, endp := range cfg.Endpoints {
				if endp == r.URL.Path {
					time.Sleep(time.Duration(cfg.EndpointsDelay[k]) * time.Millisecond)
					fmt.Printf("[REQUEST]\t%s\n", endp)
					fmt.Fprintf(w, "<h1>Hello from: %s!</h1><hl>\n", cfg.Name)
					fmt.Fprintln(w, "<h3>Config</h3><ul>")
					fmt.Fprintf(w, "<li>Config values: %t</li>\n", cfg.check())
					fmt.Fprintf(w, "<li>Endpoint: %s</li>\n", endp)
					fmt.Fprintf(w, "<li>CPU usage: %d</li>\n", cfg.EndpointsCPU[k])
					fmt.Fprintf(w, "<li>Delay time: %d</li>\n", cfg.EndpointsDelay[k])
					fmt.Fprintf(w, "<li>Call out: %s</li>\n", cfg.EndpointsCall[k])
					//fmt.Fprintf(w, "\t(%s)\n", strings.Split(cfg.EndpointsCall[k], ";"))
					if cfg.EndpointsCall[k] != "pass" {
						fmt.Fprintf(w, "<ul>")
						for i, callOut := range strings.Split(cfg.EndpointsCall[k], "__") {
							fmt.Printf("[CALL_OUT]\t#no%d --> %s\n", i, callOut)
							url := "http://" + callOut
							resp, err := http.Get(url)

							if err != nil {
								fmt.Fprintf(w, "<li>%d: <b>%s</b>: Oops, something went wrong</li>\n", i, callOut)
							} else {
								fmt.Fprintf(w, "<li>%d: <b>%s</b>: %s</li>\n", i, callOut, resp.Status)
							}

						}
						fmt.Fprintf(w, "</ul>")
					}
					fmt.Fprintln(w, "</ul>")

					fmt.Fprintln(w, "<h3>Info</h3><ul>")
					fmt.Fprintf(w, "<li>Time: %s</li>\n", time.Now())
					fmt.Fprintf(w, "<li>Method: %s</li>\n", r.Method)
					fmt.Fprintf(w, "<li>URL: %s</li>\n", r.URL)
					fmt.Fprintf(w, "<li>RemoteAddr: %s</li>\n", r.RemoteAddr)
					fmt.Fprintf(w, "<li>Host: %s</li>\n", r.Host)
					fmt.Fprintln(w, "</ul>")
				}

			}

			fmt.Fprintf(w, "Response time: %s", time.Now().Sub(start))
		})
	}

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
