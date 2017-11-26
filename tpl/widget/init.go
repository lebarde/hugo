// Copyright 2017 The Hugo Authors. All rights reserved.
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

package widget

import (
	"github.com/gohugoio/hugo/deps"
	"github.com/gohugoio/hugo/tpl/internal"
)

const name = "widget"

func init() {
	f := func(d *deps.Deps) *internal.TemplateFuncsNamespace {
		ctx := New(d)

		// TODO modify the template _internal/widgets.html to
		// allow output even if no widget area have been set in
		// config file (see tpl/tplimpl/template_embedded.go).
		/*examples := [][1]string{
			{`{{ widgets "footer" . }}`, `<div class="widget-area widget-area-footer">YOUR WIDGETS GENERATED HERE</div>`},
		}*/

		ns := &internal.TemplateFuncsNamespace{
			Name:    name,
			Context: func(...interface{}) interface{} { return ctx },
		}

		ns.AddMethodMapping(ctx.Widgets,
			[]string{"widgets"},
			[][2]string{},
		)

		return ns

	}

	internal.AddTemplateFuncsNamespace(f)
}
