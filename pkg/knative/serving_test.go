package knative_test

import (
	"testing"

	"github.com/julz/knife/knative"
)

func TestSimpleService(t *testing.T) {
	s := knative.NewRunLatestService("foo")

	errorIfNotEqual(t, s.ObjectMeta.Name, "foo",
		"expected service to have name '%s' but was '%s'",
	)

	errorIfNotEqual(t, s.TypeMeta.Kind, "Service",
		"expected service to have kind '%s' but was '%s'",
	)
}

func TestRunLatestService(t *testing.T) {
	s := knative.NewRunLatestService("foo", knative.WithBuild(
		knative.WithBuildTemplate("buildpack", nil, nil),
	))

	errorIfNotEqual(t, s.Spec.RunLatest.Configuration.Build.Template.Name, "buildpack", "expected service to have runLatest type with build template '%s' but was '%s'")
}

func TestPinnedService(t *testing.T) {
	s := knative.NewPinnedService("foo", "revision", knative.WithBuild(
		knative.WithBuildTemplate("buildpack", nil, nil),
	))

	errorIfNotEqual(t, s.Spec.Pinned.RevisionName, "revision", "expected service to have revision name '%s' but was '%s'")

	errorIfNotEqual(t, s.Spec.Pinned.Configuration.Build.Template.Name, "buildpack", "expected service to have runLatest type with build template '%s' but was '%s'")
}
