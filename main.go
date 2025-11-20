package main

import (
	"net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.LoadHTMLGlob("templates/*.html")
    
	// set the port to run the system
	r.Run(":3000")
}