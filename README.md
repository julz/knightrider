Knightrider /kay-nightrider/ - a simple tool for working with knative

# What dis?

It's a little tool for knative. Instead of being a whole CLI, it just gives you some commands to spit out yml which you can pipe to `kubectl apply`. This is kind of nice, especially for development, cos you get to see what's happening and mess with stuff. There's also a kubectl plugin to add a bit of sugar to your knative kubectl-ing. As a bit of sugar, you can also do `kr apply/create/delete/update` etc, which just pipes the yml to kubectl for you.

# Use as Plugin

If you're using the latest kubectl, then you can use this as a kubectl plugin! Install and use as follows:

~~~~
go install github.com/julz/knightrider/cmd/kubectl-knative
kubectl knative generate service my-service docker.io/my-repo/my-image | kubectl apply -f -

# or, with sugar:
kubectl knative create service my-service docker.io/my-repo/my-image

# .. with sugar, and watching for completion with nice status output:
kubectl knative create service my-service docker.io/my-repo/my-image --watch
~~~~

# Use as a standalone binary for generating knative yml

If you're not using the latest kubectl, you can just install via `go get github.com/julz/knightrider/cmd/kr` and use via `kr` as described below.

To create a build, you can do:

~~~~
kr generate build -t buildpack -a arg1=foo -a arg2=bar -s github.com/foo/bar mybuild | kubectl apply -f -
~~~~

Or, without the piping:

~~~~
kr create build -t buildpack -a arg1=foo -a arg2=bar -s github.com/foo/bar mybuild
~~~~

(Please remember to pronounce this "kay-nightrider create" in your head).

To set up a source-to-service build you can do:

~~~~
kr generate service myservice -s github.com/foo/bar -t buildpack -a IMAGE=docker.io/busybox docker.io/busybox echo hello | kubectl apply -f -
~~~~

Now let's do some routing!

~~~~
kr generate route my-route -r revision1:100 -c configuration1:0:v2 | kubectl apply -f -
kr generate route my-route -r revision1:80 -c configuration1:20:v2 | kubectl apply -f -
kr generate route my-route -c configuration2:100 | kubectl apply -f -
~~~~

Similar stuff works for most other things.

*TIP*: For a diff showing what will change if you apply a generated object, you can pipe to `kubectl alpha diff -f - LAST LOCAL` instead of `kubectl apply -f -`.

# How about rapid local development?

Glad you asked! You can do a super-nice local-build-and-run-on-cluster for Go programs using the fantastic `ko apply` instead of `kubectl apply`:

~~~~
ko apply -L -f <( kr generate-service hello-world github.com/julz/kr/test/cmd/hello-world )
~~~~

NOTE: `ko apply -f` doesnt support `-` for stdin, so you can't just pipe to `ko apply -f -` :-(

What the above did is generate a Knative Service YML with 'github.com/julz/kr/test/cmd/hello-world/' as the Image, and then use `ko` to turn that in to a YML with a proper docker image URI and apply the manifest. The image gets built in your local minikube's docker (this also works fine for remote clusters, just lose the `-L` in the above command) so it's _blaaazing_ fast. See [the go-containerregistry repo](https://github.com/google/go-containerregistry/tree/master/cmd/ko) for more about `ko`.

Here's how to get a diff of what's about to be changed, you should see the image updated to the new built sha:

~~~~
ko resolve -f <(kr generate-service hello-world github.com/julz/kr/test/cmd/hello-world) | kubectl alpha diff -f - LAST LOCAL
~~~~

# What about Secrets and ServiceAccounts?

Sure!

~~~~
kr generate secret the-secret -t git:github.com -t docker:docker.io | kubectl apply -f -
> Username: ...
> Password: ...

kr generate service-account buildbot -s the-secret | kubectl apply -f -

kr generate build mybuild --service-account buildbot -s github.com/julz/myapp -t kaniko | kubectl apply -f -

# or just..
kr apply build mybuild --service-account buildbot -s github.com/julz/myapp -t kaniko
~~~~

# Anything else?

You can also use knightrider as a nice go library for building knative yml. e.g.

~~~~golang
func main() {
  json.NewEncoder(os.Stdout).Encode( 
    knative.NewBuild("my-build", knative.WithGitSource("github.com/foo/bar", "master"), knative.WithBuildTemplate("buildpack", knative.WithBuildTemplateArg("key", "value")
  )

  json.NewEncoder(os.Stdout).Encode( 
    knative.NewRunLatestService("my-service",
      knative.WithRevisionTemplate("docker.io/busybox", nil, nil), 
      knative.WithGitSource(knative.WithGitSource("github.com/foo/bar", "master"), 
      knative.WithBuildTemplate("buildpack"), 
      knative.WithBuildTemplateArg("key", "value"),
      knative.WithServiceAccount("buildbot"),
    ))
}
~~~~
