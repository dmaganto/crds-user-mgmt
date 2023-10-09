package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func main() {
	// Use kubeconfig to create config
	//kubeconfig := "/your/path/to/kubeconfig"
	//config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//if err != nil {
	//	log.Fatal(err)
	//}
	config, err := rest.InClusterConfig()
	if err != nil {
		// handle error
	}
	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Define the CRD we want to watch
	crdGVR := schema.GroupVersionResource{
		Group:    "dmaganto.infra",
		Version:  "v1alpha1",
		Resource: "developers",
	}

	stopCh := make(chan struct{})
	signalCh := make(chan os.Signal, 2)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalCh
		close(stopCh)
	}()
	// For all namespaces
	//factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, corev1.NamespaceAll, nil)
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, "default", nil)
	informer := factory.ForResource(crdGVR)

	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)

			data, err := unstructuredObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(bytes.NewBuffer(data))
			//_, err = http.Post("http://tekton-pipeline-example", "application/json", bytes.NewBuffer(data))
			//if err != nil {
			//	log.Fatal(err)
			//}
		},
	})

	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	<-stopCh
}
