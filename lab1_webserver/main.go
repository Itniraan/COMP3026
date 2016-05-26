package main

// Import the net/http package
import (
    "net/http"
)

// Main function
func main() {
    // Initilize the fileserver, point it to the "html" directory as its root
    fileServer := http.FileServer(http.Dir("html"))
    
    // Handle any http requests by sending them to the fileserver
    http.Handle("/", fileServer)
    
    // Set the listening port
    http.ListenAndServe(":8000", nil)
}