package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/pflag"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/inspection-cli-plugin/pkg/cmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}


func main() {
	config := pflag.String("kubeconfig", filepath.Join(homeDir(), ".kube", "config"), "absolute path to the kubeconfig file")
	pflag.Parse()

	root := cmd.NewCmdInspectionOptions(config, genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr})
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
