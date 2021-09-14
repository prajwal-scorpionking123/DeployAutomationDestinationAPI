package deploycontroller

import (
	"archive/zip"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//checking whether file in that entered path is exist or not
// func FileExists(name string) bool {
// 	if _, err := os.Stat(name); err != nil {
// 		if os.IsNotExist(err) {
// 			return false
// 		}
// 	}
// 	return true
// }
// func uploadMedia(file multipart.File, filename string) {
// 	defer file.Close()
// 	tmpfile, _ := os.Create("../SOURCE/" + filename)
// 	defer tmpfile.Close()
// 	io.Copy(tmpfile, file)
// }

// func getMetadata(r *http.Request) ([]byte, error) {
// 	f, _, err := r.FormFile("metadata")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get metadata form file: %v", err)
// 	}

// 	metadata, errRead := ioutil.ReadAll(f)
// 	if errRead != nil {
// 		return nil, fmt.Errorf("failed to read metadata: %v", errRead)
// 	}

// 	return metadata, nil
// }

// func verifyRequest(r *http.Request) error {
// 	if _, ok := r.MultipartForm.File["media"]; !ok {
// 		return fmt.Errorf("media is absent")
// 	}

// 	if _, ok := r.MultipartForm.File["metadata"]; !ok {
// 		return fmt.Errorf("metadata is absent")
// 	}

// 	return nil
// }
// func DeployFiles(c *gin.Context) {
// 	parseErr := c.Request.ParseMultipartForm(32 << 20)
// 	if parseErr != nil {

// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"err": "failed to parse multipart message",
// 		})
// 		return
// 	}

// 	if c.Request.MultipartForm == nil || c.Request.MultipartForm.File == nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"err": "expecting multipart form file",
// 		})
// 		return
// 	}

// 	if err := verifyRequest(c.Request); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"err": err.Error(),
// 		})
// 		return
// 	}

// 	metadata, errMeta := getMetadata(c.Request)
// 	if errMeta != nil {

// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"err": "failed to get metadata",
// 		})
// 		return
// 	}
// 	log.Print(string(metadata))

// 	for _, h := range c.Request.MultipartForm.File["media"] {
// 		file, err := h.Open()
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"err": "failed to get media form file",
// 			})
// 			return
// 		}
// 		uploadMedia(file, h.Filename)
// 	}
// }
func DeployFiles(c *gin.Context) {
	var backupList []string
	var destination string
	contentType, params, parseErr := mime.ParseMediaType(c.Request.Header.Get("Content-Type"))
	if parseErr != nil || !strings.HasPrefix(contentType, "multipart/") {

		c.JSON(http.StatusBadRequest, gin.H{
			"err": "expecting a multipart message",
		})
		return
	}

	multipartReader := multipart.NewReader(c.Request.Body, params["boundary"])
	defer c.Request.Body.Close()

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "unexpected error when retrieving a part of the message",
			})
			return
		}
		defer part.Close()

		fileBytes, err := ioutil.ReadAll(part)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"err": "failed to read content of the part",
			})
			return
		}

		switch part.Header.Get("Content-ID") {
		case "metadata":
			log.Print(string(fileBytes))
			destination = string(fileBytes)
			log.Print(destination)

		case "media":
			log.Printf("filesize = %d", len(fileBytes))
			log.Println(part.Header.Get("Content-Filepath"))
			isAlreadyThere, _ := isExists(part.Header.Get("Content-Filepath"))
			println(isAlreadyThere)
			if isAlreadyThere {
				backupList = append(backupList, filepath.ToSlash(part.Header.Get("Content-Filepath")))
			}
			f, _ := os.Create(part.Header.Get("Content-Filepath"))
			f.Write(fileBytes)
			f.Close()
		}
	}
	fmt.Println(backupList)
	TakeBackup(backupList, destination)
}
func TakeBackup(backupList []string, destination string) {

	currentTime := time.Now().Format("01-02-2006")
	RandomCrypto, _ := rand.Prime(rand.Reader, 128)
	err := os.MkdirAll("../BACKUP/"+currentTime+"/"+destination+"/", 0755)
	// Get a Buffer to Write To
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.Create("../BACKUP/" + currentTime + "/" + destination + "/" + RandomCrypto.String() + "_backup.zip")
	if err != nil {
		log.Println("Failed to open zip for writing: %s", err)
	}
	defer file.Close()
	zipw := zip.NewWriter(file)
	defer zipw.Close()
	for _, filename := range backupList {
		if err := appendFiles(filename, zipw); err != nil {
			log.Println("Failed to add file %s to zip: %s", filename, err)
		}
	}
}
func appendFiles(filename string, zipw *zip.Writer) error {
	fmt.Println(filename)
	file, err := os.Open(filename)
	// dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Failed to open %s: %s", filename, err)
	}

	wr, err := zipw.Create(filepath.Base(filename))
	if err != nil {
		msg := "Failed to create entry for %s in zip file: %s"
		return fmt.Errorf(msg, filename, err)
	}

	if _, err := io.Copy(wr, file); err != nil {
		return fmt.Errorf("Failed to write %s to zip: %s", filename, err)
	}

	return nil
}
func isExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// func addFiles(w *zip.Writer, basePath, baseInZip string) {
// 	// Open the Directory
// 	files, err := ioutil.ReadDir(basePath)
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	for _, file := range files {
// 		fmt.Println(basePath + file.Name())
// 		if !file.IsDir() {
// 			dat, err := ioutil.ReadFile(basePath + file.Name())
// 			if err != nil {
// 				fmt.Println(err)
// 			}

// 			// Add some files to the archive.
// 			f, err := w.Create(baseInZip + file.Name())
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			_, err = f.Write(dat)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 		} else if file.IsDir() {

// 			// Recurse
// 			newBase := basePath + file.Name() + "/"
// 			fmt.Println("Recursing and Adding SubDir: " + file.Name())
// 			fmt.Println("Recursing and Adding SubDir: " + newBase)

// 			addFiles(w, newBase, baseInZip+file.Name()+"/")
// 		}
// 	}
// }
// func ZipWriter(filename string) {
// 	err := os.MkdirAll(filename, 0755)
// 	name := filepath.Base(filename)
// 	// Get a Buffer to Write To
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	srcFile, err := os.Open("../SOURCE/" + filename)
// 	fmt.Println(srcFile)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	outFile, err := os.Create("../BACKUP/" + name)
// 	fmt.Println(outFile)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	files, err := io.Copy(outFile, srcFile)
// 	fmt.Println(files)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer outFile.Close()

// 	// Create a new zip archive.
// 	// w := zip.NewWriter(outFile)

// 	// Add some files to the archive.
// 	// addFiles(w, baseFolder, "")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Make sure to check the error on Close.
// }
func CutSource(source string) string {
	s := strings.ReplaceAll(source, "../SOURCE", "")
	return s
}
