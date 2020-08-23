package cmd

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type nodeResult struct {
	totalCount int
	notReady   []string
}

//CheckNodes xx
func (o *InspectionOptions) CheckNodes() nodeResult {
	var result = nodeResult{
		totalCount: 0,
		notReady:   make([]string, 0),
	}
	nodes, err := o.kubeCli.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if errors.IsNotFound(err) {
		fmt.Fprintf(o.ErrOut, "%s", "Node not found\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Fprintf(o.ErrOut, "Error getting node %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		fmt.Fprintln(o.ErrOut, err.Error())
	}
	result.totalCount = len(nodes.Items)
	if len(nodes.Items) > 0 {
		for i, _ := range nodes.Items {
			for _, status := range nodes.Items[i].Status.Conditions {
				if status.Type == "Ready" && status.Status != "True" {
					result.notReady = append(result.notReady, nodes.Items[i].ObjectMeta.Name)
				}
			}
		}
	}
	return result
}
