// Copyright 2016-present The Hugo Authors. All rights reserved.
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

package helpers

import (
	"fmt"

	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/hugofs"
	"github.com/spf13/cast"
)

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
	BaseURL

	disablePathToLower bool
	removePathAccents  bool
	uglyURLs           bool
	canonifyURLs       bool

	language *Language

	// pagination path handling
	paginatePath string

	theme string

	// Directories
	themesDir  string
	layoutDir  string
	widgetsDir string
	workingDir string
	staticDirs []string

	// The PathSpec looks up its config settings in both the current language
	// and then in the global Viper config.
	// Some settings, the settings listed below, does not make sense to be set
	// on per-language-basis. We have no good way of protecting against this
	// other than a "white-list". See language.go.
	defaultContentLanguageInSubdir bool
	defaultContentLanguage         string
	multilingual                   bool

	// The file systems to use
	Fs *hugofs.Fs

	// The config provider to use
	Cfg config.Provider
}

func (p PathSpec) String() string {
	return fmt.Sprintf("PathSpec, language %q, prefix %q, multilingual: %T", p.language.Lang, p.getLanguagePrefix(), p.multilingual)
}

// NewPathSpec creats a new PathSpec from the given filesystems and Language.
func NewPathSpec(fs *hugofs.Fs, cfg config.Provider) (*PathSpec, error) {

	baseURLstr := cfg.GetString("baseURL")
	baseURL, err := newBaseURLFromString(baseURLstr)

	if err != nil {
		return nil, fmt.Errorf("Failed to create baseURL from %q: %s", baseURLstr, err)
	}

	var staticDirs []string

	for i := -1; i <= 10; i++ {
		staticDirs = append(staticDirs, getStringOrStringSlice(cfg, "staticDir", i)...)
	}

	ps := &PathSpec{
		Fs:                             fs,
		Cfg:                            cfg,
		disablePathToLower:             cfg.GetBool("disablePathToLower"),
		removePathAccents:              cfg.GetBool("removePathAccents"),
		uglyURLs:                       cfg.GetBool("uglyURLs"),
		canonifyURLs:                   cfg.GetBool("canonifyURLs"),
		multilingual:                   cfg.GetBool("multilingual"),
		defaultContentLanguageInSubdir: cfg.GetBool("defaultContentLanguageInSubdir"),
		defaultContentLanguage:         cfg.GetString("defaultContentLanguage"),
		paginatePath:                   cfg.GetString("paginatePath"),
		BaseURL:                        baseURL,
		themesDir:                      cfg.GetString("themesDir"),
		layoutDir:                      cfg.GetString("layoutDir"),
		widgetsDir:                     cfg.GetString("widgetsDir"),
		workingDir:                     cfg.GetString("workingDir"),
		staticDirs:                     staticDirs,
		theme:                          cfg.GetString("theme"),
	}

	if language, ok := cfg.(*Language); ok {
		ps.language = language
	}

	return ps, nil
}

func getStringOrStringSlice(cfg config.Provider, key string, id int) []string {

	if id >= 0 {
		key = fmt.Sprintf("%s%d", key, id)
	}

	var out []string

	sd := cfg.Get(key)

	if sds, ok := sd.(string); ok {
		out = []string{sds}
	} else if sd != nil {
		out = cast.ToStringSlice(sd)
	}

	return out
}

// PaginatePath returns the configured root path used for paginator pages.
func (p *PathSpec) PaginatePath() string {
	return p.paginatePath
}

// WorkingDir returns the configured workingDir.
func (p *PathSpec) WorkingDir() string {
	return p.workingDir
}

// StaticDirs returns the relative static dirs for the current configuration.
func (p *PathSpec) StaticDirs() []string {
	return p.staticDirs
}

// LayoutDir returns the relative layout dir in the current configuration.
func (p *PathSpec) LayoutDir() string {
	return p.layoutDir
}

// WidgetsDir returns the relative layout dir in the currenct Hugo project.
func (p *PathSpec) WidgetsDir() string {
	return p.widgetsDir
}

// Theme returns the theme name if set.
func (p *PathSpec) Theme() string {
	return p.theme
}

// Theme returns the theme relative theme dir.
func (p *PathSpec) ThemesDir() string {
	return p.themesDir
}
