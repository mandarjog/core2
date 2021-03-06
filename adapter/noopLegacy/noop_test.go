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

package noopLegacy

import (
	"testing"

	"istio.io/core/pkg/adapter"
	"istio.io/core/pkg/adapterManager"
	"istio.io/core/pkg/config"
)

func TestRegisteredForAllAspects(t *testing.T) {
	builders := adapterManager.BuilderMap([]adapter.RegisterFn{Register})
	supKindCnt := uint(config.NumKinds)
	var i uint
	for i = 0; i < supKindCnt; i++ {
		if i == uint(config.Unspecified) {
			continue
		}
		k := config.Kind(i)
		found := false
		for _, noop := range builders {
			found = found || noop.Kinds.IsSet(k)
		}
		if !found {
			t.Errorf("Noop is not registered for kind %s", k)
		}
	}
}
