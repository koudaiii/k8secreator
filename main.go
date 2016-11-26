package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	namespace  = flag.String("namespace", DefaultNamespace(), "namespace")
	name       = flag.String("name", "dotenv", "namespace")
)

func DefaultNamespace() string {
	var basename string
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_, err = os.Stat(pwd + "/kubernetes")
	if err == nil {
		dirs := strings.Split(pwd, "/")
		basename = strings.Join(dirs[len(dirs)-1:], "")
	} else {
		basename = "default"
	}
	return basename
}

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	if *kubeconfig == "" {
		*kubeconfig = clientcmd.RecommendedHomeFile
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// creates the clientset

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// create secret dataset

	s := &v1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      *name,
			Namespace: *namespace,
		},
		StringData: map[string]string{
			"key": "value",
		},
		Type: "Opaque",
	}

	// create secret

	secret := clientset.Core().Secrets(*namespace)
	result, err := secret.Create(s)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// display result

	// fmt.Printf("Secret %s in the cluster\n", result)
	fmt.Printf("Name: %s Namespace: %s\n", result.ObjectMeta.Name, result.ObjectMeta.Namespace)
	var v string
	for key, value := range result.Data {
		if false {
			v = base64.StdEncoding.EncodeToString(value)
		} else {
			v = strconv.Quote(string(value))
		}
		fmt.Printf("Key: %s Value: %s\n", key, v)
	}
}
