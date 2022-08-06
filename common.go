package main

import (
	"database/sql"
	"log"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getCheckBoxValue(c *gin.Context, field string) bool {
	fieldValue := c.PostForm(field)
	var fieldValueBool bool
	if fieldValue == "" {
		fieldValueBool = false
	} else {
		fieldValueBool = true
	}
	return fieldValueBool
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
