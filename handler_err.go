package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	ResponseWithError(w, 400, "test error")
}
