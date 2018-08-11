package knative

import (
	build "github.com/knative/build/pkg/apis/build/v1alpha1"
	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRunLatestService generates a new serving.Service which tracks the latest revision
func NewRunLatestService(name string, options ...ConfigurationOption) *serving.Service {
	s := defaultService(name)
	s.Spec.RunLatest = &serving.RunLatestType{}
	for _, o := range options {
		o(&s.Spec.RunLatest.Configuration)
	}

	return s
}

// NewPinnedService generates a new serving.Service pinned to a particular revision
func NewPinnedService(name, revisionName string, options ...ConfigurationOption) *serving.Service {
	s := defaultService(name)
	s.Spec.Pinned = &serving.PinnedType{
		RevisionName: revisionName,
	}
	for _, o := range options {
		o(&s.Spec.Pinned.Configuration)
	}

	return s
}

// NewConfiguration generates a new configuration with the given name and options
func NewConfiguration(name string, options ...ConfigurationOption) *serving.Configuration {
	c := &serving.Configuration{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "serving.knative.dev/v1alpha1",
			Kind:       "Configuration",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	for _, o := range options {
		o(&c.Spec)
	}

	return c
}

// ConfigurationOption is an option that can configure a ConfigurationSpec
type ConfigurationOption func(*serving.ConfigurationSpec)

// WithBuild adds a BuildSpec to the Configuration
func WithBuild(options ...BuildSpecOption) ConfigurationOption {
	return func(t *serving.ConfigurationSpec) {
		t.Build = &build.BuildSpec{}

		for _, o := range options {
			o(t.Build)
		}
	}
}

// WithRevisionTemplate adds a RevisionTemplate to the ConfigurationSpec
func WithRevisionTemplate(image string, args []string, env map[string]string) ConfigurationOption {
	return func(t *serving.ConfigurationSpec) {
		var cenv []corev1.EnvVar
		for k, v := range env {
			cenv = append(cenv, corev1.EnvVar{Name: k, Value: v})
		}

		t.RevisionTemplate.Spec.Container.Image = image
		t.RevisionTemplate.Spec.Container.Args = args
		t.RevisionTemplate.Spec.Container.Env = cenv
	}
}

// WithSingleConcurrency sets the RevisionRequestConcurrencyModel to Single
func WithSingleConcurrency(s *serving.ConfigurationSpec) {
	s.RevisionTemplate.Spec.ConcurrencyModel = serving.RevisionRequestConcurrencyModelSingle
}

// WithMultiConcurrency sets the RevisionRequestConcurrencyModel to Multi
func WithMultiConcurrency(s *serving.ConfigurationSpec) {
	s.RevisionTemplate.Spec.ConcurrencyModel = serving.RevisionRequestConcurrencyModelMulti
}

// WithImagePullPolicyAlways sets the image pull policy to Always
func WithImagePullPolicyAlways(s *serving.ConfigurationSpec) {
	s.RevisionTemplate.Spec.Container.ImagePullPolicy = corev1.PullAlways
}

func defaultService(name string) *serving.Service {
	return &serving.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "serving.knative.dev/v1alpha1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}
