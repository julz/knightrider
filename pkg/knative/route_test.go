package knative_test

import (
	"testing"

	"github.com/julz/knightrider/pkg/knative"
	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
)

func TestSimpleRoute(t *testing.T) {
	r := knative.NewRoute("sudo")

	errorIfNotEqual(t, r.ObjectMeta.Name, "sudo",
		"expected route to have name '%s' but was '%s'",
	)

	errorIfNotEqual(t, r.TypeMeta.Kind, "Route",
		"expected route to have kind '%s' but was '%s'",
	)

	errorIfNotEqual(t, r.TypeMeta.APIVersion, "serving.knative.dev/v1alpha1",
		"expected route to have version '%s' but was '%s'",
	)
}

func TestRouteWithTraffic(t *testing.T) {
	r := knative.NewRoute("sudo", knative.WithTrafficToRevision("name", "revision1", 80), knative.WithTrafficToConfiguration("", "config", 20))

	errorIfNotEqual(t, r.Spec.Traffic, []serving.TrafficTarget{
		{
			Name:         "name",
			RevisionName: "revision1",
			Percent:      80,
		},
		{
			Name:              "",
			ConfigurationName: "config",
			Percent:           20,
		},
	}, "expected route traffic to be '%s' but was '%s'")
}
