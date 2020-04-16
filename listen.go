package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"
)

func main() {

	app := &cli.App{
		Name:  "listen",
		Usage: "listen to a port and log requests",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Value: "4000",
				Usage: "port ",
			},
			&cli.StringFlag{
				Name:  "logPath",
				Value: "",
				Usage: "path to log file",
			},
		},
		Action: func(c *cli.Context) error {
			port := 4000
			if c.NArg() > 0 {
				if p, err := strconv.Atoi(c.Args().Get(0)); err == nil {
					port = p
				}
			}
			start(port, c.String("logPath"))
			return nil
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func start(port int, logPath string) {

	openLogFile(logPath)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	http.HandleFunc("/", rootHandler)

	fmt.Printf("listening on %v\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), logRequest(http.DefaultServeMux))
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		line := fmt.Sprintf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, b)
		fmt.Printf(line)
		log.Printf(line)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		fmt.Printf("Logging to %v\n", logfile)
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			log.Fatal("OpenLogfile", err)
		}
		log.SetOutput(lf)
	}
}
