package cmd

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type workloadResult struct {
	deployment workload
	sts        workload
}

type workload struct {
	workloads       int
	noLimitWorkload []string
	oneReplicaWorkload []string
}

//CheckWorkloads xx
func (o *InspectionOptions) CheckWorkloads() workloadResult {
	var (
		deploymentResult = workload{
			workloads:       0,
			noLimitWorkload: make([]string, 0),
		}
		stsResult = workload{
			workloads:       0,
			noLimitWorkload: make([]string, 0),
		}
	)

	deployments, err := o.kubeCli.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
	if errors.IsNotFound(err) {
		fmt.Fprintf(o.ErrOut, "%s", "deploy not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Fprintf(o.ErrOut, "Error getting deploy %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Fprintln(o.ErrOut, err.Error())
	} else {
		deploymentResult.workloads = len(deployments.Items)
		for i, _ := range deployments.Items {
			if *deployments.Items[i].Spec.Replicas == int32(1) {
				deploymentResult.oneReplicaWorkload = append(deploymentResult.oneReplicaWorkload,
					fmt.Sprintf("%s.%s", deployments.Items[i].ObjectMeta.Namespace, deployments.Items[i].ObjectMeta.Name))
			}
			for _, container := range deployments.Items[i].Spec.Template.Spec.Containers {
				if container.Resources.Limits == nil {
					deploymentResult.noLimitWorkload = append(deploymentResult.noLimitWorkload,
						fmt.Sprintf("%s.%s", deployments.Items[i].ObjectMeta.Namespace, deployments.Items[i].ObjectMeta.Name))
					break
				}
			}
		}
	}

	sts, err := o.kubeCli.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
	if errors.IsNotFound(err) {
		fmt.Fprintf(o.ErrOut, "%s", "sts not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Fprintf(o.ErrOut, "Error getting sts %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Fprintln(o.ErrOut, err.Error())
	} else {
		stsResult.workloads = len(sts.Items)
		for i, _ := range sts.Items {
			if *sts.Items[i].Spec.Replicas == int32(1) {
				stsResult.oneReplicaWorkload = append(stsResult.oneReplicaWorkload,
					fmt.Sprintf("%s.%s", sts.Items[i].ObjectMeta.Namespace, sts.Items[i].ObjectMeta.Name))
			}
			for _, container := range sts.Items[i].Spec.Template.Spec.Containers {
				if container.Resources.Limits == nil {
					stsResult.noLimitWorkload = append(deploymentResult.noLimitWorkload,
						fmt.Sprintf("%s.%s", sts.Items[i].ObjectMeta.Namespace, sts.Items[i].ObjectMeta.Name))
				}
			}
		}
	}

	return workloadResult{
		deployment: deploymentResult,
		sts:        stsResult,
	}
}
