package util

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Zip(dst, src string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 w
		_, err = io.Copy(w, fr)
		if err != nil {
			return
		}
		// 输出压缩的内容
		//fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetAllFile(pathname string) (files []string, err error) {
	err = filepath.Walk(pathname, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	return
}

func CreateFileIfAbsent(fileName string) (f *os.File, err error) {

	filePath := fileName[:strings.LastIndex(fileName, "/")]

	err = os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	exists, err := PathExists(fileName)
	if err != nil {
		return nil, err
	}

	if exists {
		f, err = os.Open(fileName)
		if err != nil {
			return nil, err
		}
	} else {
		f, err = os.Create(fileName)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

func CreateZip(zipPath, zipFile string) ([]byte, error) {
	zipPaths := strings.Split(zipPath, "/")
	err := Zip(zipFile, zipPaths[0])
	if err != nil {
		return nil, fmt.Errorf("压缩zip包失败,error=%s", err.Error())
	}
	zip, err := os.Open(zipFile)
	if err != nil {
		return nil, err
	}

	fileInfo, err := zip.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = zip.Read(buffer)
	if err != nil {
		return nil, err
	}
	defer zip.Close()
	return buffer, nil
}

func RemovePath(zipPath, zipFile string) error {
	zipPaths := strings.Split(zipPath, "/")
	err := os.RemoveAll(zipPaths[0])
	if err != nil {
		return fmt.Errorf("删除文件夹失败,error=%s", err.Error())
	}
	if zipFile != "" {
		err = os.Remove(zipFile)
		if err != nil {
			return fmt.Errorf("删除zip包失败,error=%s", err.Error())
		}
	}
	return nil
}
