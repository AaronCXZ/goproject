package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	threads, err := data.Thread()
	if err == nil {
		_, err := seesion(w, r)
		if err != nil {
			generateHTML(w, threads, "layout", "public.navbar", "index")
		} else {
			generateHTML(w, threads, "layout", "private.navbar", "index")
		}
	}
}
