package knative

import (
	build "github.com/knative/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewBuild creates a new Build object
func NewBuild(name string, options ...BuildSpecOption) *build.Build {
	b := build.Build{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "build.knative.dev/v1alpha1",
			Kind:       "Build",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: build.BuildSpec{},
	}

	for _, option := range options {
		option(&b.Spec)
	}

	return &b
}

// BuildSpecOption is an option that can configure a BuildSpec
type BuildSpecOption func(*build.BuildSpec)

// WithGitSource configures a Build with a Git source
func WithGitSource(url, revision string) BuildSpecOption {
	return func(b *build.BuildSpec) {
		b.Source = &build.SourceSpec{}
		b.Source.Git = &build.GitSourceSpec{
			Url:      url,
			Revision: revision,
		}
	}
}

// WithBuildTemplate configures a BuildTemplate for a Build
func WithBuildTemplate(name string, args map[string]string, env map[string]string) BuildSpecOption {
	return func(b *build.BuildSpec) {
		b.Template = &build.TemplateInstantiationSpec{
			Name: name,
		}

		for name, value := range args {
			b.Template.Arguments = append(b.Template.Arguments, build.ArgumentSpec{
				Name:  name,
				Value: value,
			})
		}

		for name, value := range env {
			b.Template.Env = append(b.Template.Env, corev1.EnvVar{
				Name:  name,
				Value: value,
			})
		}
	}
}

// WithStep adds a Step to a Build
func WithStep(name, image string, args ...string) BuildSpecOption {
	return func(b *build.BuildSpec) {
		b.Steps = append(b.Steps, corev1.Container{
			Name:  name,
			Image: image,
			Args:  args,
		})
	}
}

// WithServiceAccount adds a ServiceAccount to the Build
func WithServiceAccount(name string) BuildSpecOption {
	return func(b *build.BuildSpec) {
		b.ServiceAccountName = name
	}
}
