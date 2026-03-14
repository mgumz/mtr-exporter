package main

import (
	"io"
	"net/http"
)

func mtrHealthPage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK\n")
}
