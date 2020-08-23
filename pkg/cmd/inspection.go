package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	inspectionExample = `
	# view the current namespace in your KUBECONFIG
	%[1]s inspection
`

	errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
)

// InspectionOptions provides information required to get deployment and statefulset ...
type InspectionOptions struct {
	args       []string
	kubeConfig string
	kubeCli    *kubernetes.Clientset
	genericclioptions.IOStreams
}

// NewInspectionOptions provides an instance of InspectionOptions with default values
func NewInspectionOptions(streams genericclioptions.IOStreams, config *string) *InspectionOptions {
	return &InspectionOptions{
		kubeConfig: *config,
		IOStreams:  streams,
	}
}

// NewCmdInspectionOptions provides a cobra command wrapping NewInspectionOptions
func NewCmdInspectionOptions(config *string, streams genericclioptions.IOStreams) *cobra.Command {
	o := NewInspectionOptions(streams, config)
	cmd := &cobra.Command{
		Use:          "ns [kubeconfig] [flags]",
		Short:        "View or set the current namespace",
		Example:      fmt.Sprintf(inspectionExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.CreateKubeCli(); err != nil {
				return err
			}
			return o.Run()
		},
	}
	return cmd
}

// Validate ensures that all required arguments and flag values are provided
func (o *InspectionOptions) Validate() error {
	if _, err := os.Stat(o.kubeConfig); os.IsNotExist(err) {
		return err
	}
	return nil
}

// CreateKubeCli ensures that we can connect master successfully
func (o *InspectionOptions) CreateKubeCli() error {
	config, err := clientcmd.BuildConfigFromFlags("", o.kubeConfig)
	if err != nil {
		panic(err.Error())
	}

	o.kubeCli, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

// Run lists all available deployment,statefulset on a user's KUBECONFIG
func (o *InspectionOptions) Run() error {
	nodeResult := o.CheckNodes()
	if nodeResult.totalCount == 0 {
		fmt.Fprintf(o.ErrOut, "Total node count is [%d]\n", nodeResult.totalCount)
	} else {
		fmt.Fprintf(o.Out, "Total node count is [%d]\n", nodeResult.totalCount)
		fmt.Fprintf(o.Out, "Total notReady node count is [%d]\n", len(nodeResult.notReady))
		if len(nodeResult.notReady) >0{
			fmt.Fprintf(o.Out, "notReady nodes: %v \n", nodeResult.notReady)
		}
	}

	workloadResult := o.CheckWorkloads()
	if workloadResult.deployment.workloads == 0 {
		fmt.Fprintf(o.ErrOut, "Total deployment count is [%d]\n", workloadResult.deployment.workloads)
	} else {
		fmt.Fprintf(o.Out, "Total deployment count is [%d]\n", workloadResult.deployment.workloads)
		fmt.Fprintf(o.Out, "Total noLimit deployment count is [%d]\n", len(workloadResult.deployment.noLimitWorkload))
		fmt.Fprintf(o.Out, "Total oneReplica deployment count is [%d]\n", len(workloadResult.deployment.oneReplicaWorkload))

		if len(workloadResult.deployment.noLimitWorkload) > 0 {
			fmt.Fprintf(o.Out, "noLimit deployment: %v \n", workloadResult.deployment.noLimitWorkload)
		}
		if len(workloadResult.deployment.oneReplicaWorkload) > 0 {
			fmt.Fprintf(o.Out, "oneReplica deployment: %v \n", workloadResult.deployment.oneReplicaWorkload)
		}
	}
	if workloadResult.sts.workloads == 0 {
		fmt.Fprintf(o.ErrOut, "Total sts count is [%d]\n", workloadResult.sts.workloads)
	} else {
		fmt.Fprintf(o.Out, "Total sts count is [%d]\n", workloadResult.sts.workloads)
		fmt.Fprintf(o.Out, "Total noLimit sts count is [%d]\n", len(workloadResult.sts.noLimitWorkload))
		fmt.Fprintf(o.Out, "Total oneReplica sts count is [%d]\n", len(workloadResult.sts.oneReplicaWorkload))

		if len(workloadResult.sts.noLimitWorkload) > 0 {
			fmt.Fprintf(o.Out, "noLimit sts: %v \n", workloadResult.sts.noLimitWorkload)
		}
		if len(workloadResult.sts.oneReplicaWorkload) > 0 {
			fmt.Fprintf(o.Out, "oneReplica sts: %v \n", workloadResult.sts.oneReplicaWorkload)
		}
	}

	return nil
}
