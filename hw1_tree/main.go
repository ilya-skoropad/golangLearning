package main

import (
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

const (
	file_print_pattern       = `%s (%sb)`
	empty_file_print_pattern = `%s (empty)`
	folder_print_pattern     = `%s`
)

const (
	empty_arrow    = "\t"
	straight_arrow = "│\t"
	t_arrow        = "├───"
	l_arrow        = "└───"
)

func dirTree(out io.Writer, path string, printFiles bool) error {
	state := make([]bool, 10, 255)
	printTree(out, path, 1, state, printFiles)

	return nil
}

func printTree(out io.Writer, path string, lvl int, state []bool, printFiles bool) {
	msg := ``
	i := 1
	last := false
	files, err := ioutil.ReadDir(path)

	if err != nil {
		log.Fatal(err)
	}

	SortFiles(files)

	for _, f := range files {
		last = isLast(files, i, printFiles)
		state[lvl-1] = !last
		msg = printArrow(state, lvl, last)

		if f.IsDir() {
			msg += printFolderName(f)
			fmt.Fprintln(out, msg)
			printTree(out, path+`/`+f.Name(), lvl+1, state, printFiles)
		} else if printFiles {
			msg += printFileName(f)
			fmt.Fprintln(out, msg)
		}
		i++
	}
}

func SortFiles(files []os.FileInfo) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})
}

func isLast(files []fs.FileInfo, current int, printFiles bool) bool {
	last := len(files)
	if !printFiles {
		last = getLastFolder(files)
	}

	if current == last {
		return true
	}

	return false
}

func getLastFolder(files []fs.FileInfo) int {
	i := 1
	last := 0

	for _, f := range files {
		if f.IsDir() {
			last = i
		}
		i++
	}

	return last
}

func printFileName(file os.FileInfo) string {
	msg := fmt.Sprintf(file_print_pattern, file.Name(), fmt.Sprint(file.Size()))
	if file.Size() == 0 {
		msg = fmt.Sprintf(empty_file_print_pattern, file.Name())
	}

	return msg
}

func printFolderName(folder os.FileInfo) string {
	msg := folder.Name()

	return msg
}

func printArrow(state []bool, level int, last bool) string {
	msg := printPreviewArrows(state, level)

	add := t_arrow
	if last {
		add = l_arrow
	}

	return msg + add
}

func printPreviewArrows(state []bool, level int) string {
	msg := ``
	for i := 0; i < level-1; i++ {
		add := empty_arrow
		if state[i] {
			add = straight_arrow
		}

		msg += add
	}

	return msg
}
