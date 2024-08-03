package renderer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/fadyat/ggt/internal"
	"github.com/fadyat/ggt/internal/plugins"
)

type Renderer struct {
	f *internal.Flags
}

func NewRenderer(f *internal.Flags) *Renderer {
	return &Renderer{
		f: f,
	}
}

func (r *Renderer) Render(file *plugins.PluggableFile) error {
	flag, perms := os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(0666)
	if _, err := os.Stat(r.f.OutputFile); err == nil {
		flag, perms = os.O_RDWR|os.O_APPEND, os.FileMode(0644)
	}

	f, err := os.OpenFile(r.f.OutputFile, flag, perms)
	if err != nil {
		return fmt.Errorf("open output file: %w", err)
	}
	defer f.Close()

	return renderTemplate(f, r.f.TmplFile, file)
}

func renderTemplate(out io.Writer, templatePath string, data any) error {
	templateName := filepath.Base(templatePath)
	tmpl, err := template.
		New(templateName).
		Funcs(funcMap()).
		ParseFiles(filepath.Clean(templatePath))

	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	if err = tmpl.ExecuteTemplate(out, templateName, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"collect":       collect,
		"prefix":        prefix,
		"to_got":        toGot,
		"join":          join,
		"generics":      generics,
		"generics_args": genericsArgs,
		"test_call":     testCall,
		"arg_define":    argDefine,
		"call_args":     callArgs,
	}
}

// collect is a helper function, which helps to access the field of the struct
// and store it in the slice of strings.
func collect(field string, slice any) []string {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic("collect: second argument must be a slice")
	}

	var values = make([]string, 0, v.Len())
	for i := 0; i < v.Len(); i++ {
		var elem = v.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		values = append(values, elem.FieldByName(field).String())
	}

	return values
}

// prefix is a helper function, which adds the prefix to each element of the slice.
func prefix(p string, slice []string) []string {
	for i, v := range slice {
		slice[i] = p + v
	}

	return slice
}

// join it is a wrapper around strings.Join, to support the pipeline operator.
func join(sep string, slice []string) string {
	return strings.Join(slice, sep)
}

// toGot is a function, which converts the variables names from the
// want-like names to the got-like names
func toGot(value []string) []string {
	for i, v := range value {
		if strings.HasPrefix(v, "want") {
			value[i] = strings.Replace(v, "want", "got", 1)
		}
	}

	return value
}

// generics is a helper function, which generates the go syntax for the typed arguments.
func generics(value []*internal.Identifier) string {
	if len(value) == 0 {
		return ""
	}

	var args = make([]string, 0, len(value))
	for _, v := range value {
		args = append(args, fmt.Sprintf("%s %s", v.Name, v.Type))
	}

	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}

func genericsArgs(value []*internal.Identifier) string {
	if len(value) == 0 {
		return ""
	}

	// it's hard to determine, which type of generics user
	// wants to use, so we are setting any as a default type.
	var args = make([]string, 0, len(value))
	for range value {
		args = append(args, "any")
	}

	return fmt.Sprintf("[%s]", strings.Join(args, ", "))
}

func testCall(fn *plugins.PluggableFn) string {
	var sb strings.Builder
	if fn.Receiver != nil {
		sb.WriteString(fmt.Sprintf("%s.", fn.Receiver.Name))
	}

	sb.WriteString(fn.Name)
	return sb.String()
}

func argDefine(t string) string {
	if strings.HasPrefix(t, "...") {
		var tt = strings.TrimPrefix(t, "...")
		return fmt.Sprintf("[]%s", tt)
	}

	return t
}

func argCallable(arg *internal.Identifier) string {
	var name = fmt.Sprintf("tt.args.%s", arg.Name)
	if strings.HasPrefix(arg.Type, "...") {
		name += "..."
	}

	return name
}

func callArgs(args []*internal.Identifier) string {
	var call = make([]string, 0, len(args))
	for _, arg := range args {
		call = append(call, argCallable(arg))
	}

	return strings.Join(call, ", ")
}