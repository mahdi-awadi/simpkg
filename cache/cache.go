package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-per/simpkg/helpers"
	"github.com/go-per/simpkg/std"
)

// Instance is cache instance
var Instance *Cache

// Cache struct
type Cache struct {
	root       string
	extension  string
	filePrefix string
}

// ICache interface
type ICache interface {
	SetRoot(root string)
	GetRoot() string
	SetExtension(extension string)
	SetFilePrefix(prefix string)
	Path(cacheName string, addTimestamp ...bool) (string, string)
	Write(cacheName string, data []byte, addTimestamp ...bool) error
	WriteAsync(cacheName string, data []byte, addTimestamp ...bool)
	Get(cacheName string) ([]byte, error)
	Remove(cacheName string) error
	Clear() error
}

// create default instance
func init() {
	Instance = New()
}

// New returns new cache
func New() *Cache {
	return &Cache{
		root:      "./cache",
		extension: "txt",
	}
}

// SetRoot sets cache root
func (cache *Cache) SetRoot(root string) {
	cache.root = root
}

// GetRoot returns cache root
func (cache *Cache) GetRoot() string {
	return cache.root
}

// SetExtension sets cache extension
func (cache *Cache) SetExtension(extension string) {
	cache.extension = strings.TrimLeft(extension, ".")
}

// SetFilePrefix sets cache file prefix
func (cache *Cache) SetFilePrefix(prefix string) {
	cache.filePrefix = prefix
}

// Path return cache path
func (cache *Cache) Path(cacheName string, addTimestamp ...bool) (string, string) {
	timeStamp := ""
	if addTimestamp != nil && len(addTimestamp) > 0 && addTimestamp[0] {
		timeStamp += fmt.Sprintf("-%v", time.Now().UnixMilli())
	}

	var cachePath []string
	cachePath = append(cachePath, cache.root)

	cacheNameParts := strings.Split(cacheName, filepath.FromSlash("/"))
	partsLen := len(cacheNameParts) - 1
	cacheFileName := cacheNameParts[partsLen]
	if len(cacheNameParts) > 1 {
		for i, part := range cacheNameParts {
			if i < partsLen {
				cachePath = append(cachePath, part)
			}
		}
	}

	extension := cache.extension
	basePath := filepath.Join(cachePath...)
	prefix := ""
	if cache.filePrefix != "" {
		prefix = cache.filePrefix + "-"
	}
	if filepath.Ext(cacheFileName) != "" {
		extension = filepath.Ext(cacheFileName)
		cacheFileName = strings.TrimSuffix(cacheFileName, extension)
	}

	return basePath, filepath.Join(basePath, fmt.Sprintf("%s%s%s.%v", prefix, cacheFileName, timeStamp, strings.Replace(extension, ".", "", -1)))
}

// WriteAsync writes cache file asynchronously
func (cache *Cache) WriteAsync(cacheName string, data []byte, addTimestamp ...bool) {
	go func() {
		err := cache.Write(cacheName, data, addTimestamp...)
		if err != nil {
			std.Error("Could not write cache file %v", err)
		}
	}()
}

// Write writes cache file
func (cache *Cache) Write(cacheName string, data []byte, addTimestamp ...bool) (err error) {
	dir, file := cache.Path(cacheName, addTimestamp...)
	_ = helpers.EnsureDir(dir)
	err = helpers.WriteFile(file, data)
	return
}

// Get return cached version of data
func (cache *Cache) Get(cacheName string) ([]byte, error) {
	_, file := cache.Path(cacheName)
	return helpers.ReadFile(file)
}

// Remove removes cache file if available
func (cache *Cache) Remove(cacheName string) error {
	_, file := cache.Path(cacheName)
	return os.Remove(file)
}

// Clear clears cache
func (cache *Cache) Clear() error {
	contents, err := filepath.Glob(filepath.Join(cache.root, "*."+cache.extension))
	if err != nil {
		return err
	}
	for _, item := range contents {
		err = os.RemoveAll(item)
		if err != nil {
			return err
		}
	}
	return nil
}
