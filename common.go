package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/gin-gonic/gin"
	"github.com/gookit/color"
)

/* read settings functions */
func removeJsonComments(jsonStr string) string {
	// remove comments from json file
	var jsonStrNoComments string
	for _, settingLine := range strings.Split(jsonStr, "\n") {
		if !strings.HasPrefix(strings.TrimSpace(settingLine), "//") {
			jsonStrNoComments += settingLine + "\n"
		}
	}

	return jsonStrNoComments
}

func readSettingsInterface(settingFile string) map[string]interface{} {
	// read file
	settingsBytes, err1 := os.ReadFile(settingFile)
	if err1 != nil {
		log.Fatal(err1)
	}

	settingsStr := string(settingsBytes)

	settingsStrNoComments := removeJsonComments(settingsStr)

	rawData := []byte(settingsStrNoComments)
	var payload interface{}                  //The interface where we will save the converted JSON data.
	err := json.Unmarshal(rawData, &payload) // Convert JSON data into interface{} type
	if err != nil {
		log.Fatal(err)
	}

	m := payload.(map[string]interface{}) // To use the converted data we will need to convert it into a map[string]interface{}

	return m
}

func readSettings(settingsFile string) map[string]string {
	/* Read settings */
	settingsMap = make(map[string]string)

	if _, err := os.Stat(settingsFile); err == nil {
		settingsBytes, err := os.ReadFile(settingsFile)
		if err != nil {
			log.Fatal(err)
		}
		settingsStr := string(settingsBytes)

		// remove json comments
		settingsStrNoComments := removeJsonComments(settingsStr)

		json.Unmarshal([]byte(settingsStrNoComments), &settingsMap)

	} else {
		log.Fatalf(color.Red.Sprint(settingsFile + " does not exist!"))
	}

	return settingsMap
}

// Converts an interface slice of string to float
func interfaceStringToFloat(inputSlice interface{}) []float64 {
	var outputSlice []float64
	if inputSlice == nil {
		return outputSlice
	}
	for _, is := range inputSlice.([]interface{}) {
		iss := is.(string)
		floatNum, err := strconv.ParseFloat(iss, 64)
		if err == nil {
			outputSlice = append(outputSlice, floatNum)
		}
	}

	return outputSlice
}

// json bind
func BindWihtoutBindingTag(c *gin.Context, dest interface{}) error {
	if c.ContentType() != "application/json" {
		return errors.New("'BindWithoutBindingTag' only serves for application/json not for " + c.ContentType())
	}
	buf, e := io.ReadAll(c.Request.Body)
	if e != nil {
		return e
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(buf))
	a := json.Unmarshal(buf, &dest)
	fmt.Printf("%s", a)

	return a
}

/* General */

// Get current working directory
func getcwd() string {
	path, _ := os.Getwd()

	return path
}

// Get current source code file path
func getSrcPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	return filepath.Dir(filename)
}

// Get executable path
func getExePath() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	return filepath.Dir(ex)
}

// Read standard input text
// :return: (string) text entered
func readStdinOriginal() string {
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	return inputText
}

// Check if regex is matched
func reMatched(s string, regex string) bool {
	r, _ := regexp.Compile(regex)
	matches := r.FindAllStringSubmatch(s, -1)

	return len(matches) > 0
}

// Mac address functions
func validateMac(mac string) bool {
	r, _ := regexp.Compile("^[0-9A-Fa-f]{12}$")
	matches := r.FindAllStringSubmatch(mac, -1)

	return len(matches) > 0
}

/* Strings */
// Clean strings, remove non-ascii chars from string
func cleanString(s string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	t := re.ReplaceAllLiteralString(s, "")

	return t
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

// contains checks if a string is present in a interface
func containsInterface(s interface{}, str string) bool {
	a := 1
	_ = a
	sSlice := s.([]interface{})
	for _, v := range sSlice {
		if v == str {
			return true
		}
	}

	return false
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

// Remove element from slice
func RemoveIndexSlice(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// Remove element from slice
func RemoveIndexSliceNested(s [][]string, index int) [][]string {
	return append(s[:index], s[index+1:]...)
}

// Insert int value in Slice
func insertIntInNestedSlice(a [][]int, index int, value []int) [][]int {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// Insert string value in Slice
func insertStringInNestedSlice(a [][]string, index int, value []string) [][]string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// Insert int value in Slice
func insertIntInSlice(a []int, index int, value int) []int {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

// Insert string value in Slice
func insertStringInSlice(a []string, index int, value string) []string {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
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
