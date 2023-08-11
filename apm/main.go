package main

import (
	"fmt"
	"go.elastic.co/apm/module/apmhttp"
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", Invoke)

	if err := http.ListenAndServe(":8080", apmhttp.Wrap(mux)); err != nil {
		log.Fatalln(err)
	}
}

func Invoke(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`Hello: %s`, time.Now())))
}
