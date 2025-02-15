package v2alpha1

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestHandlerDeploymentList(t *testing.T) {
	t.Run("filter not ready deployment", func(t *testing.T) {
		list := &appsv1.DeploymentList{
			Items: []appsv1.Deployment{
				{
					Status: appsv1.DeploymentStatus{
						Replicas: 0,
					},
				},
				{
					Status: appsv1.DeploymentStatus{
						Replicas:      1,
						ReadyReplicas: 0,
					},
				},
				{
					Status: appsv1.DeploymentStatus{
						Replicas:      1,
						ReadyReplicas: 1,
					},
				},
			},
		}
		assert.Len(t, handlerDeploymentList(list), 1)
	})

	t.Run("sort deployment list", func(t *testing.T) {
		list := &appsv1.DeploymentList{
			Items: []appsv1.Deployment{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "emqx-1",
						CreationTimestamp: metav1.Time{Time: time.Now().AddDate(0, 0, 1)},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:      1,
						ReadyReplicas: 1,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:              "emqx-0",
						CreationTimestamp: metav1.Time{Time: time.Now().AddDate(0, 0, -1)},
					},
					Status: appsv1.DeploymentStatus{
						Replicas:      1,
						ReadyReplicas: 1,
					},
				},
			},
		}

		var l []string
		for _, d := range handlerDeploymentList(list) {
			l = append(l, d.DeepCopy().Name)
		}
		assert.ElementsMatch(t, []string{"emqx-0", "emqx-1"}, l)
	})
}

func TestHandlerEventList(t *testing.T) {
	t.Run("filter event", func(t *testing.T) {
		list := &corev1.EventList{
			Items: []corev1.Event{
				{
					Reason: "SuccessfulCreate",
				},
				{
					Reason: "ScalingReplicaSet",
				},
			},
		}
		assert.Len(t, handlerEventList(list), 1)
	})

	t.Run("sort event list", func(t *testing.T) {
		list := &corev1.EventList{
			Items: []corev1.Event{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "emqx-1",
					},
					LastTimestamp: metav1.Time{Time: time.Now().AddDate(0, 0, 1)},
					Reason:        "ScalingReplicaSet",
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "emqx-0",
					},
					LastTimestamp: metav1.Time{Time: time.Now().AddDate(0, 0, -1)},
					Reason:        "ScalingReplicaSet",
				},
			},
		}

		var l []string
		for _, e := range handlerEventList(list) {
			l = append(l, e.DeepCopy().Name)
		}
		assert.ElementsMatch(t, []string{"emqx-0", "emqx-1"}, l)
	})
}
