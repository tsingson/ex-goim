package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

func SaveToml(v interface{}, filename string) error {
	// currentPath+"/toml1.toml"
	var err error
	b := &bytes.Buffer{}
	encoder := toml.NewEncoder(b)
	if err = encoder.Encode(v); err != nil {

	}
	WriteToFile(b.Bytes(), filename)
	return err
}

func WriteToFile(c []byte, filename string) error {
	// 将指定内容写入到文件中
	err := ioutil.WriteFile(filename, c, 0666)
	return err
}

func GetCurrentPath() (string, error) {
	return filepath.Abs(filepath.Dir(os.Args[0]))
}

func GetCurrentExecDir() (dir string, err error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		fmt.Printf("exec.LookPath(%s), err: %s\n", os.Args[0], err)
		return "", err
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("filepath.Abs(%s), err: %s\n", path, err)
		return "", err
	}
	dir = filepath.Dir(absPath)
	return dir, nil
}
