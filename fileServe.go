package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type folderData struct {
	Path  []string
	Dirs  []dirInfo
	Files []fileInfo
}

type fileInfo struct {
	Name         string
	Size         int64
	CreationDate time.Time
	Type         string
}

type dirInfo struct {
	Name         string
	Size         int64
	CreationDate time.Time
}

var rootFiles = "/home/dartie/Downloads" // TODO: load from settings.json file.
var rootPath_ = "list"                   // TODO: load from settings.json file.
var tmpPath = "tmp"                      // TODO: load from settings.json file.

func getPathsFromUrl(c *gin.Context, rootPath string) (string, []string) {
	fullUrl := c.Request.URL.Path // TODO : check html/url escape https://pkg.go.dev/github.com/henrylee2cn/gin/template
	urlParameter := strings.TrimPrefix(fullUrl, rootPath)
	serverRelativePathList := strings.Split(urlParameter, "/")
	serverRelativePathList = removeEmptyStrings(serverRelativePathList)
	rootFilesList := strings.Split(rootFiles, string(os.PathSeparator))
	serverPathList := append(rootFilesList, serverRelativePathList...)
	serverPath := filepath.Join(serverPathList...)

	// TODO: relative paths with dots should not be allowed

	// TODO: test on Windows
	if os.PathListSeparator == 58 {
		if len(serverPathList) > 0 {
			if serverPathList[0] == "" {
				serverPath = "/" + serverPath
			}
		}
	}

	return serverPath, serverRelativePathList
}

func getRootPath(url string) string {
	var rootPath string
	urlSplit := strings.Split(url, "/")
	for _, path := range urlSplit {
		if strings.HasPrefix(path, "*") {
			break
		}
		if path == "" {
			continue
		}
		rootPath += "/" + path
	}

	return rootPath
}

func getBase64String(localFile string) string {
	// Open file on disk.
	f, _ := os.Open(localFile)

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := io.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	return encoded
}

func listFileHandler(c *gin.Context) {
	rootPath := getRootPath(c.FullPath()) // TODO: Use c.Params[0] instead?
	serverPath, serverRelativePathList := getPathsFromUrl(c, rootPath)

	if file, err := os.Stat(serverPath); errors.Is(err, os.ErrNotExist) {
		message := "Requested path does not exists"
		status := "2"
		c.HTML(http.StatusOK, "home.html", gin.H{"Feedback": map[string]string{message: status}, "Url": "/"}) // "User": userInfoMap,
		return
	} else {
		if !file.IsDir() {
			// it's a file, download it
			fileName := c.Param("filename")
			targetPath := filepath.Join(serverPath, fileName)
			//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
			if !strings.HasPrefix(filepath.Clean(targetPath), serverPath) {
				c.String(403, "Look like you attacking me")
				return
			}
			//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
			c.Header("Content-Description", "File Transfer")
			c.Header("Content-Transfer-Encoding", "binary")
			c.Header("Content-Disposition", "attachment; filename="+fileName)
			c.Header("Content-Type", "application/octet-stream")
			c.File(targetPath)
			return
		}
	}

	/* Create struct with files info */
	files, err_ := os.ReadDir(serverPath)
	checkErr(err_)
	_ = files

	//***
	/* 	fileOpen, err := os.Open(serverPath)
	   	checkErr(err)
	   	defer fileOpen.Close()
	   	files, _ := fileOpen.Readdir(0)
	   	SortFileNameAscend(files)
	*/
	//***

	var thisFolder folderData

	thisFolder.Path = serverRelativePathList
	for _, file := range files {
		if file.IsDir() {
			dirFullpath := filepath.Join(serverPath, file.Name())
			folderOsInfo, err := os.Stat(dirFullpath)
			checkErr(err)

			// Build file struct info
			var thisDirectory dirInfo

			thisDirectory.Name = folderOsInfo.Name()
			thisDirectory.CreationDate = folderOsInfo.ModTime()
			thisDirectory.Size = folderOsInfo.Size()

			thisFolder.Dirs = append(thisFolder.Dirs, thisDirectory)

		} else {
			fileFullpath := filepath.Join(serverPath, file.Name())
			fileOsInfo, err := os.Stat(fileFullpath)
			checkErr(err)

			// Build file struct info
			var thisFile fileInfo

			thisFile.Name = file.Name()
			thisFile.CreationDate = fileOsInfo.ModTime()
			thisFile.Size = fileOsInfo.Size()
			thisFile.Type = strings.Trim(strings.ToLower(filepath.Ext(thisFile.Name)), ".")

			thisFolder.Files = append(thisFolder.Files, thisFile)
		}
	}

	c.HTML(http.StatusOK, "filelist.html", gin.H{"thisFolder": thisFolder, "rootPath": rootPath})
}

