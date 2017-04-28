// Copyright 2017-present The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tplimpl

import (
	"fmt"
	"html/template"
	"strings"
	texttemplate "text/template"

	bp "github.com/gohugoio/hugo/bufferpool"
	"github.com/gohugoio/hugo/deps"

	"github.com/spf13/cast"
)

// Some of the template funcs are'nt entirely stateless.
type templateFuncster struct {
	funcMap template.FuncMap

	*deps.Deps
}

func newTemplateFuncster(deps *deps.Deps) *templateFuncster {
	return &templateFuncster{
		Deps: deps,
	}
}

// Partial executes the named partial and returns either a string,
// when called from text/template, for or a template.HTML.
func (t *templateFuncster) partial(name string, contextList ...interface{}) (interface{}, error) {
	var prefix = "partials"
	if strings.HasPrefix("partials/", name) {
		name = name[8:]
	}
	var context interface{}

	if len(contextList) == 0 {
		context = nil
	} else if pr, err := cast.ToStringE(contextList[0]); err == nil && len(contextList) >= 2 {
		// The first parameter of the list (second of the partial
		// call) is the prefix
		prefix = pr
		context = contextList[1]
	} else {
		context = contextList[0]
	}

	prefix += "/"
	for _, n := range []string{prefix + name, "theme/" + prefix + name} {
		templ := t.Tmpl.Lookup(n)
		if templ == nil {
			// For legacy reasons.
			templ = t.Tmpl.Lookup(n + ".html")
		}
		if templ != nil {
			b := bp.GetBuffer()
			defer bp.PutBuffer(b)

			if err := templ.Execute(b, context); err != nil {
				return "", err
			}

			if _, ok := templ.Template.(*texttemplate.Template); ok {
				return b.String(), nil
			}

			return template.HTML(b.String()), nil

		}
	}

	return "", fmt.Errorf("Partial %q not found", name)
}

/*
// Retrieves and display a widget area using the /widgets/ shortcode
func (t *templateFuncster) widgets(name string, context interface{}) (interface{}, error) {
	// Add (_wa: name) index/value to context to access it inside
	// the embedded template (as Widget Area)
	outcontext := make(map[string]interface{})
	outcontext["c"] = context
	outcontext["_wa"] = name

	// See in template_embedded for widgets.html
	templ := t.Tmpl.Lookup("_internal/widgets.html")
	if templ != nil {
		b := bp.GetBuffer()
		defer bp.PutBuffer(b)

		if err := templ.Execute(b, outcontext); err != nil {
			return "", err
		}

		return template.HTML(b.String()), nil
	}
	return "", fmt.Errorf("Widget area %q not found", name)
}
*/
