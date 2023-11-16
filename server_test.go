// SPDX-License-Identifier: GPL-2.0

package main

import (
	"net/http"
	"testing"
	"bytes"
)

func TestParseHeaders(t *testing.T) {
	var error error
	data := []byte("")
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(data))

	_, error = parseHeaders(req)
	if error == nil {
		t.Fatalf("parseHeaders(req) = %v, want %v", error, "error")
	}

	req.Header.Set("Secret", "testing")
	_, error = parseHeaders(req)
	if error == nil {
		t.Fatalf("parseHeaders(req) = %v, want %v", error, "error")
	}

	req.Header.Set("Project-Root", "/tmp")
	_, error = parseHeaders(req)
	if error == nil {
		t.Fatalf("parseHeaders(req) = %v, want %v", error, "error")
	}

	req.Header.Set("X-Hub-Signature-256", "testing")
	_, error = parseHeaders(req)
	if error != nil {
		t.Fatalf("parseHeaders(req) = %v, want %v", error, "error")
	}
}

func TestParseBody(t *testing.T) {
	data := []byte("")
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(data))
	_, error := parseBody(req)
	if error != nil {
		t.Fatalf("parseBody(req) = %v, want %v", error, "error")
	}
}

func TestVerifySignature(t *testing.T) {
	var ret bool

	ret = verifySignature("secret", "expected_signature", []byte("payload"))
	if ret {
		t.Fatalf("verifySignature([]byte(\"secret\"), \"expected_signature\", []byte(\"payload\")) = %v, want %v", ret, true)
	}

	ret = verifySignature("It's a Secret to Everybody", "sha256=757107ea0eb2509fc211221cce984b8a37570b6d7586c22c46f4379c8b043e17", []byte("Hello, World!"))
	if !ret {
		t.Fatalf("verifySignature([]byte(\"It's a Secret to Everybody\"), \"sha256=757107ea0eb2509fc211221cce984b8a37570b6d7586c22c46f4379c8b043e17\", []byte(\"Hello, World!\")) = %v, want %v", ret, true)
	}
}
