/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kubernetes

import (
	"testing"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stretchr/testify/assert"
)

func TestGetLastTransitionTimeOfUnreadyNode(t *testing.T) {
	start := time.Now()
	later := start.Add(10 * time.Minute)

	readyCondition := apiv1.NodeCondition{
		Type:               apiv1.NodeReady,
		Status:             apiv1.ConditionTrue,
		LastTransitionTime: metav1.NewTime(later),
	}
	unreadyCondition := apiv1.NodeCondition{
		Type:               apiv1.NodeReady,
		Status:             apiv1.ConditionFalse,
		LastTransitionTime: metav1.NewTime(later),
	}

	readyNode := &apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "readyNode",
			CreationTimestamp: metav1.NewTime(start),
		},
		Status: apiv1.NodeStatus{
			Capacity:    apiv1.ResourceList{},
			Allocatable: apiv1.ResourceList{},
			Conditions:  []apiv1.NodeCondition{readyCondition},
		},
	}
	unreadyNode := &apiv1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "unreadyNode",
			CreationTimestamp: metav1.NewTime(start),
		},
		Status: apiv1.NodeStatus{
			Capacity:    apiv1.ResourceList{},
			Allocatable: apiv1.ResourceList{},
			Conditions:  []apiv1.NodeCondition{unreadyCondition},
		},
	}

	lastTransitionTime, err := GetLastTransitionTimeOfUnreadyNode(readyNode)
	assert.Error(t, err)

	lastTransitionTime, err = GetLastTransitionTimeOfUnreadyNode(unreadyNode)
	assert.NotNil(t, lastTransitionTime)
	assert.Equal(t, later, lastTransitionTime)
	assert.Nil(t, err)
}
