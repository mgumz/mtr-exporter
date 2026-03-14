package main

import (
	"io"
	"net/http"
)

func mtrIndexPage(w http.ResponseWriter, r *http.Request) {

	const txt = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<title>mtr-exporter</title>
</head>
<body>
	mtr-exporter - <a href="https://github.com/mgumz/mtr-exporter">https://github.com/mgumz/mtr-exporter<a><br>
	see <a href="/metrics">/metrics</a>.
</body>`

	io.WriteString(w, txt)
}
