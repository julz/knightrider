package knative

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret Creates a new Secret for use with Knative
func NewSecret(name string, options ...SecretOption) corev1.Secret {
	s := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Annotations: make(map[string]string),
		},
		Data:       make(map[string][]byte),
		StringData: make(map[string]string),
	}

	for _, o := range options {
		o(&s)
	}

	return s
}

// SecretOption is a function that can configure a Secret
type SecretOption func(*corev1.Secret)

// WithBasicAuth adds a username and password to the data, and a BasicAuth type
func WithBasicAuth(user, pass string) SecretOption {
	return func(s *corev1.Secret) {
		s.Type = corev1.SecretTypeBasicAuth
		s.StringData["username"] = user
		s.StringData["password"] = pass
	}
}

// WithSSHAuth adds an ssh privaatekey to the secret
func WithSSHAuth(privateKey []byte) SecretOption {
	return func(s *corev1.Secret) {
		s.Type = corev1.SecretTypeSSHAuth
		s.Data["ssh-privatekey"] = privateKey
	}
}

// WithGitTarget marks the Secret as applying to a particular git host
func WithGitTarget(host string) SecretOption {
	return func(s *corev1.Secret) {
		s.Annotations[nextSuffixedKey(s.Annotations, "build.knative.dev/git-")] = host
	}
}

// WithDockerTarget marks the Secret as applying to a particular docker registry host
func WithDockerTarget(host string) SecretOption {
	return func(s *corev1.Secret) {
		s.Annotations[nextSuffixedKey(s.Annotations, "build.knative.dev/docker-")] = host
	}
}

func nextSuffixedKey(m map[string]string, prefix string) string {
	key := fmt.Sprintf("%s%d", prefix, 0)
	for i := 0; hasKey(m, fmt.Sprintf("%s%d", prefix, i)); i++ {
		key = fmt.Sprintf("%s%d", prefix, i+1)
	}

	return key
}

func hasKey(m map[string]string, k string) bool {
	_, ok := m[k]
	return ok
}
