package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	ResponseWithJSON(w, 200, struct{}{})
}
