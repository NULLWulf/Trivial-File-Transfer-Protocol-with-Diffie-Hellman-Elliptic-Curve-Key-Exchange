package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func RunClientMode() {
	router := httprouter.New() // Create HTTP router
	router.GET("/", homepage)  // Services index.html
	router.GET("/getImage", getImage)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Failed to Listen and Serve: ", err)
		return
	}
}

func homepage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Serving homepage")
	http.ServeFile(w, r, "./html/index.html")
	return
}

func getImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	imageUrl := r.URL.Query().Get("url")
	log.Printf("Serving image: %s\n", imageUrl)

	var img []byte
	client, err := NewTFTPClient()
	if err != nil {
		log.Printf("Error Creating TFTP Client: %s\n", err)
		return
	}
	defer client.Close()
	img, err = client.RequestFile(imageUrl)
	if err != nil {
		log.Printf("Error Requesting File over TFTP: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	_, err = w.Write(img)
	return
}
