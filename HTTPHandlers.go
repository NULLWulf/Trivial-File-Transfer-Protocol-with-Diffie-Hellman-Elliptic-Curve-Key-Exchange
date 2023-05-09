package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

// Mutex to lock thread

// RunClientMode starts the client mode
func RunClientMode() {
	router := httprouter.New() // Create HTTP router
	router.GET("/", homepage)  // Services index.html
	router.GET("/getImage", getImage2)

	log.Printf("Listening on port 40500\n")
	err := http.ListenAndServe(":8080", router) // ports for serving web page
	if err != nil {
		log.Fatal("Failed to Listen and Serve: ", err)
		return
	}
}

// homepage serves the index.html file
func homepage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("Serving homepage")
	http.ServeFile(w, r, "./html/index.html")
	return
}

// getImage serves the image
func getImage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Lock thread
	imageUrl := r.URL.Query().Get("url")
	log.Printf("Serving image: %s\n", imageUrl)

	client, err := NewTFTPClient() // instantiate a new TFTP client
	if err != nil {
		log.Printf("Error Creating TFTP Client: %s\n", err)
		return
	}
	defer client.Close()
	err, img, _ := client.RequestFile(imageUrl) // request the file via url
	log.Printf("Image Size: %d\n", len(img))
	if err != nil {
		log.Printf("Error Requesting File over TFTP: %s\n", err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg") // set the content type
	n, err := w.Write(img)
	log.Printf("Serving image of size: %d\n", n)
	// Save it to client files folder
	err = os.WriteFile("./client_files/"+imageUrl, img, 0644)
	if err != nil {
		log.Printf("Error saving file: %s\n", err)
		return
	}
	return
}

func getImage2(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Lock thread
	imageUrl := r.URL.Query().Get("url")
	log.Printf("Serving image: %s\n", imageUrl)

	client, err := NewTFTPClient() // instantiate a new TFTP client
	if err != nil {
		log.Printf("Error Creating TFTP Client: %s\n", err)
		return
	}
	defer client.Close()

	var img []byte
	var numTries int
	for numTries < 10 {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Panic occurred while requesting file over TFTP: %v\n", r)
				}
			}()
			for !client.SendKeyPair() {
				log.Printf("Failed to send key pair, retrying...\n")
			}
			err, img, _ = client.RequestFile(imageUrl) // request the file via url
			if err != nil {
				log.Printf("Error Requesting File over TFTP: %s\n", err)
				panic(err)
			}
		}()

		if err == nil {
			break
		}
		log.Printf("Retry #%d requesting file %s...\n", numTries+1, imageUrl)
		numTries++
	}

	if err != nil {
		log.Printf("Max retries reached, failed to request file %s: %s\n", imageUrl, err)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg") // set the content type
	n, err := w.Write(img)
	log.Printf("Serving image of size: %d\n", n)

	// Save it to client files folder
	err = os.WriteFile("./client_files/"+imageUrl, img, 0644)
	if err != nil {
		log.Printf("Error saving file: %s\n", err)
		return
	}
	return
}
