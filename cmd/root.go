package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fuseml/fuseml/cli/cmd/internal/client"
	"github.com/fuseml/fuseml/cli/kubernetes/config"
	pconfig "github.com/fuseml/fuseml/cli/paas/config"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Version = "0.1"
)

var (
	flagConfigFile string
	kubeconfig     string
)

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	ExitfIfError(checkDependencies(), "Cannot operate")

	rootCmd := &cobra.Command{
		Use:           "fuseml",
		Short:         "Fuseml cli",
		Long:          `fuseml cli is the official command line interface for Fuseml PaaS `,
		Version:       fmt.Sprintf("%s", Version),
		SilenceErrors: true,
	}

	pf := rootCmd.PersistentFlags()
	argToEnv := map[string]string{}

	pf.StringVarP(&flagConfigFile, "config-file", "", pconfig.DefaultLocation(),
		"set path of configuration file")
	viper.BindPFlag("config-file", pf.Lookup("config-file"))
	argToEnv["config-file"] = "FUSEML_CONFIG"

	config.KubeConfigFlags(pf, argToEnv)
	config.LoggerFlags(pf, argToEnv)

	pf.IntP("verbosity", "", 0, "Only print progress messages at or above this level (0 or 1, default 0)")
	viper.BindPFlag("verbosity", pf.Lookup("verbosity"))
	argToEnv["verbosity"] = "VERBOSITY"

	client.CmdPush.Flags().StringVarP(&client.FlagServe, "serve", "s", "", "inference service to serve the model (kfserving, seldon_mlflow, seldon_sklearn, knative, deployment)")

	config.AddEnvToUsage(rootCmd, argToEnv)

	rootCmd.AddCommand(CmdCompletion)
	rootCmd.AddCommand(client.CmdInstall)
	rootCmd.AddCommand(client.CmdUninstall)
	rootCmd.AddCommand(client.CmdInfo)
	rootCmd.AddCommand(client.CmdOrgs)
	rootCmd.AddCommand(client.CmdCreateOrg)
	rootCmd.AddCommand(client.CmdPush)
	rootCmd.AddCommand(client.CmdDeleteApp)
	rootCmd.AddCommand(client.CmdApps)
	rootCmd.AddCommand(client.CmdTarget)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func checkDependencies() error {
	ok := true

	dependencies := []struct {
		CommandName string
	}{
		{CommandName: "kubectl"},
		{CommandName: "helm"},
		{CommandName: "sh"},
		{CommandName: "git"},
		{CommandName: "openssl"},
	}

	for _, dependency := range dependencies {
		_, err := exec.LookPath(dependency.CommandName)
		if err != nil {
			fmt.Println(emoji.Sprintf(":fire:Not found: %s", dependency.CommandName))
			ok = false
		}
	}

	if ok {
		return nil
	}

	return errors.New("Please check your PATH, some of our dependencies were not found")
}
