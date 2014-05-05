package webhog

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// Create a tar.gz compressed dir and add in found files
// for upload.
func ArchiveFinalFiles(entDir string) (string, error) {
	lastEl := strings.Split(entDir, "/")
	zippedName := lastEl[len(lastEl)-1] + ".tar.gz"
	finalDir := strings.Join(lastEl[0:len(lastEl)-1], "/")

	fw, err := os.Create(finalDir + "/" + zippedName)
	if err != nil {
		return "", err
	}

	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	dir, err := ioutil.ReadDir(entDir)
	if err != nil {
		return "", err
	}

	for _, file := range dir {
		curPath := entDir + "/" + file.Name()
		err := writeTar(curPath, tw, file)
		if err != nil {
			return "", err
		}
	}

	return finalDir + "/" + zippedName, err
}

// Write the files into the tar dir
func writeTar(path string, tw *tar.Writer, fi os.FileInfo) error {
	fr, err := os.Open(path)
	if err != nil {
		return err
	}

	defer fr.Close()

	h := new(tar.Header)
	h.Name = fi.Name()
	h.Size = fi.Size()
	h.Mode = int64(fi.Mode())
	h.ModTime = fi.ModTime()

	err = tw.WriteHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(tw, fr)
	if err != nil {
		return err
	}

	return err
}
