/*
 * Copyright 2022 CoreLayer BV
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package models

import (
	"sort"
	"strings"
)

type DataMapWriter interface {
	appendData(source map[string]string, destination map[string]string) (map[string]string, error)
}

type Framework struct {
	Release  Release   `yaml:"release"`
	Prefixes []Prefix  `yaml:"prefixes"`
	Packages []Package `yaml:"packages,omitempty"`
}

func (f *Framework) GetPrefixWithVersion(section string) string {
	return strings.Join([]string{section, f.Release.GetVersionAsString()}, "_")
}

func (f *Framework) GetPrefixes() (map[string]Prefix, error) {
	output := make(map[string]Prefix)
	var err error

	for _, prefix := range f.Prefixes {
		prefix.Prefix = f.GetPrefixWithVersion(prefix.Section)
		output, err = AppendPrefixes(output, prefix)
		if err != nil {
			break
		}
	}

	return output, err
}

func (f *Framework) GetElements() (map[string]Element, error) {
	output := make(map[string]Element)
	var elements map[string]Element
	var err error

	for _, p := range f.Packages {
		elements, err = p.GetElements()
		if err != nil {
			break
		}
		output, err = AppendElements(output, elements)
		if err != nil {
			break
		}
	}

	return output, err
}

func (f *Framework) GetFields() (map[string]string, error) {
	output := make(map[string]string)
	var fields map[string]string
	var err error

	for _, p := range f.Packages {
		fields, err = p.GetFields()
		if err != nil {
			break
		}
		output, err = AppendFields(output, fields)
		if err != nil {
			break
		}
	}

	return output, err
}

//
//func (f *Framework) GetExpressions(kind string, tagFilter []string) (map[string]interface{}, error) {
//	output := make(map[string]interface{})
//	var err error
//
//	if kind == "install" {
//		output, err = f.getInstallExpressions(tagFilter)
//	} else if kind == "uninstall" {
//		output, err = f.getUninstallExpressions(tagFilter)
//	}
//
//	return output, err
//}

//func (f *Framework) getInstallExpressions(tagFilter []string) (map[string]interface{}, error) {
//	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get install expressions from packages"))
//
//	output := make(map[string]interface{})
//	var expressions map[string]interface{}
//	var err error
//
//	for _, p := range f.Packages {
//		expressions, err = p.GetInstallExpressions(tagFilter)
//		if err != nil {
//			break
//		}
//
//		output, err = f.appendData(output, expressions)
//		if err != nil {
//			break
//		}
//
//	}
//
//	return output, err
//}
//
//func (f *Framework) getUninstallExpressions(tagFilter []string) (map[string]interface{}, error) {
//	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get uninstall expressions from packages"))
//
//	output := make(map[string]interface{})
//	var expressions map[string]interface{}
//	var err error
//
//	for _, p := range f.Packages {
//		expressions, err = p.GetUninstallExpressions(tagFilter)
//		if err != nil {
//			break
//		}
//
//		output, err = f.appendData(output, expressions)
//		if err != nil {
//			break
//		}
//
//	}
//
//	return output, err
//}

func (f *Framework) SortPrefixes(prefixes []Prefix) {
	sort.Slice(prefixes, func(i, j int) bool {
		return prefixes[i].ProcessingOrder < prefixes[j].ProcessingOrder
	})
}
