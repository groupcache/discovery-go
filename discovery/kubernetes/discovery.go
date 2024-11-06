/*
 * Copyright 2024 Arsene Tochemey Gandote
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kubernetes

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"

	goset "github.com/deckarep/golang-set/v2"
	"go.uber.org/atomic"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/strings/slices"

	"github.com/groupcache/discovery-go/discovery"
)

// Discovery represents the kubernetes discovery
type Discovery struct {
	config *Config
	client kubernetes.Interface
	mu     sync.Mutex

	stopChan chan struct{}
	// states whether the actor system has started or not
	initialized *atomic.Bool
}

// enforce compilation error
var _ discovery.Provider = &Discovery{}

// NewDiscovery returns an instance of the kubernetes discovery provider
func NewDiscovery(config *Config) *Discovery {
	// create an instance of
	discovery := &Discovery{
		mu:          sync.Mutex{},
		stopChan:    make(chan struct{}, 1),
		initialized: atomic.NewBool(false),
		config:      config,
	}

	return discovery
}

// ID returns the discovery provider id
func (d *Discovery) ID() string {
	return "kubernetes"
}

// Initialize initializes the plugin: registers some internal data structures, clients etc.
func (d *Discovery) Initialize() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.initialized.Load() {
		return discovery.ErrAlreadyInitialized
	}

	return d.config.Validate()
}

// Register registers this node to a service discovery directory.
func (d *Discovery) Register() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.initialized.Load() {
		return discovery.ErrAlreadyRegistered
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		return fmt.Errorf("failed to get the in-cluster config of the kubernetes provider: %w", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create the kubernetes client api: %w", err)
	}

	d.client = client
	d.initialized = atomic.NewBool(true)
	return nil
}

// Deregister removes this node from a service discovery directory.
func (d *Discovery) Deregister() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.initialized.Load() {
		return discovery.ErrNotInitialized
	}
	d.initialized = atomic.NewBool(false)
	close(d.stopChan)
	return nil
}

// DiscoverPeers returns a list of known nodes.
func (d *Discovery) DiscoverPeers() ([]string, error) {
	if !d.initialized.Load() {
		return nil, discovery.ErrNotInitialized
	}

	ctx := context.Background()

	pods, err := d.client.CoreV1().Pods(d.config.Namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(d.config.PodLabels).String(),
	})

	if err != nil {
		return nil, err
	}

	validPortNames := []string{d.config.DiscoveryPortName, d.config.PortName}

	// define the addresses list
	addresses := goset.NewSet[string]()

MainLoop:
	for _, pod := range pods.Items {
		pod := pod

		if pod.Status.Phase != corev1.PodRunning {
			continue MainLoop
		}
		// If there is a Ready condition available, we need that to be true.
		// If no ready condition is set, then we accept this pod regardless.
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.Status != corev1.ConditionTrue {
				continue MainLoop
			}
		}

		// iterate the pod containers and find the named port
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				if !slices.Contains(validPortNames, port.Name) {
					continue
				}

				if port.Name == d.config.DiscoveryPortName {
					addresses.Add(net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(int(port.ContainerPort))))
				}
			}
		}
	}
	return addresses.ToSlice(), nil
}

// Close closes the provider
func (d *Discovery) Close() error {
	return nil
}
