// Copyright (C) 2022 CGI France
//
// This file is part of SIGO.
//
// SIGO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// SIGO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with SIGO.  If not, see <http://www.gnu.org/licenses/>.

package sigo

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Rules struct {
	Name string `yaml:"name"`
}

type Definition struct {
	Version     string   `yaml:"version"`
	K           int      `yaml:"kAnonymity"`
	L           int      `yaml:"lDiversity"`
	Sensitive   []string `yaml:"sensitives"`
	Aggregation string   `yaml:"aggregation"`
	Rules       []Rules  `yaml:"rules"`
}

// LoadConfigurationFromYAML returns the configuration of the yaml file in a Definition object.
func LoadConfigurationFromYAML(filename string) (Definition, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return Definition{}, fmt.Errorf("%w", err)
	}

	var conf Definition
	err = yaml.Unmarshal(source, &conf)

	if err != nil {
		return conf, fmt.Errorf("%w", err)
	}

	return conf, nil
}

// Exist return true if the file is present in the current directory.
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		// Checking if the given file exists or not
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// Contains return true if str is in slice.
func Contains(slice []string, str string) bool {
	for i := range slice {
		if slice[i] == str {
			return true
		}
	}

	return false
}
