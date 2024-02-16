package main

import (
	"context"
	"encoding/base64"
	"github.com/stretchr/testify/assert"
	"go-lab/common"
	apps "k8s.io/api/apps/v1"
	core "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"testing"
)

var kube *kubernetes.Clientset

var values *common.Values

func TestMain(m *testing.M) {
	var err error
	if values, err = common.LoadValues("../config.yml"); err != nil {
		panic(err.Error())
	}
	cadata, _ := base64.StdEncoding.DecodeString(values.KUBERNETES.CAData)
	certdata, _ := base64.StdEncoding.DecodeString(values.KUBERNETES.CertData)
	keydata, _ := base64.StdEncoding.DecodeString(values.KUBERNETES.KeyData)
	kube, err = kubernetes.NewForConfig(&rest.Config{
		Host: values.KUBERNETES.Host,
		TLSClientConfig: rest.TLSClientConfig{
			CAData:   cadata,
			CertData: certdata,
			KeyData:  keydata,
		},
	})
	if err != nil {
		panic(err.Error())
	}
	os.Exit(m.Run())
}

func TestServerVersion(t *testing.T) {
	info, err := kube.ServerVersion()
	assert.NoError(t, err)
	t.Log(info)
	t.Log(info.GoVersion)
	t.Log(info.Major)
	t.Log(info.Minor)
	t.Log(info.Compiler)
	t.Log(info.Platform)
}

func TestNodes(t *testing.T) {
	nodes, err := kube.CoreV1().Nodes().
		List(context.TODO(), meta.ListOptions{})
	assert.NoError(t, err)
	for _, node := range nodes.Items {
		t.Log(node.Status)
	}
}

func TestDeploymentsList(t *testing.T) {
	data, err := kube.AppsV1().Deployments("kube-system").List(context.TODO(), meta.ListOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestDeploymentsCreate(t *testing.T) {
	ctx := context.TODO()
	container := core.Container{
		Name:            "schedule",
		Image:           "ccr.ccs.tencentyun.com/weplanx/schedule:v1.3.4",
		ImagePullPolicy: core.PullAlways,
		Env: []core.EnvVar{
			{Name: "MODE", Value: "release"},
			{Name: "NAMESPACE", Value: "example"},
			{Name: "NODE", Value: "xxx"},
			{Name: "NATS_HOSTS", Value: values.NATS.Url},
			{Name: "NATS_NKEY", Value: values.NATS.NKey},
		},
	}
	deployment := &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name: "schedule",
		},
		Spec: apps.DeploymentSpec{
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{"app": "schedule"},
			},
			Template: core.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{"app": "schedule"},
				},
				Spec: core.PodSpec{
					Containers: []core.Container{container},
				},
			},
		},
	}

	r, err := kube.AppsV1().
		Deployments("default").
		Create(ctx, deployment, meta.CreateOptions{})
	assert.NoError(t, err)
	t.Log(r)
}

func TestDeploymentsDelete(t *testing.T) {
	ctx := context.TODO()
	err := kube.AppsV1().Deployments("default").
		Delete(ctx, "schedule", meta.DeleteOptions{})
	assert.NoError(t, err)
}

func TestConfigMapCreate(t *testing.T) {
	configMap := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "kube-system",
			Name:      "test",
		},
		Data: map[string]string{
			"ini": "engine = On",
		},
	}

	data, err := kube.CoreV1().
		ConfigMaps("kube-system").
		Create(context.TODO(), configMap, meta.CreateOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestConfigMapUpdate(t *testing.T) {
	configMap := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "kube-system",
			Name:      "test",
		},
		Data: map[string]string{
			"ini": "engine = On",
			"zx":  "cccc",
		},
	}
	data, err := kube.CoreV1().
		ConfigMaps("kube-system").
		Update(context.TODO(), configMap, meta.UpdateOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestConfigMapDelete(t *testing.T) {
	err := kube.CoreV1().
		ConfigMaps("kube-system").
		Delete(context.TODO(), "test", meta.DeleteOptions{})
	assert.NoError(t, err)
}

func TestIngressList(t *testing.T) {
	data, err := kube.NetworkingV1().
		Ingresses("kube-system").
		List(context.TODO(), meta.ListOptions{})
	assert.NoError(t, err)
	for _, x := range data.Items {
		t.Log(x.Name)
	}
}

func TestIngressCreate(t *testing.T) {
	prefix := networking.PathTypePrefix
	ingress := &networking.Ingress{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "kube-system",
			Name:      "test",
			Annotations: map[string]string{
				"traefik.ingress.kubernetes.io/router.entrypoints":        "web,websecure",
				"traefik.ingress.kubernetes.io/router.tls.certresolver":   "kainonly",
				"traefik.ingress.kubernetes.io/router.tls.domains.0.main": "kainonly.com",
				"traefik.ingress.kubernetes.io/router.tls.domains.0.sans": "*.kainonly.com",
			},
		},
		Spec: networking.IngressSpec{
			Rules: []networking.IngressRule{
				{
					Host: "test.hnvane.com",
					IngressRuleValue: networking.IngressRuleValue{
						HTTP: &networking.HTTPIngressRuleValue{
							Paths: []networking.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &prefix,
									Backend: networking.IngressBackend{
										Service: &networking.IngressServiceBackend{
											Name: "nginx",
											Port: networking.ServiceBackendPort{
												Number: 8080,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	data, err := kube.NetworkingV1().
		Ingresses("kube-system").
		Create(context.TODO(), ingress, meta.CreateOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestIngressDelete(t *testing.T) {
	err := kube.NetworkingV1().
		Ingresses("kube-system").
		Delete(context.TODO(), "test", meta.DeleteOptions{})
	assert.NoError(t, err)
}

func TestSecret(t *testing.T) {
	ctx := context.TODO()
	b, err := os.ReadFile("../nats-accounts/accounts.conf")
	assert.NoError(t, err)
	secret := &core.Secret{
		ObjectMeta: meta.ObjectMeta{
			Namespace: "default",
			Name:      "includx",
		},
		Data: map[string][]byte{
			"data": b,
		},
		Type: "Opaque",
	}

	_, err = kube.CoreV1().
		Secrets("default").
		Update(ctx, secret, meta.UpdateOptions{})
	assert.NoError(t, err)
}
