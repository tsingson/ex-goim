package utils

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/karrick/godirwalk"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const bufferSize = 65536

func Exists(name string) bool {
	afs := afero.NewOsFs()
	b, e := afero.Exists(afs, name)
	if e != nil {
		return false
	}
	return b
}

func DirExists(name string) bool {
	afs := afero.NewOsFs()
	b, e := afero.DirExists(afs, name)
	if e != nil {
		return false
	}
	return b
}

func WriteToFile(c []byte, filename string) error {
	// 将指定内容写入到文件中
	err := ioutil.WriteFile(filename, c, 0666)
	return err
}

// 获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListFiles(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) // 忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { // 匹配文件
			files = append(files, dirPth+PthSep+fi.Name())
		}
	}
	return files, nil
}

func ListDir(dirPth string) (files []string, err error) {
	files = make([]string, 0, 100)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			files = append(files, StrBuilder(dirPth, PthSep, fi.Name()))
		}
	}
	// litter.Dump(files)
	return files, nil
}

func Md5CheckSum(filename string) (string, error) {
	if info, err := os.Stat(filename); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", nil
	}

	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	for buf, reader := make([]byte, bufferSize), bufio.NewReader(file); ; {
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		hash.Write(buf[:n])
	}

	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	return checksum, nil
}
func ListSubPath(osDirname string) ([]string, error) {

	children, err := godirwalk.ReadDirnames(osDirname, nil)
	if err != nil {
		err1 := errors.Wrap(err, "cannot get list of directory children")
		return nil, err1
	}
	sort.Strings(children)
	var sublist []string
	sublist = make([]string, len(children))
	for _, child := range children {
		pathNode := StrBuilder(osDirname, "/", child, "/")
		// 	fmt.Printf("%s\n", pathNode)
		sublist = append(sublist, pathNode)
	}

	return sublist, nil
}

func GetCurrentPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func GetCurrentExecDir() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		// 	fmt.Printf("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		// 	fmt.Printf("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}

// B2S converts a byte slice to a string.
// It's fasthttpgx, but not safe. Use it only if you know what you're doing.
func B2S(b []byte) string {
	bytesHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	strHeader := reflect.StringHeader{
		Data: bytesHeader.Data,
		Len:  bytesHeader.Len,
	}
	return *(*string)(unsafe.Pointer(&strHeader))
}

// S2B converts a string to a byte slice.
// It's fasthttpgx, but not safe. Use it only if you know what you're doing.
func S2B(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bytesHeader := reflect.SliceHeader{
		Data: strHeader.Data,
		Len:  strHeader.Len,
		Cap:  strHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bytesHeader))
}

func StrBuilder(args ...string) string {
	var str strings.Builder

	for _, v := range args {
		str.WriteString(v)
	}
	return str.String()
}
