// Package index creates pages for displaying directory listings.
package index

import (
	"path"
	"sort"
	"strings"
	"sync"

	"github.com/FooSoft/goldsmith"
)

type directory struct {
	entries   EntriesByName
	indexFile *goldsmith.File
}

// Entry contains information about a directory item.
type Entry struct {
	Name  string
	Path  string
	IsDir bool
	File  *goldsmith.File
}

// Index chainable plugin context.
type Index struct {
	indexName string
	filesKey  string
	indexMeta map[string]interface{}

	dirLists    map[string]*directory
	dirsHandled map[string]bool
	mutex       sync.Mutex
}

// New creates a new instance of the Index plugin.
// The meta parameter allows additional metadata to be provided for generated indices.
func New(meta map[string]interface{}) *Index {
	return &Index{
		indexName:   "index.html",
		indexMeta:   meta,
		filesKey:    "Files",
		dirsHandled: make(map[string]bool),
		dirLists:    make(map[string]*directory),
	}
}

// IndexFilename sets the name of the file to be created as the directory index (default: "index.html").
func (plugin *Index) IndexFilename(filename string) *Index {
	plugin.indexName = filename
	return plugin
}

// FilesKey sets the metadata key used to access the files in the current directory (default: "Files").
func (plugin *Index) FilesKey(key string) *Index {
	plugin.filesKey = key
	return plugin
}

func (*Index) Name() string {
	return "index"
}

func (plugin *Index) Process(context *goldsmith.Context, inputFile *goldsmith.File) error {
	plugin.mutex.Lock()
	defer plugin.mutex.Unlock()

	currentPath := inputFile.Path()
	currentIsDir := false

	for {
		if handled, _ := plugin.dirsHandled[currentPath]; handled {
			break
		}

		plugin.dirsHandled[currentPath] = true

		currentDir := path.Dir(currentPath)
		currentBase := path.Base(currentPath)

		list, ok := plugin.dirLists[currentDir]
		if !ok {
			list = new(directory)
			plugin.dirLists[currentDir] = list
		}

		if !currentIsDir {
			if currentBase == plugin.indexName {
				list.indexFile = inputFile
			} else {
				context.DispatchFile(inputFile)
			}
		}

		entry := Entry{Name: currentBase, Path: currentPath, IsDir: currentIsDir, File: inputFile}
		list.entries = append(list.entries, entry)

		if currentDir == "." {
			break
		}

		currentPath = currentDir
		currentIsDir = true
	}

	return nil
}

func (plugin *Index) Finalize(context *goldsmith.Context) error {
	for name, list := range plugin.dirLists {
		sort.Sort(list.entries)

		indexFile := list.indexFile
		if indexFile == nil {
			indexFile = context.CreateFileFromData(path.Join(name, plugin.indexName), make([]byte, 0))
			for name, value := range plugin.indexMeta {
				indexFile.Meta[name] = value
			}
		}

		indexFile.Meta[plugin.filesKey] = list.entries
		context.DispatchFile(indexFile)
	}

	return nil
}

type EntriesByName []Entry

func (d EntriesByName) Len() int {
	return len(d)
}

func (d EntriesByName) Less(i, j int) bool {
	d1, d2 := d[i], d[j]

	if d1.IsDir && !d2.IsDir {
		return true
	}
	if !d1.IsDir && d2.IsDir {
		return false
	}

	return strings.Compare(d1.Name, d2.Name) == -1
}

func (d EntriesByName) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
