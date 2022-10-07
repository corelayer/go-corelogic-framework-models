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
	"regexp"
	"strings"
)

func AppendPrefixes(destination map[string]Prefix, source Prefix) (map[string]Prefix, error) {
	var err error

	if _, isMapContainsKey := destination[source.Section]; isMapContainsKey {
		err = fmt.Errorf("duplicate key %q found in framework", source.Section)
	} else {
		destination[source.Section] = source
	}

	return destination, err
}

func AppendElements(destination map[string]Element, source map[string]Element) (map[string]Element, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in framework", k)
			break
		} else {
			destination[k] = v
		}
	}

	return destination, err
}

func AppendFields(destination map[string]string, source map[string]string) (map[string]string, error) {
	var err error

	for k, v := range source {
		if _, isMapContainsKey := destination[k]; isMapContainsKey {
			err = fmt.Errorf("duplicate key %q found in framework", k)
			break
		} else {
			destination[k] = v
		}
	}

	return destination, err
}

func UnfoldFields(fields map[string]string, prefixes map[string]Prefix) map[string]string {
	fields = unfoldFieldData(fields)
	fields = unfoldPrefixes(fields, prefixes)

	return fields
}

func unfoldFieldData(fields map[string]string) map[string]string {
	re := regexp.MustCompile(`<<[a-zA-Z0-9_.]*/[a-zA-Z0-9_]*>>`)
	for key := range fields {
		loop := true
		for loop {
			foundKeys := re.FindAllString(fields[key], -1)
			for _, foundKey := range foundKeys {
				searchKey := strings.ReplaceAll(foundKey, "<<", "")
				searchKey = strings.ReplaceAll(searchKey, ">>", "")
				fields[key] = strings.ReplaceAll(fields[key], foundKey, fields[searchKey])
			}

			if !re.MatchString(fields[key]) {
				loop = false
			}
		}
	}

	return fields
}

func unfoldPrefixes(fields map[string]string, prefixes map[string]Prefix) map[string]string {
	re := regexp.MustCompile(`<<[a-zA-Z0-9_.]*>>`)
	for key := range fields {
		loop := true
		for loop {
			foundKeys := re.FindAllString(fields[key], -1)
			for _, foundKey := range foundKeys {
				searchKey := strings.ReplaceAll(foundKey, "<<", "")
				searchKey = strings.ReplaceAll(searchKey, ">>", "")
				fields[key] = strings.ReplaceAll(fields[key], foundKey, prefixes[searchKey].Prefix)
			}

			if !re.MatchString(fields[key]) {
				loop = false
			}
		}
	}

	return fields

}
