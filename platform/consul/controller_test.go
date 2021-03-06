// Copyright 2017 Istio Authors
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

package consul

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"

	"istio.io/core/model"
)

var (
	services = map[string][]string{
		"productpage": {"version|v1"},
		"reviews":     {"version|v1", "version|v2", "version|v3"},
	}
	productpage = []*api.CatalogService{
		{
			Node:           "istio",
			Address:        "172.19.0.5",
			ID:             "111-111-111",
			ServiceName:    "productpage",
			ServiceTags:    []string{"version|v1"},
			ServiceAddress: "172.19.0.11",
			ServicePort:    9080,
		},
	}
	reviews = []*api.CatalogService{
		{
			Node:           "istio",
			Address:        "172.19.0.5",
			ID:             "222-222-222",
			ServiceName:    "reviews",
			ServiceTags:    []string{"version|v1"},
			ServiceAddress: "172.19.0.6",
			ServicePort:    9080,
		},
		{
			Node:           "istio",
			Address:        "172.19.0.5",
			ID:             "333-333-333",
			ServiceName:    "reviews",
			ServiceTags:    []string{"version|v2"},
			ServiceAddress: "172.19.0.7",
			ServicePort:    9080,
		},
		{
			Node:           "istio",
			Address:        "172.19.0.5",
			ID:             "444-444-444",
			ServiceName:    "reviews",
			ServiceTags:    []string{"version|v3"},
			ServiceAddress: "172.19.0.8",
			ServicePort:    9080,
			NodeMeta:       map[string]string{protocolTagName: "tcp"},
		},
	}
)

type mockServer struct {
	Server      *httptest.Server
	Services    map[string][]string
	Productpage []*api.CatalogService
	Reviews     []*api.CatalogService
}

func newServer() *mockServer {
	m := mockServer{
		Productpage: make([]*api.CatalogService, len(productpage)),
		Reviews:     make([]*api.CatalogService, len(reviews)),
		Services:    make(map[string][]string),
	}

	copy(m.Reviews, reviews)
	copy(m.Productpage, productpage)
	for k, v := range services {
		m.Services[k] = v
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/catalog/services" {
			data, _ := json.Marshal(&m.Services)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(data))
		} else if r.URL.Path == "/v1/catalog/service/reviews" {
			data, _ := json.Marshal(&m.Reviews)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(data))
		} else if r.URL.Path == "/v1/catalog/service/productpage" {
			data, _ := json.Marshal(&m.Productpage)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintln(w, string(data))
		} else {
			fmt.Fprintln(w, r.URL.Path)
		}
	}))

	m.Server = server
	return &m
}

func TestInstances(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	hostname := serviceHostname("reviews")
	instances := controller.Instances(hostname, []string{}, model.LabelsCollection{})
	if len(instances) != 3 {
		t.Errorf("Instances() returned wrong # of service instances => %q, want 3", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != hostname {
			t.Errorf("Instances() returned wrong service instance => %v, want %q",
				inst.Service.Hostname, hostname)
		}
	}

	filterTagKey := "version"
	filterTagVal := "v3"
	instances = controller.Instances(hostname, []string{}, model.LabelsCollection{
		model.Labels{filterTagKey: filterTagVal},
	})
	if len(instances) != 1 {
		t.Errorf("Instances() did not filter by tags => %q, want 1", len(instances))
	}
	for _, inst := range instances {
		found := false
		for key, val := range inst.Labels {
			if key == filterTagKey && val == filterTagVal {
				found = true
			}
		}
		if !found {
			t.Errorf("Instances() did not match by tag {%q:%q}", filterTagKey, filterTagVal)
		}
	}

	filterPort := "http"
	instances = controller.Instances(hostname, []string{filterPort}, model.LabelsCollection{})
	if len(instances) != 2 {
		t.Errorf("Instances() did not filter by port => %q, want 2", len(instances))
	}
	for _, inst := range instances {
		if inst.Endpoint.ServicePort.Name != filterPort {
			t.Errorf("Instances() did not filter by port => %q, want %q",
				inst.Endpoint.ServicePort.Name, filterPort)
		}
	}
}

func TestInstancesBadHostname(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	instances := controller.Instances("badhostname", []string{}, model.LabelsCollection{})
	if len(instances) != 0 {
		t.Errorf("Instances() returned wrong # of service instances => %q, want 0", len(instances))
	}
}

func TestGetService(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	service, exists := controller.GetService("productpage.service.consul")
	if !exists {
		t.Error("service should exist")
	}

	if service.Hostname != serviceHostname("productpage") {
		t.Errorf("GetService() incorrect service returned => %q, want %q",
			service.Hostname, serviceHostname("productpage"))
	}
}

func TestGetServiceBadHostname(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	service, exists := controller.GetService("badshostname")
	if exists {
		t.Error("service should not exist")
	}

	if service != nil {
		t.Error("service should be nil")
	}
}

func TestGetServiceNoInstances(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	ts.Productpage = []*api.CatalogService{}

	service, exists := controller.GetService("productpage.service.consul")
	if exists {
		t.Error("service should not exist")
	}

	if service != nil {
		t.Error("service should be nil")
	}
}

func TestServices(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	services := controller.Services()
	serviceMap := make(map[string]*model.Service)
	for _, svc := range services {
		name, err := parseHostname(svc.Hostname)
		if err != nil {
			t.Errorf("Services() error parsing hostname: %v", err)
		}
		serviceMap[name] = svc
	}

	for _, name := range []string{"productpage", "reviews"} {
		if _, exists := serviceMap[name]; !exists {
			t.Errorf("Services() missing: %q", name)
		}
	}
	if len(services) != 2 {
		t.Errorf("Services() returned wrong # of services: %q, want 2", len(services))
	}
}

func TestServicesError(t *testing.T) {
	ts := newServer()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		ts.Server.Close()
		t.Errorf("could not create Consul Controller: %v", err)
	}

	ts.Server.Close()
	services := controller.Services()
	if len(services) != 0 {
		t.Errorf("Services() returned wrong # of services: %q, want 0", len(services))
	}
}

func TestHostInstances(t *testing.T) {
	ts := newServer()
	defer ts.Server.Close()
	controller, err := NewController(ts.Server.URL, "datacenter", 3*time.Second)
	if err != nil {
		t.Errorf("could not create Consul Controller: %v", err)
	}

	services := controller.HostInstances(map[string]bool{"172.19.0.11": true})
	if len(services) != 1 {
		t.Errorf("HostInstances() returned wrong # of endpoints => %q, want 1", len(services))
	}

	if services[0].Service.Hostname != serviceHostname("productpage") {
		t.Errorf("HostInstances() wrong service instance returned => hostname %q, want %q",
			services[0].Service.Hostname, serviceHostname("productpage"))
	}
}
