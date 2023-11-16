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

func parseHeaders(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("Secret")
	if secret == "" {
		log.Println("Secret not found")
		return
	}
	log.Println("Found Secret: " + secret)

	project_root := r.Header.Get("Project-Root")
	if project_root == "" {
		log.Println("Project-Root not found")
		return
	}

	log.Println("Found Project-Root: " + project_root)
	_, error := os.Stat(project_root)
	if error != nil {
		log.Println("Project-Root path does not exist")
		return
	}

	expected_signature := r.Header.Get("X-Hub-Signature-256")
	if expected_signature == "" {
		log.Println("X-Hub-Signature-256 not found")
		return
	}
	log.Println("Found X-Hub-Signature-256: " + expected_signature)

	payload, error := io.ReadAll(r.Body)
	if payload == nil || error != nil {
		log.Println("Payload not found")
		return
	}

	digest := hmac.New(sha256.New, []byte(secret))
	digest.Write(payload)
	signature := fmt.Sprintf("sha256=%x", digest.Sum(nil))
	log.Println("Signature: " + expected_signature)

	if signature != expected_signature {
		log.Println("Signatures do not match")
		return
	}

	cmd := exec.Command("git", "-C", project_root, "pull")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	error = cmd.Run()
	if error != nil {
		log.Println(error)
	}
	log.Println(out.String())
}

func main() {
	port := flag.String("port", "32777", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", parseHeaders)
	http.ListenAndServe(":" + *port, nil)
}
