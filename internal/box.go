//go:generate go run generator.go

package box

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type embedBox struct {
	storage map[string][]byte
}

// Create new box for embed files
func newEmbedBox() *embedBox {
	return &embedBox{storage: make(map[string][]byte)}
}

// Add a file to box
func (e *embedBox) Add(file string, content []byte) {
	e.storage[file] = content
}

// Get file's content
func (e *embedBox) Get(file string) []byte {
	if f, ok := e.storage[file]; ok {
		return f
	}
	return nil
}

// Embed box expose
var box = newEmbedBox()

// Add a file content to box
func Add(file string, content []byte) {
	box.Add(file, content)
}

// Get a file from box
func Get(file string) []byte {
	return box.Get(file)
}

// FileSystem interface
// inspiration from https://dev.to/koddr/the-easiest-way-to-embed-static-files-into-a-binary-file-in-your-golang-app-no-external-dependencies-43pc
type Dir string

type FileSystem interface {
	Open(name string) (File, error)
}

type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

func (d Dir) Open(name string) (File, error) {
	return nil, errors.New("open not implemented")
}

type fileHandler struct {
	root FileSystem
}

func FileServer(root FileSystem) http.Handler {
	return &fileHandler{root}
}

func (f *fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}
	ServeFile(w, r, f.root, path.Clean(upath), true)
}

var rootDir = "./static"

func GetFileBytes(filename string) []byte {
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		log.Infof("could not find dir %s", rootDir)
		return Get(filename)
	}

	log.Infof("serve file from disk %s", rootDir+filename)
	fileContent, err := ioutil.ReadFile(rootDir + filename)
	if err != nil {
		log.WithError(err).Fatalf("could not read file %s", rootDir+filename)
	}
	return fileContent
}

func ServeFile(w io.Writer, r *http.Request, fs FileSystem, name string, redirect bool) {
	log.Infof("serve file %s ", name)
	if name == "/" {
		_, err := w.Write(GetFileBytes("/index.html"))
		if err != nil {
			log.WithError(err).Fatalf("failed to serve file %s", name)
		}
	} else {
		_, err := w.Write(GetFileBytes(name))
		if err != nil {
			log.WithError(err).Fatalf("failed to serve file %s", name)
		}
	}
}
