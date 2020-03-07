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

// AppConfigsGetter has a method to return a AppConfigInterface.
// A group's client should implement this interface.
type AppConfigsGetter interface {
	AppConfigs(namespace string) AppConfigInterface
}

// AppConfigInterface has methods to work with AppConfig resources.
type AppConfigInterface interface {
	Create(*v1.AppConfig) (*v1.AppConfig, error)
	Update(*v1.AppConfig) (*v1.AppConfig, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.AppConfig, error)
	List(opts metav1.ListOptions) (*v1.AppConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.AppConfig, err error)
	AppConfigExpansion
}

// appConfigs implements AppConfigInterface
type appConfigs struct {
	client rest.Interface
	ns     string
}

// newAppConfigs returns a AppConfigs
func newAppConfigs(c *ArgovueV1Client, namespace string) *appConfigs {
	return &appConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the appConfig, and returns the corresponding appConfig object, and an error if there is any.
func (c *appConfigs) Get(name string, options metav1.GetOptions) (result *v1.AppConfig, err error) {
	result = &v1.AppConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AppConfigs that match those selectors.
func (c *appConfigs) List(opts metav1.ListOptions) (result *v1.AppConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.AppConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("appconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested appConfigs.
func (c *appConfigs) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("appconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a appConfig and creates it.  Returns the server's representation of the appConfig, and an error, if there is any.
func (c *appConfigs) Create(appConfig *v1.AppConfig) (result *v1.AppConfig, err error) {
	result = &v1.AppConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("appconfigs").
		Body(appConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a appConfig and updates it. Returns the server's representation of the appConfig, and an error, if there is any.
func (c *appConfigs) Update(appConfig *v1.AppConfig) (result *v1.AppConfig, err error) {
	result = &v1.AppConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("appconfigs").
		Name(appConfig.Name).
		Body(appConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the appConfig and deletes it. Returns an error if one occurs.
func (c *appConfigs) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *appConfigs) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("appconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched appConfig.
func (c *appConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.AppConfig, err error) {
	result = &v1.AppConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("appconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}