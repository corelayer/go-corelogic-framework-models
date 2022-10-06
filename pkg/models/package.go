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

type Package struct {
	Name    string   `yaml:"name"`
	Modules []Module `yaml:"modules,omitempty"`
}

func (p *Package) GetElements() (map[string]string, error) {
	output := make(map[string]string)
	var fields map[string]string
	var err error

	for _, m := range p.Modules {
		fields, err = m.GetElements(p.Name)
		if err != nil {
			break
		} else {
			output, err = p.AppendData(fields, output)
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
		} else {
			output, err = p.AppendData(fields, output)
		}
	}

	return output, err
}

func (p *Package) GetInstallExpressions(tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var expressions map[string]string
	var err error

	for _, m := range p.Modules {
		expressions, err = m.GetInstallExpressions(p.Name, tagFilter)
		if err != nil {
			break
		} else {
			output, err = p.AppendData(expressions, output)
		}
	}

	return output, err
}

func (p *Package) GetUninstallExpressions(tagFilter []string) (map[string]string, error) {
	output := make(map[string]string)
	var expressions map[string]string
	var err error

	for _, m := range p.Modules {
		expressions, err = m.GetUninstallExpressions(p.Name, tagFilter)
		if err != nil {
			break
		} else {
			output, err = p.AppendData(expressions, output)
		}
	}

	return output, err
}

func (p *Package) AppendData(source map[string]string, destination map[string]string) (map[string]string, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in package %q", k, p.Name)
			log.Fatal(err)
		} else {
			destination[k] = v
		}
	}

	return destination, err
}
