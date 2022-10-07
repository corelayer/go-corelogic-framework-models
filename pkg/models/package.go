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
)

type Package struct {
	Name    string   `yaml:"name"`
	Modules []Module `yaml:"modules,omitempty"`
}

func (p *Package) GetElements() (map[string]Element, error) {
	output := make(map[string]Element)
	var elements map[string]Element
	var err error

	for _, m := range p.Modules {
		elements, err = m.GetElements(p.Name)
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

func (p *Package) GetFields() (map[string]string, error) {
	output := make(map[string]string)
	var fields map[string]string
	var err error

	for _, m := range p.Modules {
		fields, err = m.GetFields(p.Name)
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

func (p *Package) GetInstallExpressions(tagFilter []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	for _, m := range p.Modules {
		expressions, err = m.GetInstallExpressions(p.Name, tagFilter)
		if err != nil {
			break
		}

		output, err = p.appendData(output, expressions)
		if err != nil {
			break
		}

	}

	return output, err
}

func (p *Package) GetUninstallExpressions(tagFilter []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	var expressions map[string]interface{}
	var err error

	for _, m := range p.Modules {
		expressions, err = m.GetUninstallExpressions(p.Name, tagFilter)
		if err != nil {
			break
		}

		output, err = p.appendData(output, expressions)
		if err != nil {
			break
		}

	}

	return output, err
}

func (p *Package) appendData(destination map[string]interface{}, source map[string]interface{}) (map[string]interface{}, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in package %q", k, p.Name)
			break
		} else {
			destination[k] = v
		}
	}

	return destination, err
}