func viewFileHandler(c *gin.Context) {
	rootPath := getRootPath(c.FullPath())
	serverPath, _ := getPathsFromUrl(c, rootPath)
	fileName := c.Param("filename")
	targetPath := filepath.Join(serverPath, fileName)
	//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), serverPath) {
		c.String(403, "Look like you attacking me")
		return
	}
	//This headers needed for some browsers (for example without this headers Chrome will download files as txt)
	fileType := detectFileType(targetPath)

	c.Header("Content-Disposition", "filename="+fileName)
	c.Header("Content-Type", fileType)
	c.File(targetPath)
}

func downloadFolderHandler(c *gin.Context) {
	rootPath := getRootPath(c.FullPath())
	serverPath, _ := getPathsFromUrl(c, rootPath)
	fileName := c.Param("filename")
	targetPath := filepath.Join(serverPath, fileName)
	//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), serverPath) {
		c.String(403, "Look like you attacking me")
		return
	}

	cwd, err := os.Getwd()
	checkErr(err)
	zipDst := filepath.Join(cwd, filepath.Clean(tmpPath), targetPath+".zip")

	if _, err := os.Stat(filepath.Dir(zipDst)); errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(filepath.Dir(zipDst), os.ModePerm)
	}

	// create archive to download
	zipError := zipFolder(targetPath, zipDst)
	if zipError != nil {
		// TODO: to handle
		checkErr(zipError)
		c.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	//This headers needed for some browsers (for example without this headers Chrome will download files as txt)
	fileType := detectFileType(zipDst)

	c.Header("Content-Disposition", "filename="+fileName)
	c.Header("Content-Type", fileType)
	c.File(zipDst)

	// remove file
	e := os.Remove(zipDst)
	checkErr(e)
}

func viewFileHandler_(c *gin.Context) {
	rootPath := getRootPath(c.FullPath())
	serverPath, _ := getPathsFromUrl(c, rootPath)
	fileName := c.Param("filename")
	targetPath := filepath.Join(serverPath, fileName)
	//This ckeck is for example, I not sure is it can prevent all possible filename attacks - will be much better if real filename will not come from user side. I not even tryed this code
	if !strings.HasPrefix(filepath.Clean(targetPath), serverPath) {
		c.String(403, "Look like you are attacking me")
		return
	}
	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	/* 	c.Header("Content-Disposition", "filename="+fileName)
	   	c.Header("Content-Type", "application/octet-stream") */
	c.Header("Content-Disposition", "filename=a.png")
	c.Header("Content-Type", "image/png")
	c.File(targetPath)
	return
	fileBytes, err := os.ReadFile(targetPath)
	if err != nil {
		panic(err)
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(fileBytes)
}

func getUserAllowance(userId float64) map[string]string {
	accessQuery := `SELECT path, type FROM "Access" WHERE user_id = $1`
	rows, err := db.Query(accessQuery, userId)
	checkErr(err)
	defer rows.Close()

	pathList := make(map[string]string)
	for rows.Next() {
		var path, accessType string
		err := rows.Scan(&path, &accessType)
		checkErr(err)
		//pathList = append(pathList, path)
		pathList[path] = accessType
	}
	errRows := rows.Err()
	checkErr(errRows)

	return pathList
}

// middleware for checking whether the user is admin
func canViewTheContent(c *gin.Context) {
	userMap := getCurrentUserMap(c)
	userId := userMap["id"].(float64)

	pathList := getUserAllowance(userId)

	// Get folder path requested
	requestedPath := c.Params[0].Value

	_ = pathList
	_ = requestedPath

	if val, ok := pathList[requestedPath]; !ok {
		_ = val
		c.HTML(http.StatusForbidden, "403-Forbidden.html", nil)
		c.Abort()
		return
	}
	c.Next()

}

// DONE: file size for files
// DONE: view/download for files
// DONE: different icon according to the file type
// DONE: view list/grid
// DONE: table header sticky css
// DONE: download folder content as zip
// TODO: add table sorting
