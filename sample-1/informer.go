package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"shawnrong.github.com/k8s-crd-demo/sample1/api/types/v1alpha1"
	clientV1alpha1 "shawnrong.github.com/k8s-crd-demo/sample1/clientset/v1alpha1"
	"time"
)

func WatchResources(clientSet clientV1alpha1.ExampleV1Alpha1Interface) cache.Store {
	projectStore, projectController := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(lo metav1.ListOptions) (result runtime.Object, err error) {
				return clientSet.Projects("some-namespace").List(lo)
			},
			WatchFunc: func(lo metav1.ListOptions) (watch.Interface, error) {
				return clientSet.Projects("some-namespace").Watch(lo)
			},
		},
		&v1alpha1.Project{},
		1*time.Minute,
		cache.ResourceEventHandlerFuncs{},
	)

	go projectController.Run(wait.NeverStop)
	return projectStore
}
