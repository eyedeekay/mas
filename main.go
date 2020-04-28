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
		ret += `
      <link rel="stylesheet" href="/` + strings.Replace(style, "site/", "", -1) + `">`
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
		ret += `
      <script src="/` + strings.Replace(script, "site/", "", -1) + `"></script>`
	}
	return ret
}

func findTitle(markdown string) string {
	return "Placeholder"
}

func top(title, jdir, cdir string) []byte {
	return []byte(`<!DOCTYPE html>
  <html lang="en">
    <head>
      <meta charset="utf-8">
      <title>` + title + `</title>` + checkcssdir(cdir) + `` + checkjsdir(jdir) + `
    </head>
    <body>
    <!-- NAV AREA -->
`)
}

func argCat() string {
	var args string
	if len(os.Args) < 2 {
		return "MAS is a simple site generator"
	}
	for _, arg := range os.Args[1:] {
		args += arg + " "
	}
	return strings.TrimSuffix(args, " ")
}

func deSuffix(name string) string {
	if len(strings.SplitN(name, ".", 2)) > 1 {
		r := strings.Replace(strings.Replace(name, strings.SplitN(name, ".", 2)[1], "", -1), ".", "", -1)
		fmt.Println("Stripping suffix from filename", name)
		return r
	}
	return ""
}

var bottom = `  </body>
</html>`

func classify(bytes []byte, name, side string) string {
	var str string
	if side != "" {
		str = strings.Replace(string(bytes), `<!-- NAV AREA -->`, side, -1)
	} else {
		str = string(bytes)
	}
	if name != "" {
		ps := strings.Replace(
			str,
			"<p>",
			`<p class="`+deSuffix(name)+`">`,
			-1)
		hs := strings.Replace(
			ps,
			"<h1>",
			`<p class="`+deSuffix(name)+`">`,
			-1,
		)
		return hs
	}
	return string(bytes)

}

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
	head := top(argCat(), jsbdir, cssbdir)
	for _, file := range files {
		if strings.Contains(file, "/") {
			if err := os.MkdirAll(filepath.Join(builddir, filepath.Dir(file)), 0755); err == nil {
				if sitefile, err := os.Create(filepath.Join(builddir, filepath.Dir(file), filepath.Base(file)+".html")); err == nil {
					if bytes, err := ioutil.ReadFile(file); err == nil {
						sitefile.Write([]byte(md.RenderToString(bytes)))
					}
					sitefile.Close()
				}
			}
		} else {
			if sitefile, err := os.Create(filepath.Join(builddir, file+".html")); err == nil {
				if bytes, err := ioutil.ReadFile(file); err == nil {
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
	var sidebar string
	for _, dir := range dirs {
		tmpfile, err := os.Create(filepath.Join(dir, "index.md.html"))
		tmpfile.Write(head)
		count := 0
		if err == nil {
			if files, err := ioutil.ReadDir(dir); err == nil {
				for _, file := range files {
					if len(file.Name()) > 8 {
						if file.Name()[len(file.Name())-8:] == ".md.html" {
							if file.Name() != "index.md.html" {
								if bytes, err := ioutil.ReadFile(filepath.Join(dir, file.Name())); err == nil {
									tmpfile.Write([]byte(classify(bytes, file.Name(), "")))
								}
								os.Remove(filepath.Join(dir, file.Name()))
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
			os.Remove(filepath.Join(dir, "index.md.html"))
		} else {
			fp := filepath.Join(dir, "index.md.html")

			sidebar = `[` + dir + `](` + strings.Replace(strings.Replace(fp, ".md", "", -1), "site/", "", 1) + `)`
			f, _ := ioutil.ReadFile(filepath.Join(dir, "index.md.html"))
			fmt.Println(string(f))
		}

	}
	for _, dir := range dirs {
		tmpfile, err := os.Create(filepath.Join(dir, "index.html"))
		count := 0
		if err == nil {
			if files, err := ioutil.ReadDir(dir); err == nil {
				for _, file := range files {
					if file.Name() == "index.md.html" {
						if bytes, err := ioutil.ReadFile(filepath.Join(dir, file.Name())); err == nil {
							tmpfile.Write([]byte(classify(bytes, file.Name(), md.RenderToString([]byte(sidebar)))))
						}
						//os.Remove(filepath.Join(dir, file.Name()))
						count++
					}
				}
			}
		}
		tmpfile.Close()
		if count == 0 {
			os.Remove(filepath.Join(dir, "index.html"))
		} else {
			os.Remove(filepath.Join(dir, "index.md.html"))
			f, _ := ioutil.ReadFile(filepath.Join(dir, "index.html"))
			fmt.Println(string(f))
		}
	}
}
