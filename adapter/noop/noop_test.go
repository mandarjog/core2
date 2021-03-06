// Copyright 2017 Istio Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package noop

// NOTE: This test will eventually be auto-generated so that it automatically supports all templates
//       known to Mixer. For now, it's manually curated.

import (
	"context"
	"reflect"
	"testing"
	"time"

	rpc "github.com/googleapis/googleapis/google/rpc"

	"istio.io/core/pkg/adapter"
	"istio.io/core/pkg/adapter/test"
	"istio.io/core/template/checknothing"
	"istio.io/core/template/listentry"
	"istio.io/core/template/logentry"
	"istio.io/core/template/metric"
	"istio.io/core/template/quota"
	"istio.io/core/template/reportnothing"
)

func TestBasic(t *testing.T) {
	info := GetInfo()

	if !contains(info.SupportedTemplates, checknothing.TemplateName) ||
		!contains(info.SupportedTemplates, reportnothing.TemplateName) ||
		!contains(info.SupportedTemplates, listentry.TemplateName) ||
		!contains(info.SupportedTemplates, logentry.TemplateName) ||
		!contains(info.SupportedTemplates, metric.TemplateName) ||
		!contains(info.SupportedTemplates, quota.TemplateName) {
		t.Error("Didn't find all expected supported templates")
	}

	cfg := info.DefaultConfig
	b := info.NewBuilder().(*builder)
	b.SetAdapterConfig(cfg)

	if err := b.Validate(); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	handler, buildErr := b.Build(context.Background(), test.NewEnv(t))
	if buildErr != nil {
		t.Errorf("Got error %v, expecting success", buildErr)
	}

	checkNothingHandler := handler.(checknothing.Handler)
	if result, err := checkNothingHandler.HandleCheckNothing(context.TODO(), nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else {
		if !reflect.DeepEqual(result.Status, rpc.Status{Code: int32(rpc.OK)}) {
			t.Errorf("Got status %v, expecting %v", result.Status, rpc.Status{Code: int32(rpc.OK)})
		}
		if result.ValidDuration < 1000*time.Second {
			t.Errorf("Got duration of %v, expecting at least 1000 seconds", result.ValidDuration)
		}
		if result.ValidUseCount < 1000 {
			t.Errorf("Got use count of %d, expecting at least 1000", result.ValidUseCount)
		}
	}

	reportNothingHandler := handler.(reportnothing.Handler)
	if err := reportNothingHandler.HandleReportNothing(context.TODO(), nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	listEntryHandler := handler.(listentry.Handler)
	if result, err := listEntryHandler.HandleListEntry(context.TODO(), nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else {
		if !reflect.DeepEqual(result.Status, rpc.Status{Code: int32(rpc.OK)}) {
			t.Errorf("Got status %v, expecting %v", result.Status, rpc.Status{Code: int32(rpc.OK)})
		}
		if result.ValidDuration < 1000*time.Second {
			t.Errorf("Got duration of %v, expecting at least 1000 seconds", result.ValidDuration)
		}
		if result.ValidUseCount < 1000 {
			t.Errorf("Got use count of %d, expecting at least 1000", result.ValidUseCount)
		}
	}

	logEntryHandler := handler.(logentry.Handler)
	if err := logEntryHandler.HandleLogEntry(context.TODO(), nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	metricHandler := handler.(metric.Handler)
	if err := metricHandler.HandleMetric(context.TODO(), nil); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}

	quotaHandler := handler.(quota.Handler)
	if result, err := quotaHandler.HandleQuota(context.TODO(), nil, adapter.QuotaArgs{QuotaAmount: 100}); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	} else {
		if result.ValidDuration < 1000*time.Second {
			t.Errorf("Got duration of %v, expecting at least 1000 seconds", result.ValidDuration)
		}
		if result.Amount != 100 {
			t.Errorf("Got %d quota, expecting 100", result.Amount)
		}
	}

	if err := handler.Close(); err != nil {
		t.Errorf("Got error %v, expecting success", err)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
