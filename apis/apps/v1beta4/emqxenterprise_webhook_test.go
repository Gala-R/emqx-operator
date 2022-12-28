/*
Copyright 2021.

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

package v1beta4

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestEnterpriseDefault(t *testing.T) {
	instance := &EmqxEnterprise{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "webhook-test",
			Namespace: "default",
			Labels: map[string]string{
				"foo": "bar",
			},
			Annotations: map[string]string{
				"foo": "bar",
			},
		},
		Spec: EmqxEnterpriseSpec{
			Template: EmqxTemplate{
				Spec: EmqxTemplateSpec{
					EmqxContainer: EmqxContainer{
						Name: "emqx",
					},
				},
			},
		},
	}
	instance.Default()

	t.Run("default labels", func(t *testing.T) {
		assert.Equal(t, map[string]string{
			"foo":                     "bar",
			"apps.emqx.io/managed-by": "emqx-operator",
			"apps.emqx.io/instance":   "webhook-test",
		}, instance.Labels)

		assert.Equal(t, map[string]string{
			"foo":                     "bar",
			"apps.emqx.io/managed-by": "emqx-operator",
			"apps.emqx.io/instance":   "webhook-test",
		}, instance.Spec.Template.Labels)
	})

	t.Run("default emqx acl", func(t *testing.T) {
		assert.ElementsMatch(t, []string{
			`{allow, {user, "dashboard"}, subscribe, ["$SYS/#"]}.`,
			`{allow, {ipaddr, "127.0.0.1"}, pubsub, ["$SYS/#", "#"]}.`,
			`{deny, all, subscribe, ["$SYS/#", {eq, "#"}]}.`,
			`{allow, all}.`,
		}, instance.Spec.Template.Spec.EmqxContainer.EmqxACL)
	})

	t.Run("default emqx config", func(t *testing.T) {
		assert.Equal(t, map[string]string{
			"name":                  "webhook-test",
			"log.to":                "console",
			"cluster.discovery":     "dns",
			"cluster.dns.type":      "srv",
			"cluster.dns.app":       "webhook-test",
			"cluster.dns.name":      "webhook-test-headless.default.svc.cluster.local",
			"listener.tcp.internal": "",
		}, instance.Spec.Template.Spec.EmqxContainer.EmqxConfig)
	})

	t.Run("default service template", func(t *testing.T) {
		assert.Equal(t, ServiceTemplate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "webhook-test",
				Namespace: "default",
				Labels: map[string]string{
					"foo":                     "bar",
					"apps.emqx.io/managed-by": "emqx-operator",
					"apps.emqx.io/instance":   "webhook-test",
				},
				Annotations: map[string]string{
					"foo": "bar",
				},
			},
			Spec: corev1.ServiceSpec{
				Selector: map[string]string{
					"foo":                     "bar",
					"apps.emqx.io/managed-by": "emqx-operator",
					"apps.emqx.io/instance":   "webhook-test",
				},
				Ports: []corev1.ServicePort{
					{
						Name:       "http-management-8081",
						Port:       8081,
						Protocol:   corev1.ProtocolTCP,
						TargetPort: intstr.FromInt(8081),
					},
				},
			},
		}, instance.Spec.ServiceTemplate)
	})
}

func TestEnterpriseValidateUpdate(t *testing.T) {
	instance := &EmqxEnterprise{}

	assert.Nil(t, instance.ValidateUpdate(&EmqxEnterprise{}))
	assert.Error(t, instance.ValidateUpdate(&EmqxEnterprise{
		Spec: EmqxEnterpriseSpec{
			Persistent: &corev1.PersistentVolumeClaimTemplate{
				Spec: corev1.PersistentVolumeClaimSpec{
					StorageClassName: &[]string{"fake"}[0],
				},
			},
		},
	}))
}
