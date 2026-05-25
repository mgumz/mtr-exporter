package main

import (
	"io"
	"net/http"
)

func mtrHealthPage(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "OK\n")
}
