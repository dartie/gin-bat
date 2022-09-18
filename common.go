package main

import (
	"archive/zip"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
	"unsafe"
)

func readStdinOriginal() string {
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	return inputText
}

// reads input from user
func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	return strings.TrimSpace(inputText)
}

// generates random string (most efficient: https://stackoverflow.com/a/31832326/4768254)
func RandStringBytesMaskImprSrcUnsafe(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// removeEmptyStrings - Use this to remove empty string values inside an array.
// This happens when allocation is bigger and empty
func removeEmptyStrings(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// detect the file type and return mimeType
func detectFileType(filepath string) string {
	file, err := os.Open(filepath)
	checkErr(err)
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	mimeType := http.DetectContentType(bytes)

	return mimeType
}

// Displays file size in human readable format "ByteCountDecimal"
func ByteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}

// Displays file size in human readable format "ByteCountBinary"
func ByteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

// this is the default sort order of golang os.Open -> ReadDir
func SortFileNameAscend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

func SortFileNameDescend(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() > files[j].Name()
	})
}

// compress a folder
func zipFolder(source, target string) error {
	// 1. Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// 2. Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 3. Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// 4. Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// 5. Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func map2(data []string, f func(string) string) []string {
	mapped := make([]string, len(data))

	for i, e := range data {
		mapped[i] = f(e)
	}

	return mapped
}

func getQueryFields(query string) []string {
	/* Get fields */
	var fields []string
	r, _ := regexp.Compile(`(?:SELECT\s+)(.*)\s+(?:FROM)`)
	fieldsMatch := r.FindAllStringSubmatch(query, -1)
	if len(fieldsMatch) > 0 {
		fields = strings.Split(fieldsMatch[0][1], ",")
		fields = map2(fields, strings.TrimSpace)
	}

	return fields
}

func sqlQuery(query string, args ...any) *sql.Rows {
	var result []string
	_ = result

	rows, err := db.Query(query, args...)

	//defer rows.Close()

	err = rows.Err()
	if err != nil {
		log.Print(err)
	}
	return rows
}
