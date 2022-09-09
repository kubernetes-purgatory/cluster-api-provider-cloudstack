/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

//+kubebuilder:object:root=true
//+kubebuilder:storageversion

// FakeKindWithInfrastructure is the Schema for testing CRDs with infrastructureTemplate under spec
type FakeKindWithInfrastructureTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec FakeKindWithInfrastructureTemplateSpec `json:"spec,omitempty"`
}

// FakeKindWithInfrastructureSpec defines the desired state of FakeKindWithInfrastructure
type FakeKindWithInfrastructureTemplateSpec struct {
	// +kubebuilder:validation:Required
	InfrastructureTemplate InfrastructureTemplate `json:"infrastructureTemplate,omitempty"`
}

// InfrastructureTemplate defines the name of the template
type InfrastructureTemplate struct {
	Name string `json:"name,omitempty"`
}

//+kubebuilder:object:root=true
// FakeKindWithInfrastructureList contains a list of FakeKindWithInfrastructure
type FakeKindWithInfrastructureTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FakeKindWithInfrastructureTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FakeKindWithInfrastructureTemplate{}, &FakeKindWithInfrastructureTemplateList{})
}
