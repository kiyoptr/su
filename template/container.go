package template

import (
	"github.com/kiyoptr/su/datastructures"
	"github.com/kiyoptr/su/errors"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

type Container struct {
	Templates map[string]*template.Template
	Functions template.FuncMap
}

func NewContainer(funcs template.FuncMap) *Container {
	return &Container{map[string]*template.Template{}, funcs}
}

// LoadDirectory recursively loads all files having .template extension as templates.
// Templates can be executed by their names which is a filepath relative to given directory and slashes replaced with
// dots, without file extension. So a file in /tmp/templates/www/index.template with root path of /tmp/templates will
// have www.index template name.
func (t *Container) LoadDirectory(path string) error {
	root := path
	pathStack := datastructures.NewStack(path)

	for !pathStack.IsEmpty() {
		top := pathStack.Pop()
		path = top.(string)

		ls, err := ioutil.ReadDir(path)
		if err != nil {
			return errors.Newif(err, "failed to read directory %s", path)
		}

		for _, fileInfo := range ls {
			if fileInfo.IsDir() {
				pathStack.Push(filepath.Join(path, fileInfo.Name()))
			} else {
				path = filepath.Join(path, fileInfo.Name())
				ext := filepath.Ext(path)

				if ext != ".template" {
					continue
				}

				templateName, err := filepath.Rel(root, path)
				if err != nil {
					return errors.Newif(err, "failed to get relative path for %s", path)
				}

				templateData, err := ioutil.ReadFile(path)
				if err != nil {
					return errors.Newif(err, "failed to read template %s", path)
				}

				templateName = strings.TrimSuffix(templateName, ext)
				templateName = strings.ReplaceAll(templateName, "/", ".")
				t.Templates[templateName], err = template.New(templateName).Parse(string(templateData))
				if err != nil {
					return errors.Newif(err, "failed to parse template %s at %s", templateName, path)
				}
				t.Templates[templateName].Funcs(t.Functions)
			}
		}
	}

	return nil
}

func (t *Container) Parse(name string, data string) (err error) {
	t.Templates[name], err = template.New(name).Parse(data)
	if err != nil {
		err = errors.Newif(err, "failed to parse template string %s", name)
		return
	}
	t.Templates[name].Funcs(t.Functions)
	return
}

func (t *Container) Execute(writer io.Writer, name string, data interface{}) error {
	tmpl, ok := t.Templates[name]
	if !ok {
		return errors.Newf("template %s isn't loaded", name)
	}

	if err := tmpl.Execute(writer, data); err != nil {
		return errors.Newif(err, "failed to execute template %s", name)
	}

	return nil
}
