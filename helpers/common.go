package helpers

import (
	"io/fs"
	"net"
	"os"
	"path/filepath"
)

var chPause = make(chan bool, 1)

// Pause the goroutine forever
func Pause() {
	<-chPause
}

// Resume process
func Resume() {
	if chPause != nil {
		chPause <- true
	}
}

// AbsolutePath returns absolute path
func AbsolutePath(base ...string) string {
	if base == nil {
		base = make([]string, 0)
	}
	if len(base) == 0 {
		base = append(base, "./")
	}

	p, err := filepath.Abs(base[0])
	if err != nil {
		return base[0]
	}

	return p
}

// EnsureDir create directories if not exists
func EnsureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func GlobFiles(path string, fn func(string) bool) []string {
	var files []string
	_ = filepath.WalkDir(path, func(s string, d fs.DirEntry, e error) error {
		if fn(s) {
			files = append(files, s)
		}
		return nil
	})
	return files
}

// ReadFile read file form root directory.
func ReadFile(path string) (content []byte, err error) {
	content, err = os.ReadFile(path)
	if err != nil {
		return
	}

	return
}

// WriteFile write file
func WriteFile(path string, content []byte) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// GetOutboundIP get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ipAddress := net.ParseIP(localAddr.String())
	if ipAddress.IsPrivate() {
		return "127.0.0.1"
	}

	return localAddr.String()
}
