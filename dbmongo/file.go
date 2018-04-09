package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// File attributs de fichier
type File struct {
	Path string    `json:"path"`
	Name string    `json:"name"`
	Date time.Time `json:"date"`
}

func listFiles(c *gin.Context) {
	fileList := []File{}
	root := "./data-raw/"
	filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		file := File{
			Path: path,
			Name: f.Name(),
			Date: f.ModTime(),
		}
		fileList = append(fileList, file)
		return nil
	})
	c.JSON(200, fileList)
}
