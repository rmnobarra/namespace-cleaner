package main

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestCleanupNamespaces(t *testing.T) {
	testCases := []struct {
		name           string
		namespaceName  string
		regexStr       string
		expectedErrMsg string
	}{
		// Test cases here...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create the NamespaceCleaner with the regex pattern
			cleaner, err := NewNamespaceCleaner(tc.regexStr)
			if err != nil {
				t.Fatalf("error creating namespace cleaner: %v", err)
			}

			// Create a fake Kubernetes clientset
			clientSet := fake.NewSimpleClientset()

			// Add the namespace to the fake clientset
			clientSet.CoreV1().Namespaces().Create(context.Background(), &metav1.Namespace{
				ObjectMeta: metav1.ObjectMeta{Name: tc.namespaceName},
			}, metav1.CreateOptions{})

			// Run the CleanupNamespaces function
			err = cleaner.CleanupNamespaces(context.Background(), clientSet)
			if tc.expectedErrMsg != "" {
				if err == nil || err.Error() != tc.expectedErrMsg {
					t.Errorf("expected error: %v, got: %v", tc.expectedErrMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}

			// Check if the namespace was deleted or not
			_, err = clientSet.CoreV1().Namespaces().Get(context.Background(), tc.namespaceName, metav1.GetOptions{})
			if tc.expectedErrMsg == "" {
				if err == nil {
					t.Errorf("expected namespace %s to be deleted, but it still exists", tc.namespaceName)
				}
			} else {
				if err != nil {
					t.Errorf("expected namespace %s to exist, but it was deleted: %v", tc.namespaceName, err)
				}
			}
		})
	}
}
