// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package actions

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/processors"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/common/jsontransform"
)

type parseFields struct {
	fields            []string
	maxDepth          int
	overwriteKeys     bool
	processArray      bool
	end_remove_fields bool
	target            *string
}

type parseConfig struct {
	Fields          []string `config:"fields"`
	OverwriteKeys   bool     `config:"overwrite_keys"`
	Target          *string  `config:"target"`
	EndRemoveFields bool     `config:"end_remove_fields"`
}

var (
	defaultParseConfig = parseConfig{
		EndRemoveFields: true,
		OverwriteKeys:   true,
	}
	metrics_prefix = "metric_"
	logger_name    = "metrics"
)

func init() {
	processors.RegisterPlugin("parse_metrics_fields",
		configChecked(newParseFields,
			requireFields("fields"),
			allowedFields("fields", "max_depth", "overwrite_keys", "process_array", "target", "when")))
}

func newParseFields(c *common.Config) (processors.Processor, error) {
	parseConfig := defaultParseConfig

	err := c.Unpack(&parseConfig)
	if err != nil {
		logp.Warn("Error unpacking config for decode_json_fields")
		return nil, fmt.Errorf("fail to unpack the decode_json_fields configuration: %s", err)
	}

	f := &parseFields{fields: parseConfig.Fields, overwriteKeys: parseConfig.OverwriteKeys, target: parseConfig.Target}
	return f, nil
}

func (f *parseFields) Run(event *beat.Event) (*beat.Event, error) {
	var errs []string

	for _, field := range f.fields {
		data, err := event.GetValue(field)

		if err != nil && errors.Cause(err) != common.ErrKeyNotFound {
			debug("Error trying to GetValue for field : %s in event : %v", field, event)
			errs = append(errs, err.Error())
			continue
		}

		text, ok := data.(string)
		if !ok {
			// ignore non string fields when unmarshaling
			continue
		}

		var output interface{}

		err = unmarshal(f.maxDepth, text, &output, f.processArray)
		if err != nil {
			debug("Error trying to unmarshal %s", text)
			errs = append(errs, err.Error())
			continue
		}

		target := field
		if f.target != nil {
			target = *f.target
		}

		if target != "" {
			_, err = event.PutValue(target, output)
		} else {
			switch t := output.(type) {
			case map[string]interface{}:
				jsontransform.WriteJSONKeys(event, t, f.overwriteKeys)
			default:
				errs = append(errs, "failed to add target to root")
			}
		}

		if err != nil {
			debug("Error trying to Put value %v for field : %s", output, field)
			errs = append(errs, err.Error())
			continue
		}
		value ,isMetricsLog:= getValue(field, &output)
		if isMetricsLog == 1 {
			splitValue(value, event)
			event.Delete(field)
		}
	}

	if len(errs) > 0 {
		return event, fmt.Errorf(strings.Join(errs, ", "))
	}
	return event, nil
}

func splitValue(value string, event *beat.Event) {
	s := strings.Split(value, ",")
	if len(s) > 0 {
		for _, k := range s {
			tmp := strings.Split(k, "=")
			if len(tmp) == 2 {
				event.PutValue(metrics_prefix+strings.Trim(tmp[0], " "), tmp[1])
			}
		}
	}
	return
}

func getValue(key string, fields *interface{}) (value string,isMetricsLog int) {
	value = ""
	isMetricsLog = 0
	switch O := interface{}(*fields).(type) {
	case map[string]interface{}:
		for k, v := range O {
			if k == "logger_name" && v == logger_name {
				isMetricsLog = 1
			}
			if k == key {
				value = v.(string)
			}
		}
	}
	return
}

func (f parseFields) String() string {
	return "parse_metrics_fields=" + strings.Join(f.fields, ", ")
}
