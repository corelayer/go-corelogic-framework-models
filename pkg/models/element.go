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
	"strings"
)

type Element struct {
	Name        string     `yaml:"name"`
	Tags        []string   `yaml:"tags,omitempty"`
	Fields      []Field    `yaml:"fields,omitempty"`
	Expressions Expression `yaml:"expressions,omitempty"`
}

func (e *Element) GetFullName(moduleName string) string {
	return moduleName + "." + e.Name
}

func (e *Element) GetFields(moduleName string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, f := range e.Fields {
		outputKey := f.GetFullName(e.GetFullName(moduleName))
		if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
			err = fmt.Errorf("duplicate key in fields: %q", outputKey)
			break
		} else {
			output[outputKey] = f.Data
		}
	}

	return output, err
}

func (e *Element) GetFullyQualifiedExpression(expression string, moduleName string) (string, error) {
	// Expand field names in expression to fully qualified name for element
	expression = strings.ReplaceAll(expression, "<<", "<<"+e.GetFullName(moduleName)+"/")

	fields, err := e.GetFields(moduleName)
	if err == nil {
		// Replace all field names with their actual value
		for k, v := range fields {
			expression = strings.ReplaceAll(expression, "<<"+k+">>", v)
		}
	}

	return expression, err
}

func (e *Element) HasFilteredTag(tagFilter []string) bool {
	filterElement := false
	for _, t := range e.Tags {
		for _, f := range tagFilter {
			if t == f {
				filterElement = true
				break
			}
		}
	}
	return filterElement
}
