package deploycontroller

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

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

		case "media":
			log.Printf("filesize = %d", len(fileBytes))
			log.Println(part.Header.Get("Content-Filepath"))
			f, _ := os.Create(part.Header.Get("Content-Filepath"))
			f.Write(fileBytes)
			f.Close()
		}
	}
}

// func isExists(name string) (bool, error) {
// 	_, err := os.Stat("../PRODUCTION/" + name)
// 	if err == nil {
// 		return true, nil
// 	}
// 	if errors.Is(err, os.ErrNotExist) {
// 		return false, nil
// 	}
// 	return false, err
// }
// func takeBackup(filename string) {
// 	ZipWriter(filename)
// }
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
// 	baseFolder := "../PRODUCTION/"
// 	output := "../BACKUP/bamu"

// 	err := os.MkdirAll(output, 0755)
// 	// Get a Buffer to Write To
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	outFile, err := os.Create("../BACKUP/" + filename + ".zip")
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	defer outFile.Close()

// 	// Create a new zip archive.
// 	w := zip.NewWriter(outFile)

// 	// Add some files to the archive.
// 	addFiles(w, baseFolder, "")

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// Make sure to check the error on Close.
// 	err = w.Close()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
