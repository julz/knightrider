package knative_test

import (
	"testing"

	"github.com/julz/knife/knative"
	corev1 "k8s.io/api/core/v1"
)

func TestSimpleSecret(t *testing.T) {
	s := knative.NewSecret("my-secret")

	errorIfNotEqual(t, s.ObjectMeta.Name, "my-secret",
		"expected secret to have name '%s' but was '%s'",
	)

	errorIfNotEqual(t, s.TypeMeta.Kind, "Secret",
		"expected secret to have kind '%s' but was '%s'",
	)

	errorIfNotEqual(t, s.TypeMeta.APIVersion, "v1",
		"expected secret to have version '%s' but was '%s'",
	)
}

func TestSecretWithBasicAuth(t *testing.T) {
	s := knative.NewSecret("my-secret", knative.WithBasicAuth("user", "pass"))

	errorIfNotEqual(t, s.StringData["username"], "user", "expected secret to have user '%s' but was '%s'")
	errorIfNotEqual(t, s.StringData["password"], "pass", "expected secret to have user '%s' but was '%s'")
	errorIfNotEqual(t, s.Type, corev1.SecretTypeBasicAuth, "expected secret to have type '%s' but was '%s'")
}

func TestSecretWithSSHAuth(t *testing.T) {
	s := knative.NewSecret("my-secret", knative.WithSSHAuth([]byte("my-private-key-base64")))

	errorIfNotEqual(t, s.Data["ssh-privatekey"], []byte("my-private-key-base64"), "expected secret to have private key '%s' but was '%s'")
	errorIfNotEqual(t, s.Type, corev1.SecretTypeSSHAuth, "expected secret to have type '%s' but was '%s'")
}

func TestSecretWithGitTarget(t *testing.T) {
	s := knative.NewSecret("my-secret", knative.WithGitTarget("gitpub.io"), knative.WithGitTarget("gcbar.io"))

	errorIfNotEqual(t, s.Annotations, map[string]string{
		"build.knative.dev/git-0": "gitpub.io",
		"build.knative.dev/git-1": "gcbar.io",
	}, "expected secret to have annotations '%s' but was '%s'")
}

func TestSecretWithDockerTarget(t *testing.T) {
	s := knative.NewSecret("my-secret", knative.WithDockerTarget("dockerr.com"), knative.WithDockerTarget("harberr.com"))

	errorIfNotEqual(t, s.Annotations, map[string]string{
		"build.knative.dev/docker-0": "dockerr.com",
		"build.knative.dev/docker-1": "harberr.com",
	}, "expected secret to have annotations '%s' but was '%s'")
}
