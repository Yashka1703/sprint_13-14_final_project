package server

import (
	"fmt"
	"log"
	"net/http"
)

func StartServ() error {
	webDir := "./web"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, webDir+"/index.html")
			return
		}
		http.ServeFile(w, r, webDir+r.URL.Path)
	})

	port := 7540
	log.Printf("Server port: %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
