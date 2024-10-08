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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	clusteraddonsv1 "sigs.k8s.io/cluster-api/exp/addons/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var linodeobjectstoragekeylog = logf.Log.WithName("linodeobjectstoragekey-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *LinodeObjectStorageKey) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1alpha2-linodeobjectstoragekey,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=linodeobjectstoragekeys,verbs=create;update,versions=v1alpha2,name=validation.linodeobjectstoragekey.infrastructure.cluster.x-k8s.io,admissionReviewVersions=v1

var _ webhook.Validator = &LinodeObjectStorageKey{}

// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha2-linodeobjectstoragekey,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=linodeobjectstoragekeys,verbs=create;update,versions=v1alpha2,name=mutation.linodeobjectstoragekey.infrastructure.cluster.x-k8s.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LinodeObjectStorageKey{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LinodeObjectStorageKey) ValidateCreate() (admission.Warnings, error) {
	linodeobjectstoragekeylog.Info("validate create", "name", r.Name)

	return r.validateLinodeObjectStorageKey()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LinodeObjectStorageKey) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	linodeobjectstoragekeylog.Info("validate update", "name", r.Name)

	return r.validateLinodeObjectStorageKey()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LinodeObjectStorageKey) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

func (r *LinodeObjectStorageKey) validateLinodeObjectStorageKey() (admission.Warnings, error) {
	var errs field.ErrorList

	if r.Spec.GeneratedSecret.Type == clusteraddonsv1.ClusterResourceSetSecretType && len(r.Spec.GeneratedSecret.Format) == 0 {
		errs = append(errs, field.Invalid(
			field.NewPath("spec").Child("generatedSecret").Child("format"),
			r.Spec.GeneratedSecret.Format,
			fmt.Sprintf("must not be empty with Secret type %s", clusteraddonsv1.ClusterResourceSetSecretType),
		))
	}

	if len(errs) > 0 {
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "LinodeObjectStorageKey"}, r.Name, errs)
	}

	return nil, nil
}

const defaultKeySecretNameTemplate = "%s-obj-key"

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LinodeObjectStorageKey) Default() {
	linodeobjectstoragekeylog.Info("default", "name", r.Name)

	// Default name and namespace derived from object metadata.
	if r.Spec.GeneratedSecret.Name == "" {
		r.Spec.GeneratedSecret.Name = fmt.Sprintf(defaultKeySecretNameTemplate, r.Name)
	}
	if r.Spec.GeneratedSecret.Namespace == "" {
		r.Spec.GeneratedSecret.Namespace = r.Namespace
	}

	// Support deprecated fields when specified and updated fields are empty.
	if r.Spec.SecretType != "" && r.Spec.GeneratedSecret.Type == "" {
		r.Spec.GeneratedSecret.Type = r.Spec.SecretType
	}
	if len(r.Spec.SecretDataFormat) > 0 && len(r.Spec.GeneratedSecret.Format) == 0 {
		r.Spec.GeneratedSecret.Format = r.Spec.SecretDataFormat
	}
}
