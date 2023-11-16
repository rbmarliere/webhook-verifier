package main

import (
	"log"
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

	project_root := r.Header.Get("Project-Root")
	if project_root == "" {
		log.Println("Project-Root not found")
		return
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		log.Println("X-Hub-Signature-256 not found")
		return
	}
	log.Println(signature)

	payload, error := io.ReadAll(r.Body)
	if error != nil {
		return
	}

	digest := hmac.New(sha256.New, []byte(secret))
	digest.Write(payload)
	expected := fmt.Sprintf("sha256=%x", digest.Sum(nil))
	log.Println(expected)

	if expected != signature {
		log.Println("Signatures do not match")
		return
	}

	cmd := exec.Command("git", "-C", project_root, "pull")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	log.Println(out.String())
	if err != nil {
		log.Println(err)
	}
}

func main() {
	port := flag.String("port", "32777", "port to listen on")
	flag.Parse()

	http.HandleFunc("/", parseHeaders)
	http.ListenAndServe(":" + *port, nil)
}
