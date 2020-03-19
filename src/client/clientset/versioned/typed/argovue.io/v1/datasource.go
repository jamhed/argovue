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

package v1

import (
	scheme "argovue/client/clientset/versioned/scheme"
	"time"

	v1 "argovue/apis/argovue.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DatasourcesGetter has a method to return a DatasourceInterface.
// A group's client should implement this interface.
type DatasourcesGetter interface {
	Datasources(namespace string) DatasourceInterface
}

// DatasourceInterface has methods to work with Datasource resources.
type DatasourceInterface interface {
	Create(*v1.Datasource) (*v1.Datasource, error)
	Update(*v1.Datasource) (*v1.Datasource, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Datasource, error)
	List(opts metav1.ListOptions) (*v1.DatasourceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Datasource, err error)
	DatasourceExpansion
}

// datasources implements DatasourceInterface
type datasources struct {
	client rest.Interface
	ns     string
}

// newDatasources returns a Datasources
func newDatasources(c *ArgovueV1Client, namespace string) *datasources {
	return &datasources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the datasource, and returns the corresponding datasource object, and an error if there is any.
func (c *datasources) Get(name string, options metav1.GetOptions) (result *v1.Datasource, err error) {
	result = &v1.Datasource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("datasources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Datasources that match those selectors.
func (c *datasources) List(opts metav1.ListOptions) (result *v1.DatasourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.DatasourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("datasources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested datasources.
func (c *datasources) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("datasources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a datasource and creates it.  Returns the server's representation of the datasource, and an error, if there is any.
func (c *datasources) Create(datasource *v1.Datasource) (result *v1.Datasource, err error) {
	result = &v1.Datasource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("datasources").
		Body(datasource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a datasource and updates it. Returns the server's representation of the datasource, and an error, if there is any.
func (c *datasources) Update(datasource *v1.Datasource) (result *v1.Datasource, err error) {
	result = &v1.Datasource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("datasources").
		Name(datasource.Name).
		Body(datasource).
		Do().
		Into(result)
	return
}

// Delete takes name of the datasource and deletes it. Returns an error if one occurs.
func (c *datasources) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("datasources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *datasources) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("datasources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched datasource.
func (c *datasources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Datasource, err error) {
	result = &v1.Datasource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("datasources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
