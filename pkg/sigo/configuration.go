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

	"gopkg.in/yaml.v2"
)

type Rules struct {
	Name string `yaml:"name"`
}

type Definition struct {
	Version     string   `yaml:"version"`
	K           int      `yaml:"kAnonymity"`
	L           int      `yaml:"lDiversity"`
	QI          []string `yaml:"quasiIdentifiers"`
	Sensitive   []string `yaml:"sensitives"`
	Aggregation string   `yaml:"aggregation"`
	Rules       []Rules  `yaml:"rules"`
}

func LoadConfigurationFromYAML(filename string) (Definition, error) {
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		//nolint: wrapcheck
		return Definition{}, err
	}

	var conf Definition
	err = yaml.Unmarshal(source, &conf)

	if err != nil {
		return conf, fmt.Errorf("%w", err)
	}

	return conf, nil
}
