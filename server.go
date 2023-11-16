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
)

func parseHeaders(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("Secret")
	if secret == "" {
		log.Fatal("Secret not found")
	}

	project_root := r.Header.Get("Project-Root")
	if project_root == "" {
		log.Fatal("Project-Root not found")
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		log.Fatal("X-Hub-Signature-256 not found")
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
		log.Fatal("Signatures do not match")
	}

	cmd := exec.Command("git", "-C", project_root, "pull")
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	log.Println(out.String())
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/", parseHeaders)
	http.ListenAndServe(":8080", nil)
}
