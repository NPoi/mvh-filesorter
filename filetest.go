package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
)

func main() {

	var root string
	if len(os.Args) > 1 && os.Args[1] != "" {
		root = os.Args[1]
	} else {
		root = "/Volumes/UNTITLED/"
	}

	var file_list []string
	var dir_list []string

	err := filepath.Walk(root,
		func(rel_path string, info os.FileInfo, err error) error {
			fmt.Println(info.Name())

			if info.Name() == ".com.apple.timemachine.donotpresent" {
				return nil
			}

			if info.Name()[0] == '.' {
				return filepath.SkipDir
			} else if info.Name() == "temp_" {
				return filepath.SkipDir
			}

			if info.IsDir() {
				rel, err := filepath.Rel(root, rel_path)
				if err != nil {
					os.Exit(2)
				}
				dir_list = append(dir_list, rel)

				return nil
			}

			rel, err := filepath.Rel(root, rel_path)
			file_list = append(file_list, rel)

			return nil
		})

	if err != nil {
		panic(err)
	}

	fmt.Println(dir_list)
	fmt.Println(file_list)

	top_dirs, err := ioutil.ReadDir(root)

	if err != nil {
		panic(err)
	}

	sort.Strings(file_list)

	err = os.Mkdir(path.Join(root, "temp_"), 0755)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(top_dirs); i++ {
		// rel = path.Join(string(root), string(rel))
		if top_dirs[i].Name()[0] == '.' {
			continue
		} else if top_dirs[i].Name() == "temp_" {
			continue
		}
		os.Rename(path.Join(root, top_dirs[i].Name()), path.Join(root, "temp_", top_dirs[i].Name()))
	}

	for i := 0; i < len(dir_list); i++ {
		os.Mkdir(path.Join(root, dir_list[i]), 0755)
	}

	for i := 0; i < len(file_list); i++ {
		// rel = path.Join(string(root), string(rel))
		// err = os.Link(path.Join(root, "temp_", file_list[i]), path.Join(root, file_list[i]))
		fmt.Println(path.Join(root, file_list[i]))
		content, err := os.Open(path.Join(root, "temp_", file_list[i]))
		if err != nil {
			panic(err)
		}
		defer content.Close()

		dst, err := os.Create(path.Join(root, file_list[i]))
		if err != nil {
			panic(err)
		}
		defer dst.Close()

		err = io.WriteFile(path.Join(root, file_list[i]), content, 0644)
		if err != nil {
			panic(err)
		}
	}

	err = os.RemoveAll(path.Join(root, "temp_"))
	if err != nil {
		panic(err)
	}
}
