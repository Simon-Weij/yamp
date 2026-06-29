package main

import (
	"net/http"
	"os"
	"strings"
)

type FileLoader struct {
	http.Handler
	allowedDirectories []string
}

func NewFileLoader(allowedDirectories []string) *FileLoader {
	//nolint:exhaustruct
	return &FileLoader{
		allowedDirectories: allowedDirectories,
	}
}

func (h *FileLoader) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		requestedFilename := req.URL.Path

		isAllowed := false
		for _, dir := range h.allowedDirectories {
			if strings.HasPrefix(requestedFilename, dir) {
				isAllowed = true
				break
			}
		}

		if isAllowed {
			fileData, err := os.ReadFile(requestedFilename)
			if err == nil {
				_, _ = res.Write(fileData)
				return
			}
		}
		next.ServeHTTP(res, req)
	})
}
