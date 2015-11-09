package vault

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*func RPack(root string, dest io.Writer) error {

  absRoot, err := filepath.Abs(root)

  if err != nil {
    return err
  }

  var files []string
  var totalSize int64
  e := filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {

    if info.IsDir() {
      return nil
    }
    totalSize += info.Size()
    files = append(files, strings.Replace(path, absRoot, "", 1))
    return nil
  })

  w := zip.NewWriter(zipFile)
  defer w.Close()

  for _, file := range files {
    if file[0] == '/' {
      file = file[1:]
    }
    fullPath := filepath.Join(absRoot, file)

    wr, _ := w.Create(file)

    fw, _ := os.Open(fullPath)
    io.Copy(wr, fw)

    fw.Close()
  }

  return e

}*/

func Pack(root string, name string) error {

	absRoot, err := filepath.Abs(root)

	if err != nil {
		return err
	}

	var files []string
	var totalSize int64
	e := filepath.Walk(absRoot, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}
		totalSize += info.Size()
		files = append(files, strings.Replace(path, absRoot, "", 1))
		return nil
	})

	zipFile, _ := os.Create(name)
	defer zipFile.Close()
	w := zip.NewWriter(zipFile)
	defer w.Close()

	for _, file := range files {
		if file[0] == '/' {
			file = file[1:]
		}
		fullPath := filepath.Join(absRoot, file)

		wr, _ := w.Create(file)

		fw, _ := os.Open(fullPath)
		io.Copy(wr, fw)

		fw.Close()
	}

	return e
}

func Unpack(zipfile string, dest string) error {

	reader, err := zip.OpenReader(zipfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer reader.Close()

	for _, f := range reader.Reader.File {

		zipped, err := f.Open()
		if err != nil {
			return err
		}

		defer zipped.Close()

		// get the individual file name and extract the current directory
		path := filepath.Join(dest, f.Name)

		dirname := filepath.Dir(path)

		if !Exists(dirname) {
			os.MkdirAll(dirname, 0700)
		}

		if f.FileInfo().IsDir() {

			err := os.MkdirAll(path, f.Mode())
			if err != nil {
				return err
			}
		} else {
			writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, f.Mode())

			if err != nil {
				return err
			}

			defer writer.Close()

			if _, err = io.Copy(writer, zipped); err != nil {

				return err
			}

		}
	}
	return nil
}
