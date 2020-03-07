/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "argovue/client/clientset/versioned/typed/argovue.io/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeArgovueV1 struct {
	*testing.Fake
}

func (c *FakeArgovueV1) AppConfigs(namespace string) v1.AppConfigInterface {
	return &FakeAppConfigs{c, namespace}
}

func (c *FakeArgovueV1) Datasets(namespace string) v1.DatasetInterface {
	return &FakeDatasets{c, namespace}
}

func (c *FakeArgovueV1) Services(namespace string) v1.ServiceInterface {
	return &FakeServices{c, namespace}
}

func (c *FakeArgovueV1) Tokens(namespace string) v1.TokenInterface {
	return &FakeTokens{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeArgovueV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}