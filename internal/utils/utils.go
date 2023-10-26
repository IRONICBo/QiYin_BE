package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"
)

// GenerateObjectName generate object name with folder and random filename.
func GenerateObjectName(filename string) string {
	// generate folder name.
	folderName := GenerateFolderName()

	// generate random file name.
	randomFileName := GenerateFileName(filename)

	return fmt.Sprintf("%s/%s", folderName, randomFileName)
}

// GenerateFolderName generate folder name.
func GenerateFolderName() string {
	// generate foler name with time string.
	return time.Now().Format("20060102")
}

// GenerateFileName generate file name.
func GenerateFileName(filename string) string {
	// generate random file name with file extension.

	// get file extension.
	ext := filepath.Ext(filename)

	// generate random file name.
	randomFileName := GenerateRandomString(64)
	// url encode.
	// randomFileName = url.QueryEscape(randomFileName)

	return randomFileName + ext
}

// GenerateRandomString generate random string.
func GenerateRandomString(length int) string {
	// dataset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	dataset := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()"
	rand.Seed(time.Now().UnixNano())

	str := make([]byte, length)
	for i := 0; i < length; i++ {
		str[i] = dataset[rand.Intn(len(dataset))]
	}

	return string(str)
}
