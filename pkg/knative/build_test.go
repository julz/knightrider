package knative_test

import (
	"reflect"
	"testing"

	"github.com/julz/knife/knative"
	build "github.com/knative/build/pkg/apis/build/v1alpha1"

	corev1 "k8s.io/api/core/v1"
)

func TestSimpleBuild(t *testing.T) {
	b := knative.NewBuild("foo")

	errorIfNotEqual(t, b.ObjectMeta.Name, "foo",
		"expected build to have name '%s' but was '%s'",
	)

	errorIfNotEqual(t, b.TypeMeta.Kind, "Build",
		"expected build to have kind '%s' but was '%s'",
	)

	errorIfNotEqual(t, b.TypeMeta.APIVersion, "build.knative.dev/v1alpha1",
		"expected build to have version '%s' but was '%s'",
	)
}

func TestBuildWithGitSource(t *testing.T) {
	b := knative.NewBuild("foo", knative.WithGitSource("github.com/foo/bar", "master"))

	if b.Spec.Source == nil {
		t.Fatalf("expected build spec to have a source")
	}

	if b.Spec.Source.Git == nil {
		t.Fatalf("expected build spec to have source type git")
	}

	errorIfNotEqual(t, b.Spec.Source.Git.Url, "github.com/foo/bar", "expected build spec to have source url '%s' but was '%s'")

	errorIfNotEqual(t, b.Spec.Source.Git.Revision, "master", "expected build spec to have source revision '%s' but was '%s'")
}

func TestBuildWithBuildTemplate(t *testing.T) {
	b := knative.NewBuild("with-build-template", knative.WithBuildTemplate("buildpack", map[string]string{"a": "b"}, map[string]string{"k": "v"}))

	if b.Spec.Template == nil {
		t.Fatalf("expected build spec to have a template")
	}

	errorIfNotEqual(t, b.Spec.Template.Name, "buildpack", "expected build template to have name '%s' but was '%s'")

	errorIfNotEqual(t, b.Spec.Template.Arguments, []build.ArgumentSpec{{Name: "a", Value: "b"}}, "expected build template to have arguments '%s' but was '%s'")

	errorIfNotEqual(t, b.Spec.Template.Env, []corev1.EnvVar{{Name: "k", Value: "v"}}, "expected build template to have arguments '%s' but was '%s'")
}

func TestBuildWithSteps(t *testing.T) {
	b := knative.NewBuild("with-steps", knative.WithStep("step1", "busybox", "echo", "foo"))

	errorIfNotEqual(t, b.Spec.Steps, []corev1.Container{
		{Name: "step1", Image: "busybox", Args: []string{"echo", "foo"}},
	}, "expected build template to have steps '%s' but was '%s'")
}

func errorIfNotEqual(t *testing.T, actual, expected interface{}, msg string) {
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf(msg, expected, actual)
	}
}
