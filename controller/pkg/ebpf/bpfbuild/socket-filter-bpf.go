// Code generated by go-bindata.
// sources:
// ../dist/socket-filter-bpf.o
// DO NOT EDIT!

package bpfbuild

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _socketFilterBpfO = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xe4\x94\xb1\x8f\x12\x41\x14\xc6\xbf\x61\x41\x16\xa4\xa0\x10\x83\xd1\x82\x92\xc2\x2c\x46\x8d\xb1\x32\x84\x44\x2a\x4c\x8c\xb1\x30\xb1\x20\x2b\x59\x83\x59\x11\xc2\x6c\xa1\xc4\xc4\xca\xc6\x8a\xc6\x86\xd2\xca\xff\xe0\xae\xdb\xf6\xfe\x0c\x8a\x2b\x2e\x57\x51\x5c\x72\x97\x6b\xe6\xb2\x33\x6f\xd8\x65\x76\xd9\xbb\xfe\xbe\x82\xb7\xef\x37\xfb\x76\xde\x7c\x79\xc3\xaf\xd7\x83\x7e\x81\x31\x68\x31\x9c\x23\xce\x62\x1d\x59\xf1\x73\x97\x7e\x2b\x60\x08\xef\x2b\xd6\x02\x50\x03\x10\x96\x55\x3e\x5a\xac\x85\xe6\xf5\x88\xdb\xc4\xff\x1c\x4b\xde\x06\x70\x2f\xe2\x15\xc5\xfd\xd5\xc9\x96\x47\x9f\xf4\xab\xa7\x42\xbd\xbf\x51\x71\x71\x26\xa3\x5f\xbd\x50\x71\x75\x29\x63\xf8\x4f\xd5\x97\x0b\xc0\x5a\x08\xd1\x34\x9a\xff\x2d\xcf\x04\x34\xe8\x54\x25\xa8\x0d\x93\x75\x9b\x9c\xba\x90\x78\x13\x80\x10\x42\xe8\xf5\x06\x79\x76\x40\xf9\xdf\xad\x7f\xca\x07\xb9\x5a\xcf\x30\xf2\x96\x4a\x5b\xc1\x16\xef\x60\xff\xbc\xcb\x1e\x12\x6b\xea\xf5\xac\xa1\x4b\xe8\xbd\xfc\xb5\x70\x68\xf0\x37\xc4\xcd\xa1\xed\x13\xb7\x33\xbe\x6b\xc1\x4a\x31\x35\xa7\x69\x5e\x93\xbc\x94\xe2\x4b\x7d\x1e\x00\x77\xa2\xfb\x61\xe4\x8f\x28\xaf\x02\x28\x46\x0f\x4e\xe0\x7d\x0f\x30\x71\x67\xbc\xc3\x3d\xce\xbf\x4c\xbf\x71\x38\x73\xef\x2b\x9f\x8e\x7c\x2f\xe8\xb8\xb3\xd9\xd0\x1d\xf9\x12\x39\xde\x78\xf8\x79\xee\x4e\x3c\x38\x3c\x98\x07\xee\x27\x38\xfc\xc7\x24\x8a\x83\x5e\xef\xc9\xf0\xb9\x0a\xcf\x54\x78\x9a\xef\xdb\x4d\xf5\x42\xb9\x98\xd2\x98\xe0\x47\x83\x9b\xb6\xb2\xc4\xd9\x93\xea\xee\xd9\xaf\x68\xe4\x0f\xae\xa9\x37\xe7\xc3\x36\xde\x6b\x00\x74\xb3\x77\xf5\x92\xfa\x6f\x25\xea\xac\x44\xbd\x9e\xcb\x32\xed\x6f\x7a\xf0\x56\xdf\x7f\x96\xdf\xff\x63\xaa\x2f\x18\x7c\x4c\xa0\x8d\xfc\xfe\xdb\x7b\xfa\xff\x60\xed\xf6\x69\x93\x47\x66\xff\xaf\x32\xf6\x8e\xb4\x24\xf8\x9f\x72\x26\xff\x03\xe3\x7a\x7d\xff\xae\x02\x00\x00\xff\xff\x2e\x02\xe6\xdc\x08\x06\x00\x00")

func socketFilterBpfOBytes() ([]byte, error) {
	return bindataRead(
		_socketFilterBpfO,
		"socket-filter-bpf.o",
	)
}

func socketFilterBpfO() (*asset, error) {
	bytes, err := socketFilterBpfOBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "socket-filter-bpf.o", size: 1544, mode: os.FileMode(420), modTime: time.Unix(1, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"socket-filter-bpf.o": socketFilterBpfO,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"socket-filter-bpf.o": &bintree{socketFilterBpfO, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

