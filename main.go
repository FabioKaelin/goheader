package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.GET("/", handler)

	goheaderGroup := r.Group("/goheader")
	{
		goheaderGroup.Any("", handler)
		goheaderGroup.POST("/upload", uploadHandler)
	}

	gonourlGroup := r.Group("/gonourl")
	{
		gonourlGroup.Any("", handler)
		gonourlGroup.POST("/upload", uploadHandler)
	}

	fmt.Println("Started 8000")
	r.Run("0.0.0.0:8000")
}

// handler echoes the HTTP request.
func handler(c *gin.Context) {
	c.String(http.StatusOK, "%s %s %s\n", c.Request.Method, c.Request.URL, c.Request.Proto)
	for k, v := range c.Request.Header {
		c.String(http.StatusOK, "Header[%q] = %q\n", k, v)
	}
	c.String(http.StatusOK, "Host = %q\n", c.Request.Host)
	c.String(http.StatusOK, "RemoteAddr = %q\n", c.Request.RemoteAddr)
	if err := c.Request.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range c.Request.Form {
		c.String(http.StatusOK, "Form[%q] = %q\n", k, v)
	}
	now := time.Now().String()
	fmt.Printf("%s | Request: %s %s %s\n", now, c.Request.Method, c.Request.URL, c.Request.Proto)
}

var MAX_UPLOAD_SIZE int64 = 1024 * 1024 * 1024 // 1GB

func randomString(length int) string {
	b := make([]byte, length+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : length+2]
}

func uploadHandler(c *gin.Context) {
	requestIdentifier := randomString(2)
	fmt.Println("-------------------New Request-------------------", requestIdentifier)

	// 32 MB is the default used by FormFile()

	multiPartFrom, err := c.MultipartForm()
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		fmt.Println(requestIdentifier, "Error:", err)
		return
	}
	// Get a reference to the fileHeaders.
	// They are accessible only after ParseMultipartForm is called
	// fmt.Println("Request received")

	for _, files := range multiPartFrom.File {
		// for key, files := range r.MultipartForm.File {

		// fmt.Printf("%d Files in %s\n", len(files), key)

		for _, fileHeader := range files {
			// Restrict the size of each uploaded file to 1MB.
			// To prevent the aggregate size from exceeding
			// a specified value, use the http.MaxBytesReader() method
			// before calling ParseMultipartForm()
			if fileHeader.Size > MAX_UPLOAD_SIZE {
				c.String(http.StatusBadRequest, "The uploaded image is too big: %s. Please use an image less than 1GB in size", fileHeader.Filename)
				fmt.Println(requestIdentifier, "The uploaded image is too big: %s. Please use an image less than 1GB in size", fileHeader.Filename)
				return
			}

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				fmt.Println(requestIdentifier, err)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				fmt.Println(requestIdentifier, err)
				return
			}

			buff = nil

			// filetype := http.DetectContentType(buff)
			// fmt.Println("Name:", fileHeader.Filename)
			// fmt.Println("Size:", fileHeader.Size)
			// fmt.Println("Header:", fileHeader.Header)
			// fmt.Println("Type:", filetype)
			// fmt.Println(".....")
			// filetype := http.DetectContentType(buff)
			// if filetype != "image/jpeg" && filetype != "image/png" {
			// 	http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
			// 	return
			// }

			// _, err = file.Seek(0, io.SeekStart)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// err = os.MkdirAll("./uploads", os.ModePerm)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusInternalServerError)
			// 	return
			// }

			// f, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }

			// defer f.Close()

			// _, err = io.Copy(f, file)
			// if err != nil {
			// 	http.Error(w, err.Error(), http.StatusBadRequest)
			// 	return
			// }
		}
	}
	c.String(http.StatusOK, "Upload successful")
	fmt.Println(requestIdentifier, "Upload successful")

	// fmt.Fprintf(w, "Upload successful")
}
