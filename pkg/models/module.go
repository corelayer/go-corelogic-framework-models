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
)

type Module struct {
	Name     string    `yaml:"name"`
	Tags     []string  `yaml:"tags,omitempty"`
	Sections []Section `yaml:"sections,omitempty"`
}

func (m *Module) GetFullModuleName(packageName string) string {
	return packageName + "." + m.Name
}

func (m *Module) GetElements(packageName string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	fullModuleName := m.GetFullModuleName(packageName)
	for _, s := range m.Sections {
		expressions, err = s.GetElements(fullModuleName)
		if err != nil {
			break
		} else {
			output, err = m.appendData(output, expressions)
		}
	}

	return output, err
}
func (m *Module) GetFields(packageName string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	fullModuleName := m.GetFullModuleName(packageName)
	for _, s := range m.Sections {
		expressions, err = s.GetFields(fullModuleName)
		if err != nil {
			break
		} else {
			output, err = m.appendData(output, expressions)
		}
	}

	return output, err
}

func (m *Module) GetInstallExpressions(packageName string, tagFilter []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	if !m.HasFilteredTag(tagFilter) {
		fullModuleName := m.GetFullModuleName(packageName)
		for _, s := range m.Sections {
			// log.Println(s.Name)
			expressions, err = s.GetInstallExpressions(fullModuleName, tagFilter)
			if err != nil {
				break
			} else {
				output, err = m.appendData(output, expressions)
			}
		}
	}

	return output, err
}

func (m *Module) GetUninstallExpressions(packageName string, tagFilter []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	if !m.HasFilteredTag(tagFilter) {
		fullModuleName := m.GetFullModuleName(packageName)
		for _, s := range m.Sections {
			expressions, err = s.GetUninstallExpressions(fullModuleName, tagFilter)
			if err != nil {
				break
			} else {
				output, err = m.appendData(output, expressions)
			}
		}
	} else {
		log.Printf("Skipping module %s", m.GetFullModuleName(packageName))

	}

	return output, err
}

func (m *Module) appendData(destination map[string]interface{}, source map[string]interface{}) (map[string]interface{}, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in %q", k, m.Name)
			break
		} else {
			destination[k] = v
		}
	}

	return destination, err
}

func (m *Module) HasFilteredTag(tagFilter []string) bool {
	filterModule := false
	for _, t := range m.Tags {
		for _, f := range tagFilter {
			if t == f {
				filterModule = true
				// log.Printf("Skipping module %s for tag %s", m.GetFullModuleName(packageName), t)
				break
			}
		}
	}

	return filterModule
}
