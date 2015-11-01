/*
 * Copyright (c) 2015 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package layout

import (
	"bytes"
	"html/template"

	"github.com/FooSoft/goldsmith"
)

type layout struct {
	tmpl *template.Template
	def  string
}

func New(glob, def string) goldsmith.Context {
	tmpl, err := template.ParseGlob(glob)
	if err != nil {
		return goldsmith.Context{Err: err}
	}

	return goldsmith.Context{
		Chainer: &layout{tmpl, def},
		Globs:   []string{"*.html", "*.html"},
	}
}

func (t *layout) ChainSingle(file goldsmith.File) goldsmith.File {
	name, ok := file.Meta["layout"]
	if !ok {
		name = t.def
	}

	file.Meta["content"] = template.HTML(file.Buff.Bytes())

	var buff bytes.Buffer
	if file.Err = t.tmpl.ExecuteTemplate(&buff, name.(string), file.Meta); file.Err == nil {
		file.Buff = &buff
	}

	return file
}