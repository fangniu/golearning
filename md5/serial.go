package main

import (
	"os"
	"fmt"
	"sort"
	"crypto/md5"
	"path/filepath"
	"io/ioutil"
)

// MD5All 读取 root 目录下的所有文件，返回一个map
// 该 map 存储了 文件路径到文件内容 md5值的映射
// 如果 Walk 执行失败，或者 ioutil.ReadFile 读取失败，
// MD5All 都会返回错误
func MD5All(root string) (map[string][md5.Size]byte, error) {
	fmt.Println(md5.Size)
	m := make(map[string][md5.Size]byte)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		m[path] = md5.Sum(data)
		return nil
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(m)
	return m, nil
}

func main() {
	// 计算特定目录下所有文件的 md5值，
	// 然后按照路径名顺序打印结果
	m, err := MD5All(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x  %s\n", m[path], path)
	}
}