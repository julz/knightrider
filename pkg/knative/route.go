package knative

import (
	serving "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewRoute creates a new route with the given options
func NewRoute(name string, options ...RouteOption) *serving.Route {
	s := &serving.Route{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "serving.knative.dev/v1alpha1",
			Kind:       "Route",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	for _, o := range options {
		o(&s.Spec)
	}

	return s
}

// RouteOption is a function that can configure a RouteSpec
type RouteOption func(*serving.RouteSpec)

// WithTrafficToRevision adds a traffic specification to a Route pointing at a Revision
func WithTrafficToRevision(name, revision string, percent int) RouteOption {
	return func(r *serving.RouteSpec) {
		r.Traffic = append(r.Traffic, serving.TrafficTarget{
			Name:         name,
			RevisionName: revision,
			Percent:      percent,
		})
	}
}

// WithTraffic adds a traffic specification to a Route pointing at the latest Revision in a Configuration
func WithTrafficToConfiguration(name, configuration string, percent int) RouteOption {
	return func(r *serving.RouteSpec) {
		r.Traffic = append(r.Traffic, serving.TrafficTarget{
			Name:              name,
			ConfigurationName: configuration,
			Percent:           percent,
		})
	}
}
