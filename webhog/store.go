package webhog

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Stored temporary directory for the entity files.
var EntityDir string

// Create a temporary dir to store entity files.
func NewEntityDir() (err error) {
	randName := randFileName(10)
	EntityDir, err = ioutil.TempDir("", randName)

	return err
}

// Generate a random file name on the fly.
func randFileName(str_size int) string {
	alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, str_size)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}

	return string(bytes) + strconv.FormatInt(time.Now().Unix(), 5)
}

// Stores the given js / css / img file into the given tempdir with a temp
// name.
func StoreResource(resource, attr, entDir string) (name string, err error) {
	// random name for the new file
	randName := randFileName(10)

	// full path of the new file to be saved
	finalPath := entDir + "/" + randName + attr

	// final name of the file prepended with ./ to be
	// found in the directory of the final entity folder
	finalName := "./" + randName + attr

	newFile, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}

	defer newFile.Close()

	resp, err := http.Get(resource)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	_, err = io.Copy(newFile, resp.Body)
	if err != nil {
		return "", err
	}

	return finalName, err
}

// Stores the final HTML string into an index.html file.
func StoreHTML(html bytes.Buffer, entDir string) (err error) {
	// default our html file name to index.html
	fileName := "index.html"

	// full path of the html filed to be saved
	finalPath := entDir + "/" + fileName
	log.Println("Ent dir: ", entDir)

	newFile, err := os.Create(finalPath)
	if err != nil {
		return err
	}

	defer newFile.Close()

	_, err = newFile.Write(html.Bytes())
	if err != nil {
		return err
	}

	return err
}

func ArchiveFinalFiles(entDir string) (err error) {
	var buf = new(bytes.Buffer)
	gz := gzip.NewWriter(buf)

	defer gz.Close()

	files, err := ioutil.ReadDir(entDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		data, err := ioutil.ReadFile(entDir + "/" + file.Name())
		if err != nil {
			return err
		}

		_, err = gz.Write(data)
		if err != nil {
			return err
		}

		f, err := os.Create(entDir + "/" + file.Name() + ".gz")
		if err != nil {
			return err
		}

		_, err = io.Copy(f, buf)
		if err != nil {
			return err
		}
	}

	err = gz.Close()

	return err
}
