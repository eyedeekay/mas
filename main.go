package main

import (
	"fmt"
	"github.com/otiai10/copy"
	"gitlab.com/golang-commonmark/markdown"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func checkcssdir(dir string) string {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 4 {
			if path[len(path)-4:] == ".css" {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	ret := ""
	for _, style := range files {
		ret += `<link rel="stylesheet" href="` + style + `">
`
	}
	return ret
}

func checkjsdir(dir string) string {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 {
			if path[len(path)-3:] == ".js" {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	ret := ""
	for _, script := range files {
		ret += `<script src="` + script + `"></script>
`
	}
	return ret
}

var bottom = `  </body>
</html>`

func main() {
	var files []string
	md := markdown.New(markdown.XHTMLOutput(true), markdown.HTML(true))

	root := "./"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if len(path) > 3 {
			if path[len(path)-3:] == ".md" {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	builddir := filepath.Join(root, "site/")
	cssdir := filepath.Join(root, "css/")
	jsdir := filepath.Join(root, "js/")
	imgdir := filepath.Join(root, "images/")
	cssbdir := filepath.Join(builddir, "css/")
	jsbdir := filepath.Join(builddir, "js/")
	imgbdir := filepath.Join(builddir, "images/")
	os.MkdirAll(builddir, 0755)
	if file, err := os.Stat(cssdir); !os.IsNotExist(err) {
		if file.IsDir() {
			copy.Copy(cssdir, cssbdir)
		}
	}
	if file, err := os.Stat(jsdir); !os.IsNotExist(err) {
		if file.IsDir() {
			copy.Copy(jsdir, jsbdir)
		}
	}
	if file, err := os.Stat(imgdir); !os.IsNotExist(err) {
		if file.IsDir() {
			copy.Copy(imgdir, imgbdir)
		}
	}
	var top = `<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title>title</title>
      ` + checkcssdir(cssbdir) + `
      ` + checkjsdir(jsbdir) + `
    </head>
    <body>
`
	for _, file := range files {
		if strings.Contains(file, "/") {
			if err := os.MkdirAll(filepath.Join(builddir, filepath.Dir(file)), 0755); err == nil {
				if sitefile, err := os.Create(filepath.Join(builddir, filepath.Dir(file), filepath.Base(file)+".html")); err == nil {
					if bytes, err := ioutil.ReadFile(file); err == nil {
						fmt.Println(md.RenderToString(bytes))
						sitefile.Write([]byte(md.RenderToString(bytes)))
					}
					sitefile.Close()
				}
			}
		} else {
			if sitefile, err := os.Create(filepath.Join(builddir, file+".html")); err == nil {
				if bytes, err := ioutil.ReadFile(file); err == nil {
					fmt.Println(md.RenderToString(bytes))
					sitefile.Write([]byte(md.RenderToString(bytes)))
				}
				sitefile.Close()
			}
		}
	}
	var dirs []string
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if !strings.Contains(path, ".git") {
				dirs = append(dirs, path)
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, dir := range dirs {
		tmpfile, err := os.Create(filepath.Join(dir, "index.html"))
		tmpfile.Write([]byte(top))
		count := 0
		if err == nil {
			if files, err := ioutil.ReadDir(dir); err == nil {
				for _, file := range files {
					if len(file.Name()) > 5 {
						if file.Name()[len(file.Name())-5:] == ".html" {
							if file.Name() != "index.html" {
								if bytes, err := ioutil.ReadFile(filepath.Join(dir, file.Name())); err == nil {
									//fmt.Println(md.RenderToString(bytes))
									tmpfile.Write([]byte(md.RenderToString(bytes)))
								}
								count++
							}
						}
					}
				}
			}
		}
		tmpfile.Write([]byte(bottom))
		tmpfile.Close()
		if count == 0 {
			os.Remove(filepath.Join(dir, "index.html"))
		}
	}
}
