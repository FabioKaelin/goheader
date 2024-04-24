package main

import (
	"fmt"
	"goheader/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Any("/*any", func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/goheader") {
			if path == "/goheader/upload" && c.Request.Method == "POST" {
				uploadHandler(c)
			} else {
				handler(c)
			}
		} else if strings.HasPrefix(path, "/gonourl") {
			if path == "/gonourl/upload" && c.Request.Method == "POST" {
				uploadHandler(c)
			} else {
				handler(c)
			}
		} else {
			handler(c)
		}
	})

	fmt.Println("Started 8000")
	r.Run("0.0.0.0:8000")
}

func handler(c *gin.Context) {
	outputString := ""
	outputString += fmt.Sprintf("c.Request.RemoteAddr: %s\n", c.Request.RemoteAddr)
	outputString += fmt.Sprintf("c.RemoteIP() %s\n", c.RemoteIP())
	outputString += fmt.Sprintf("c.Request.Proto %s\n", c.Request.Proto)
	outputString += fmt.Sprintf("c.Request.URL.Scheme %s\n", c.Request.URL.Scheme)
	outputString += fmt.Sprintf("c.Request.Method %s\n", c.Request.Method)
	outputString += fmt.Sprintf("c.Request.URL %s\n", c.Request.URL)
	outputString += fmt.Sprintf("c.Request.URL.Query() %v\n", c.Request.URL.Query())
	outputString += fmt.Sprintf("c.Request.ContentLength %d\n", c.Request.ContentLength)
	outputString += fmt.Sprintf("c.ContentType() %s\n", c.ContentType())
	outputString += fmt.Sprintln("\nHeaders:")
	for k, v := range c.Request.Header {
		outputString += fmt.Sprintf("%s: %v\n", k, v)
	}

	outputString += fmt.Sprintln("\nBody:")
	body, _ := c.GetRawData()
	outputString += fmt.Sprintf("%s\n", body)
	c.String(http.StatusOK, outputString)
	fmt.Println(outputString)
}

var MAX_UPLOAD_SIZE int64 = 1024 * 1024 * 1024 // 1GB

func uploadHandler(c *gin.Context) {
	log := logger.CreateLogger()
	fmt.Println("-------------------New Request-------------------", log.GetRequestID())

	// 32 MB is the default used by FormFile()

	multiPartFrom, err := c.MultipartForm()
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		log.Println("Error:", err)
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
				log.Println("The uploaded image is too big: %s. Please use an image less than 1GB in size", fileHeader.Filename)
				return
			}

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				log.Println(err)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				c.String(http.StatusInternalServerError, err.Error())
				log.Println(err)
				return
			}

			buff = nil

			// filetype := http.DetectContentType(buff)
			// log.Println("Name:", fileHeader.Filename)
			// log.Println("Size:", fileHeader.Size)
			// log.Println("Header:", fileHeader.Header)
			// log.Println("Type:", filetype)
			// log.Println(".....")
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
	log.Println("Upload successful")

	// fmt.Fprintf(w, "Upload successful")
}
