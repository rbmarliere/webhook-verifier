package main

import (
	"log"
	"os"
	"os/exec"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
	"flag"
)

type custom_logger struct {
	verbose bool
	*log.Logger
}

func (l *custom_logger) Println(v ...interface{}) {
	if l.verbose {
		l.Logger.Println(v...)
	}
}

var logger = &custom_logger{false, log.New(os.Stdout, "", 0)}

func parseHeaders(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("Secret")
	if secret == "" {
		logger.Println("Secret not found")
		return
	}
	logger.Println("Found Secret: " + secret)

	project_root := r.Header.Get("Project-Root")
	if project_root == "" {
		logger.Println("Project-Root not found")
		return
	}

	logger.Println("Found Project-Root: " + project_root)
	_, error := os.Stat(project_root)
	if error != nil {
		logger.Println("Project-Root path does not exist")
		return
	}

	expected_signature := r.Header.Get("X-Hub-Signature-256")
	if expected_signature == "" {
		logger.Println("X-Hub-Signature-256 not found")
		return
	}
	logger.Println("Found X-Hub-Signature-256: " + expected_signature)

	payload, error := io.ReadAll(r.Body)
	if payload == nil || error != nil {
		logger.Println("Payload not found")
		return
	}

	digest := hmac.New(sha256.New, []byte(secret))
	digest.Write(payload)
	signature := fmt.Sprintf("sha256=%x", digest.Sum(nil))
	logger.Println("Signature: " + expected_signature)

	if signature != expected_signature {
		logger.Println("Signatures do not match")
		return
	}

	cmd := exec.Command("git", "-C", project_root, "pull")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	error = cmd.Run()
	if error != nil {
		logger.Println(error)
	}
	logger.Println(out.String())
}

func main() {
	port := flag.String("port", "32777", "port to listen on")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	logger.verbose = *verbose

	http.HandleFunc("/", parseHeaders)
	http.ListenAndServe(":" + *port, nil)
}
