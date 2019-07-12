/**
* Copyright 2018 Comcast Cable Communications Management, LLC
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package influxdb

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/Comcast/trickster/internal/config"
	"github.com/Comcast/trickster/internal/proxy/model"
	"github.com/Comcast/trickster/internal/timeseries"
)

func TestSetExtent(t *testing.T) {

	start := time.Now().Add(time.Duration(-6) * time.Hour)
	end := time.Now()
	expected := "q=select+%2A+where+time+%3E%3D+" + fmt.Sprintf("%d", start.Unix()*1000) + "ms+AND+time+%3C%3D+" + fmt.Sprintf("%d", end.Unix()*1000) + "ms+group+by+time%281m%29"

	err := config.Load("trickster", "test", []string{"-origin", "none:9090", "-origin-type", "influxdb", "-log-level", "debug"})
	if err != nil {
		t.Errorf("Could not load configuration: %s", err.Error())
	}

	oc := config.Origins["default"]
	client := Client{config: oc}

	u := &url.URL{}
	tu := &url.URL{RawQuery: "q=select * where <$TIME_TOKEN$> group by time(1m)"}
	r := &model.Request{URL: u, TemplateURL: tu}
	e := &timeseries.Extent{Start: start, End: end}
	client.SetExtent(r, e)

	if expected != r.URL.RawQuery {
		t.Errorf("\nexpected [%s]\ngot    [%s]", expected, r.URL.RawQuery)
	}
}