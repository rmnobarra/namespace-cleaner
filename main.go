package main

import (
	"context"
	"fmt"
	"os"
	"regexp"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type NamespaceCleaner struct {
	clientSet *kubernetes.Clientset
	regex     *regexp.Regexp
}

func NewNamespaceCleaner(regexStr string) (*NamespaceCleaner, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("error building kubernetes in-cluster config: %v", err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error creating kubernetes clientset: %v", err)
	}

	regex, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex pattern: %v", err)
	}

	return &NamespaceCleaner{
		clientSet: clientSet,
		regex:     regex,
	}, nil
}

func (n *NamespaceCleaner) CleanupNamespaces(ctx context.Context) error {
	namespaces, err := n.clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("error getting namespaces list: %v", err)
	}

	for _, namespace := range namespaces.Items {
		if n.regex.MatchString(namespace.Name) {
			if err := n.clientSet.CoreV1().Namespaces().Delete(ctx, namespace.Name, metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("error deleting namespace %s: %v", namespace.Name, err)
			}
		}
	}
	return nil
}

func main() {
	regexStr := os.Getenv("NAMESPACE_SELECTOR")
	if regexStr == "" {
		fmt.Println("NAMESPACE_SELECTOR environment variable not set")
		return
	}

	cleaner, err := NewNamespaceCleaner(regexStr)
	if err != nil {
		fmt.Printf("error creating namespace cleaner: %v\n", err)
		return
	}

	ctx := context.Background()

	if err := cleaner.CleanupNamespaces(ctx); err != nil {
		fmt.Printf("error cleaning up namespaces: %v\n", err)
	}
}
