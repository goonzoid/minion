package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"syscall"
)

var port = flag.Int("port", 8080, "port to listen on")

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cmdJson := struct {
			Path string `json:"path"`
		}{}

		if err := json.NewDecoder(r.Body).Decode(&cmdJson); err != nil {
			http.Error(w, "unable to parse request json", http.StatusInternalServerError)
			return
		}

		cmd := exec.Command(cmdJson.Path)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		exitStatus := 0 // i'm an optimist

		err := cmd.Run()
		if err != nil {
			switch e := err.(type) {
			case *exec.ExitError: // command exited non-zero
				sys := e.Sys().(syscall.WaitStatus)
				exitStatus = sys.ExitStatus()
			default: // I/O error or something?
				http.Error(w, "failed to execute command: "+e.Error(), http.StatusInternalServerError)
				return
			}
		}

		responseJson := struct {
			Stdout     string
			Stderr     string
			ExitStatus int
		}{
			Stdout:     stdout.String(),
			Stderr:     stderr.String(),
			ExitStatus: exitStatus,
		}
		json.NewEncoder(w).Encode(responseJson)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
