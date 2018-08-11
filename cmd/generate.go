package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"

	"github.com/ghodss/yaml"
	"github.com/julz/knife/pkg/knative"
	"github.com/spf13/cobra"
)

var repo, revision, template, serviceAccount string

var generateBuild = &cobra.Command{
	Use:   "generate-build [name]",
	Short: "output a build yml to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printYaml(knative.NewBuild(args[0], buildOptions()...))
	},
}

var templateArgs, templateEnv []string
var single, alwaysPull bool

var generateConfiguration = &cobra.Command{
	Use:   "generate-configuration [name] [image] [args]",
	Short: "output a configuration yml to stdout",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		printYaml(knative.NewConfiguration(args[0], configurationOptions(args[1], args[2:])...))
	},
}

var generateService = &cobra.Command{
	Use:   "generate-service [name] [image] [args]",
	Short: "output a service yml to stdout",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		printYaml(knative.NewRunLatestService(args[0], configurationOptions(args[1], args[2:])...))
	},
}

var revisionTraffic, configurationTraffic []string

var generateRoute = &cobra.Command{
	Use:   "generate-route [name]",
	Short: "output a route yml to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printYaml(knative.NewRoute(args[0], routeOptions()...))
	},
}

var secretTargets []string

var generateSecret = &cobra.Command{
	Use:   "generate-secret [name]",
	Short: "output a secret yml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var options []knative.SecretOption
		for _, t := range secretTargets {
			parts := strings.Split(t, ":")
			switch parts[0] {
			case "git":
				options = append(options, knative.WithGitTarget(parts[1]))
			case "docker":
				options = append(options, knative.WithDockerTarget(parts[1]))
			default:
				fatalF("unrecognised secret target type: %s", parts[0])
			}
		}

		user, pass := readUserPass()
		printYaml(knative.NewSecret(args[0], append(options, knative.WithBasicAuth(user, pass))...))
	},
}

var serviceAccountSecrets []string
var generateServiceAccount = &cobra.Command{
	Use:   "generate-service-account [name]",
	Short: "output a service account yml to stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		printYaml(knative.NewServiceAccount(args[0], knative.WithSecrets(serviceAccountSecrets...)))
	},
}

func init() {
	// everything can take a build, so everything gets the build flags
	for _, cmd := range []*cobra.Command{generateBuild, generateService, generateConfiguration} {
		cmd.Flags().StringVarP(&repo, "git-repo", "u", "", "url of a git repository to use as a source")
		cmd.Flags().StringVarP(&revision, "git-revision", "r", "master", "revision (sha, tag, or branch) to use for the build")

		cmd.Flags().StringVarP(&template, "template", "t", "", "build template name")
		cmd.Flags().StringSliceVarP(&templateArgs, "template-arg", "a", nil, "build template argument in the form name=value")
		cmd.Flags().StringSliceVarP(&templateEnv, "template-env", "e", nil, "build template environment variable in the form name=value")

		cmd.Flags().StringVarP(&serviceAccount, "service-account", "s", "", "service account the build should run using")
	}

	// service and configuration have extra flags to configure the revision template
	for _, cmd := range []*cobra.Command{generateService, generateConfiguration} {
		cmd.Flags().BoolVar(&single, "single", false, "create a single threaded container")
		cmd.Flags().BoolVar(&alwaysPull, "imagePullPolicyAlways", false, "always pull a new version of the image on startup")
	}

	// route has its own arguments: a list of traffic to revisions and configurations
	generateRoute.Flags().StringSliceVarP(&revisionTraffic, "revision", "r", nil, "add traffic to a revision (in format revisionName:percent or revisionName:percent:name")
	generateRoute.Flags().StringSliceVarP(&configurationTraffic, "configuration", "c", nil, "add traffic to a configuration (in format cconfigurationName:percent or cconfigurationName:percent:name")

	// secret takes a list of targets in form [git|docker]:host
	generateSecret.Flags().StringSliceVarP(&secretTargets, "target", "t", nil, "target this secret to a particular host, format git:host or docker:host")

	// serviceAccount takes a list of secret names
	generateServiceAccount.Flags().StringSliceVarP(&serviceAccountSecrets, "secret", "s", nil, "add a secret to the generated account")

	for _, cmd := range []*cobra.Command{generateSecret, generateServiceAccount, generateBuild, generateService, generateConfiguration, generateRoute} {
		root.AddCommand(cmd)
	}
}

func buildOptions() []knative.BuildSpecOption {
	var options []knative.BuildSpecOption
	if repo != "" {
		options = append(options, knative.WithGitSource(repo, revision))
	}

	if template != "" {
		options = append(options, knative.WithBuildTemplate(template, toMap(templateArgs), toMap(templateEnv)))
	}

	if serviceAccount != "" {
		options = append(options, knative.WithServiceAccount(serviceAccount))
	}

	return options
}

func configurationOptions(image string, args []string) []knative.ConfigurationOption {
	options := []knative.ConfigurationOption{
		knative.WithRevisionTemplate(image, args, nil),
	}

	if template != "" {
		options = append(options, knative.WithBuild(buildOptions()...))
	}

	if single {
		options = append(options, knative.WithSingleConcurrency)
	} else {
		options = append(options, knative.WithMultiConcurrency)
	}

	if alwaysPull {
		options = append(options, knative.WithImagePullPolicyAlways)
	}

	return options
}

func routeOptions() []knative.RouteOption {
	var options []knative.RouteOption
	for _, r := range revisionTraffic {
		name, target, percent := parseTrafficSpec(r)
		options = append(options, knative.WithTrafficToRevision(name, target, percent))
	}

	for _, r := range configurationTraffic {
		name, target, percent := parseTrafficSpec(r)
		options = append(options, knative.WithTrafficToConfiguration(name, target, percent))
	}

	return options
}

func parseTrafficSpec(s string) (name, target string, percent int) {
	parts := strings.Split(s, ":")

	if len(parts) == 3 {
		name = parts[2]
	}

	percent, err := strconv.Atoi(parts[1])
	if err != nil {
		fatalF("Error: %s", err)
	}

	return name, parts[0], percent
}

func printYaml(o interface{}) {
	var b []byte
	var err error
	if b, err = yaml.Marshal(o); err != nil {
		fatalF("Error: %s", err)
	}

	fmt.Println(string(b))
}

func toMap(args []string) map[string]string {
	options := make(map[string]string)
	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		options[parts[0]] = parts[1]
	}

	return options
}

func readUserPass() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintf(os.Stderr, "Username: ")
	user, err := reader.ReadString('\n')
	if err != nil {
		fatalF("Error: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Password: ")
	pass, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fatalF("Error: %s", err)
	}

	fmt.Fprintf(os.Stderr, "\n")

	return strings.TrimSpace(user), string(pass)
}

func fatalF(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
