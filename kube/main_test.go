package main

import (
	"context"
	"flag"
	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
	"testing"
)

var clientset *kubernetes.Clientset

func TestMain(m *testing.M) {
	var err error
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	os.Exit(m.Run())
}

func TestNodes(t *testing.T) {
	data, err := clientset.CoreV1().Nodes().List(context.TODO(), meta.ListOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestDeployments(t *testing.T) {
	data, err := clientset.AppsV1().Deployments("kube-system").List(context.TODO(), meta.ListOptions{})
	assert.NoError(t, err)
	t.Log(data)
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

	data, err := clientset.CoreV1().
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
	data, err := clientset.CoreV1().
		ConfigMaps("kube-system").Update(context.TODO(), configMap, meta.UpdateOptions{})
	assert.NoError(t, err)
	t.Log(data)
}

func TestConfigMapDelete(t *testing.T) {
	err := clientset.CoreV1().
		ConfigMaps("kube-system").
		Delete(context.TODO(), "test", meta.DeleteOptions{})
	assert.NoError(t, err)
}
