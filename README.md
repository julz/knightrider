knife /kay-nife/ - a tool for knative yaml wrangling

# What dis?

It's a little tool for knative. Instead of being a whole CLI, it just gives you some commands to spit out yml which you can pipe to `kubectl apply`. This is kind of nice, especially for development, cos you get to see what's happening and mess with stuff.

It's pronounced _kay-nife_. Sorry.

# Examples?

Sure. To create a build, you can do:

~~~~
knife generate-build -t buildpack -a arg1=foo -a arg2=bar -s github.com/foo/bar mybuild | kubectl apply -f -
~~~~

(Please remember to pronounce this "kay-nife" in your head).

To set up a source-to-service build you can do:

~~~~
knife generate-service myservice -s github.com/foo/bar -t buildpack -a IMAGE=docker.io/busybox docker.io/busybox echo hello | kubectl apply -f -
~~~~

Now let's do some routing!

~~~~
knife generate-route my-route -r revision1:100 -c configuration1:0:v2 | kubectl apply -f -
knife generate-route my-route -r revision1:80 -c configuration1:20:v2 | kubectl apply -f -
knife generate-route my-route -c configuration2:100 | kubectl apply -f -
~~~~

Similar stuff works for most other things.

# What about Secrets and ServiceAccounts?

Sure!

~~~~
knife generate-secret the-secret -t git:github.com -t docker:docker.io | kubectl apply -f -
> Username: ...
> Password: ...

knife generate-service-account buildbot -s the-secret | kubectl apply -f -

knife generate-build mybuild --service-account buildbot -s github.com/julz/myapp -t kaniko | kubectl apply -f -
~~~~

# Anything else?

You can also use knife as a nice go library for building knative yml. e.g.

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
