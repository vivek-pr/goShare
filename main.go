package main

import (
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl", gin.H{})
    })

	router.POST("/upload", func(c *gin.Context){
		file, header, err := c.Request.FormFile("file")
		if err != nil{
			c.String(http.StatusBadRequest, "Bad Request")
			return
		}
		defer file.Close()

		dest, err := os.Create(header.Filename)
		if err != nil{
			c.String(http.StatusInternalServerError, "Internal Server Error")
			return
		}
		defer dest.Close()

		if _, err = io.Copy(dest, file); err  != nil {
			c.String(http.StatusInternalServerError, "Internal Server error")
			return
		}
		c.String(http.StatusOK, "File Uploaded")

	})

	router.GET("/download/:filename", func(c *gin.Context){
		filename := c.Param("filename")
		if _, err := os.Stat(filename); os.IsNotExist((err)){
			c.String(http.StatusNotFound, "File not found")
			return
		}

		c.Header("Content-Disposition", "attachment; filename="+filename)
		c.Header("Content-Type", "application/octet-stream")
		c.File(filename)
	})

	router.Static("/files", "./")
	router.Run(":8080")
}
