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

type Section struct {
	Name     string    `yaml:"name"`
	Elements []Element `yaml:"elements"`
}

func (s *Section) GetFullName(moduleName string) string {
	return moduleName + "." + s.Name
}

func (s *Section) expandSectionPrefix(expression string) string {
	return strings.ReplaceAll(expression, "prefix", s.Name)
}

func (s *Section) GetElements(moduleName string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		elementOutputName := e.GetFullName(s.GetFullName(moduleName))

		if _, isMapContainsKey := output[elementOutputName]; isMapContainsKey {
			err = fmt.Errorf("duplicate key in fields: %q", elementOutputName)
			break
		} else {
			output[elementOutputName] = elementOutputName
		}
	}

	return output, err
}

func (s *Section) GetFields(moduleName string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		elementOutputName := e.GetFullName(s.GetFullName(moduleName))
		for _, f := range e.Fields {
			outputKey := elementOutputName + "/" + f.Id
			if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
				err = fmt.Errorf("duplicate key in fields: %q", outputKey)
				break
			} else {
				output[outputKey] = s.expandSectionPrefix(f.Data)
			}
		}
	}

	return output, err
}

func (s *Section) GetInstallExpressions(moduleName string, tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		if !e.HasFilteredTag(tagFilter) {
			outputKey := e.GetFullName(s.GetFullName(moduleName))
			var outputValue string
			outputValue, err = e.GetFullyQualifiedExpression(e.Expressions.Install, s.GetFullName(moduleName))

			if err != nil {
				break
			} else {
				if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
					//key exist
					err = fmt.Errorf("duplicate key in section: %q", outputKey)
					break
				} else {
					output[outputKey] = s.expandSectionPrefix(outputValue)
				}
			}
		}
	}

	return output, err
}

func (s *Section) GetUninstallExpressions(moduleName string, tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var err error

	for _, e := range s.Elements {
		if !e.HasFilteredTag(tagFilter) {
			outputKey := e.GetFullName(s.GetFullName(moduleName))
			var outputValue string
			outputValue, err = e.GetFullyQualifiedExpression(e.Expressions.Uninstall, s.GetFullName(moduleName))

			if err != nil {
				break
			} else {
				if _, isMapContainsKey := output[outputKey]; isMapContainsKey {
					//key exist
					err = fmt.Errorf("duplicate key in section: %q", outputKey)
					break
				} else {
					output[outputKey] = s.expandSectionPrefix(outputValue)
				}
			}
		}
	}

	return output, err
}
