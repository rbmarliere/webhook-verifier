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

func parseHeaders(r *http.Request) ([]string, error) {
	var msg string

	secret := r.Header.Get("Secret")
	if secret == "" {
		msg = "Secret not found"
		logger.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	logger.Println("Found Secret: " + secret)

	project_root := r.Header.Get("Project-Root")
	if project_root == "" {
		msg = "Project-Root not found"
		logger.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	logger.Println("Found Project-Root: " + project_root)

	expected_signature := r.Header.Get("X-Hub-Signature-256")
	if expected_signature == "" {
		msg = "X-Hub-Signature-256 not found"
		logger.Println(msg)
		return nil, fmt.Errorf(msg)
	}
	logger.Println("Found X-Hub-Signature-256: " + expected_signature)

	return []string{secret, project_root, expected_signature}, nil
}

func parseBody(r *http.Request) ([]byte, error) {
	payload, error := io.ReadAll(r.Body)
	if payload == nil || error != nil {
		msg := "Payload not found"
		logger.Println(msg)
		return nil, fmt.Errorf(msg)
	}

	return payload, nil
}

func verifySignature(secret, project_root, expected_signature string, payload []byte) bool {
	_, error := os.Stat(project_root)
	if error != nil {
		logger.Println("Project-Root path does not exist")
		return false
	}

	digest := hmac.New(sha256.New, []byte(secret))
	digest.Write(payload)
	signature := fmt.Sprintf("sha256=%x", digest.Sum(nil))
	logger.Println("Signature: " + expected_signature)

	if signature != expected_signature {
		logger.Println("Signatures do not match")
		return false
	}

	return true
}

func updateProject(project_root string) {
	cmd := exec.Command("git", "-C", project_root, "pull")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	error := cmd.Run()
	logger.Println(out.String())
	if error != nil {
		logger.Println(error)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	headers, error := parseHeaders(r)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, error := parseBody(r)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if verifySignature(headers[0], headers[1], headers[2], payload) {
		updateProject(headers[1])
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func main() {
	port := flag.String("port", "32777", "port to listen on")
	verbose := flag.Bool("verbose", false, "verbose logging")
	flag.Parse()

	logger.verbose = *verbose

	http.HandleFunc("/", handleRequest)
	http.ListenAndServe(":" + *port, nil)
}
