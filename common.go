package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func readStdinOriginal() string {
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	return inputText
}

func readStdin() string {
	reader := bufio.NewReader(os.Stdin)
	inputText, _ := reader.ReadString('\n')

	return strings.TrimSpace(inputText)
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
