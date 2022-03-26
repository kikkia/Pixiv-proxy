package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var pixiv_domain = "https://i.pximg.net/"
var client = &http.Client{}

func main() {
	handler := http.HandlerFunc(Serve)
	port := getenv("SERVER_PORT", "80")
	http.ListenAndServe(":"+port, handler)
	log.Printf("Started server on port " + port)
}

func Serve(w http.ResponseWriter, r *http.Request) {
	url := pixiv_domain + r.URL.Path
	req, _ := http.NewRequest("GET", url, nil)
	// In order to use the image link the referer header must be set
	req.Header.Set("referer", "https://www.pixiv.net/")
	res, _ := client.Do(req)

	// Ream image data
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatalf("ioutil.ReadAll -> %v", err)
	}
	res.Body.Close()

	// Map headers over to new response
	for key, header := range res.Header {
		w.Header().Set(key, header[0])
	}
	w.WriteHeader(res.StatusCode)

	w.Write(data)
	return
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
