/*
Copyright 2023 Akamai Technologies, Inc.

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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AddressSetSpec defines the desired state of AddressSet
type AddressSetSpec struct {
	IPv4 *[]string `json:"ipv4,omitempty"`
	IPv6 *[]string `json:"ipv6,omitempty"`
}

// AddressSetStatus defines the observed state of AddressSet
type AddressSetStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=addresssets,scope=Namespaced,categories=cluster-api,shortName=addrset
// +kubebuilder:metadata:labels="clusterctl.cluster.x-k8s.io/move-hierarchy=true"

// AddressSet is the Schema for the addresssets API
type AddressSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AddressSetSpec   `json:"spec,omitempty"`
	Status AddressSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AddressSetList contains a list of AddressSet
type AddressSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AddressSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AddressSet{}, &AddressSetList{})
}
