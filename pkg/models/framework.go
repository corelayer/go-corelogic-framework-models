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
	"fmt"
	"log"
	"sort"
	"strings"
	//"github.com/corelayer/corelogic/general"
)

type DataMapWriter interface {
	appendData(source map[string]string, destination map[string]string) (map[string]string, error)
}

type Framework struct {
	Release  Release   `yaml:"release"`
	Prefixes []Prefix  `yaml:"prefixes"`
	Packages []Package `yaml:"packages,omitempty"`
}

func (f *Framework) GetPrefixMap() map[string]string {
	result := make(map[string]string)

	for _, v := range f.Prefixes {
		result[v.Section] = v.Prefix
	}

	return result
}

func (f *Framework) GetPrefixWithVersion(sectionName string) string {
	return strings.Join([]string{f.GetPrefixMap()[sectionName], f.Release.GetVersionAsString()}, "_")
}

func (f *Framework) appendData(destination map[string]interface{}, source map[string]interface{}) (map[string]interface{}, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in framework", k)
			log.Fatal(err)
		} else {
			destination[k] = v
		}
	}

	return destination, err
}

func (f *Framework) GetElements() (map[string]interface{}, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get fields from packages"))

	output := make(map[string]interface{})
	var err error

	// Get all fields in all packages
	for _, p := range f.Packages {
		var elements map[string]interface{}

		elements, err = p.GetElements()
		if err != nil {
			log.Fatal(err)
			//break
		}

		output, err = f.appendData(output, elements)
		if err != nil {
			log.Fatal(err)
			//break
		}
	}
	return output, err
}

func (f *Framework) GetFields() (map[string]interface{}, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get fields from packages"))

	output := make(map[string]interface{})
	var err error

	// Get all fields in all packages
	for _, p := range f.Packages {
		var fields map[string]interface{}

		fields, err = p.GetFields()
		if err != nil {
			log.Fatal(err)
			//break
		}

		output, err = f.appendData(output, fields)
		if err != nil {
			log.Fatal(err)
			//break
		}
	}
	return output, err
}

func (f *Framework) GetExpressions(kind string, tagFilter []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var err error

	if kind == "install" {
		output, err = f.getInstallExpressions(tagFilter)
	} else if kind == "uninstall" {
		output, err = f.getUninstallExpressions(tagFilter)
	}

	return output, err
}

func (f *Framework) getInstallExpressions(tagFilter []string) (map[string]interface{}, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get install expressions from packages"))

	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	for _, p := range f.Packages {
		expressions, err = p.GetInstallExpressions(tagFilter)
		if err != nil {
			log.Fatal(err)
		} else {
			output, err = f.appendData(output, expressions)
		}
	}

	return output, err
}

func (f *Framework) getUninstallExpressions(tagFilter []string) (map[string]interface{}, error) {
	//defer general.FinishTimer(general.StartTimer("Framework " + f.Release.GetVersionAsString() + " get uninstall expressions from packages"))

	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	for _, p := range f.Packages {
		expressions, err = p.GetUninstallExpressions(tagFilter)
		if err != nil {
			log.Fatal(err)
		} else {
			output, err = f.appendData(output, expressions)
		}
	}

	return output, err
}

func (f *Framework) SortPrefixes(prefixes []Prefix) {
	sort.Slice(prefixes, func(i, j int) bool {
		return prefixes[i].ProcessingOrder < prefixes[j].ProcessingOrder
	})
}
