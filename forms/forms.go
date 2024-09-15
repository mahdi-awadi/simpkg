package forms

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-per/simpkg/parse"
)

// locker is a locke for map items
var locker = sync.RWMutex{}

type OnFormLoad func(name string, form *Form)

// Forms is a struct for request forms
type Forms struct {
	forms      map[string]*Form
	rootPath   string
	filesExt   string
	OnFormLoad OnFormLoad
}

// Instance request forms instance
var Instance = New()

// New is a constructor for Forms
func New() *Forms {
	return &Forms{
		forms:    make(map[string]*Form),
		rootPath: "./",
		filesExt: "json",
	}
}

// SetRootPath sets root path
func (rf *Forms) SetRootPath(rootPath string) {
	rf.rootPath = rootPath
}

// SetFilesExt sets files extension
func (rf *Forms) SetFilesExt(filesExt string) {
	rf.filesExt = strings.TrimPrefix(".", filesExt)
}

// GetRootPath returns forms directory
func (rf *Forms) GetRootPath() string {
	var pathSegments []string
	pathSegments = append(pathSegments, rf.rootPath)
	return filepath.Join(pathSegments...)
}

// Load loads forms from directory
func (rf *Forms) Load() error {
	var files []string

	rootPath := rf.GetRootPath()
	_ = filepath.WalkDir(rootPath, func(p string, d fs.DirEntry, e error) error {
		if e == nil && d != nil && !d.IsDir() && strings.HasSuffix(d.Name(), "."+rf.filesExt) {
			files = append(files, p)
		}
		return nil
	})

	if files == nil || len(files) < 1 {
		return errors.New("Could not load Forms or Forms not exists: " + rootPath)
	}

	for _, file := range files {
		fileName := strings.Replace(file, rootPath, "", 1)
		fileName = strings.TrimSuffix(fileName, path.Ext(fileName))
		fileName = strings.TrimPrefix(fileName, string(os.PathSeparator))
		fileName = strings.ReplaceAll(fileName, string(os.PathSeparator), ".")

		// load file
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var form *Form
		err = json.Unmarshal(content, &form)
		if err != nil {
			return err
		}

		// add form name
		rf.AddForm(fileName, form)
	}

	// ensure Forms loaded
	if len(rf.forms) == 0 {
		return errors.New("Could not load Forms or Forms not exists:" + rootPath)
	}

	return nil
}

// AddForm to forms list
func (rf *Forms) AddForm(name string, form *Form) {
	locker.RLock()
	form.Method = strings.ToUpper(form.Method)
	form.name = name
	rf.forms[name] = form
	if rf.OnFormLoad != nil {
		rf.OnFormLoad(name, form)
	}

	// is json string
	jsonString, err := parse.ToJsonString(form.Body)
	if err != nil {
		form.Error = err
	} else {
		form.bodyString = jsonString
	}

	locker.RUnlock()
}

// Get returns a form by name
func (rf *Forms) Get(key string) (*Form, bool) {
	locker.RLock()
	defer locker.RUnlock()
	f, o := rf.forms[key]
	return f, o
}

// GetForms returns forms list
func (rf *Forms) GetForms() map[string]*Form {
	return rf.forms
}
