package knative

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewServiceAccount creates a new ServiceAccount with the given options
func NewServiceAccount(name string, options ...ServiceAccountOption) corev1.ServiceAccount {
	s := corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	for _, o := range options {
		o(&s)
	}

	return s
}

// ServiceAccountOption is a function that can configure a ServiceAccount
type ServiceAccountOption func(*corev1.ServiceAccount)

func WithSecrets(names ...string) ServiceAccountOption {
	return func(s *corev1.ServiceAccount) {
		for _, n := range names {
			s.Secrets = append(s.Secrets, corev1.ObjectReference{Name: n})
		}
	}
}
