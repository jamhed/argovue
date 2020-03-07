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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "argovue/apis/argovue.io/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// DatasetLister helps list Datasets.
type DatasetLister interface {
	// List lists all Datasets in the indexer.
	List(selector labels.Selector) (ret []*v1.Dataset, err error)
	// Datasets returns an object that can list and get Datasets.
	Datasets(namespace string) DatasetNamespaceLister
	DatasetListerExpansion
}

// datasetLister implements the DatasetLister interface.
type datasetLister struct {
	indexer cache.Indexer
}

// NewDatasetLister returns a new DatasetLister.
func NewDatasetLister(indexer cache.Indexer) DatasetLister {
	return &datasetLister{indexer: indexer}
}

// List lists all Datasets in the indexer.
func (s *datasetLister) List(selector labels.Selector) (ret []*v1.Dataset, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Dataset))
	})
	return ret, err
}

// Datasets returns an object that can list and get Datasets.
func (s *datasetLister) Datasets(namespace string) DatasetNamespaceLister {
	return datasetNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DatasetNamespaceLister helps list and get Datasets.
type DatasetNamespaceLister interface {
	// List lists all Datasets in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Dataset, err error)
	// Get retrieves the Dataset from the indexer for a given namespace and name.
	Get(name string) (*v1.Dataset, error)
	DatasetNamespaceListerExpansion
}

// datasetNamespaceLister implements the DatasetNamespaceLister
// interface.
type datasetNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Datasets in the indexer for a given namespace.
func (s datasetNamespaceLister) List(selector labels.Selector) (ret []*v1.Dataset, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Dataset))
	})
	return ret, err
}

// Get retrieves the Dataset from the indexer for a given namespace and name.
func (s datasetNamespaceLister) Get(name string) (*v1.Dataset, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("dataset"), name)
	}
	return obj.(*v1.Dataset), nil
}
