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
    log.Printf("Starting server on port " + port)
    log.Fatal(http.ListenAndServe(":"+port, handler))
}

func Serve(w http.ResponseWriter, r *http.Request) {
	url := pixiv_domain + r.URL.Path
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        http.Error(w, "Error creating request: "+err.Error(), http.StatusInternalServerError)
        log.Printf("Error creating request: %v", err)
        return
    }
    
    // In order to use the image link the referer header must be set
    req.Header.Set("referer", "https://www.pixiv.net/")
    res, err := client.Do(req)
    if err != nil {
        http.Error(w, "Error making request: "+err.Error(), http.StatusInternalServerError)
        log.Printf("Error making request to %s: %v", url, err)
        return
    }
    
    defer res.Body.Close()
    
    data, err := ioutil.ReadAll(res.Body)
    if err != nil {
        http.Error(w, "Error reading response: "+err.Error(), http.StatusInternalServerError)
        log.Printf("Error reading response: %v", err)
        return
    }
    
    for key, header := range res.Header {
        w.Header().Set(key, header[0])
    }
    w.WriteHeader(res.StatusCode)
    w.Write(data)
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
