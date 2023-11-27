package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Record start time
	startTime := time.Now()
	// Use kubeconfig to create config
	kubeconfig := "/Users/PS10409/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Print("Loading in-cluster configuration\n")
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	// handle error
	// 	fmt.Println("Error loading in-cluster configuration")
	// 	log.Fatal(err)
	// }

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating dynamic client")
		log.Fatal(err)
	}
	// Define the CRD we want to watch
	tribesGVR := schema.GroupVersionResource{
		Group:    "usermgmt.infra",
		Version:  "v1alpha1",
		Resource: "tribes",
	}
	squadsGVR := schema.GroupVersionResource{
		Group:    "usermgmt.infra",
		Version:  "v1alpha1",
		Resource: "squads",
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
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dynamicClient, 0, "privileges-azdir", nil)
	tribeInformer := factory.ForResource(tribesGVR)
	squadInformer := factory.ForResource(squadsGVR)

	tribeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// Check for new squads
		AddFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)

			// Convert to metav1.Object to get CreationTimestamp
			tribe, err := meta.Accessor(unstructuredObj)
			if err != nil {
				log.Fatal(err)
			}
			// Only process tribes created after the informer started
			if tribe.GetCreationTimestamp().Time.After(startTime) {
				data, err := unstructuredObj.MarshalJSON()
				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("tribe %v has been created\n", unstructuredObj.GetName())
				fmt.Println(bytes.NewBuffer(data))
				//_, err = http.Post("http://tekton-pipeline-example", "application/json", bytes.NewBuffer(data))
				//if err != nil {
				//	log.Fatal(err)
				//}
			}
		},
		// Check for changes in squads
		UpdateFunc: func(oldObj, newObj interface{}) {
			unstructuredNewObj := newObj.(*unstructured.Unstructured)
			newData, err := unstructuredNewObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("Updated object:", bytes.NewBuffer(newData))

			unstructuredOldObj := oldObj.(*unstructured.Unstructured)
			oldData, err := unstructuredOldObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("Old object:", bytes.NewBuffer(oldData))

			// Checking specific fields of objects
			var dataMap map[string]interface{}
			err = json.Unmarshal(newData, &dataMap)
			if err != nil {
				log.Fatal(err)
			}

			var oldMap map[string]interface{}
			var newMap map[string]interface{}

			err = json.Unmarshal(oldData, &oldMap)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(newData, &newMap)
			if err != nil {
				log.Fatal(err)
			}

			oldSpec, oldExists := oldMap["spec"]
			newSpec, newExists := newMap["spec"]

			tribeName := newMap["metadata"].(map[string]interface{})["name"]
			if oldExists && newExists {
				oldSpecMap, okOld := oldSpec.(map[string]interface{})
				newSpecMap, okNew := newSpec.(map[string]interface{})

				if okOld && okNew {
					for key, oldValue := range oldSpecMap {
						newValue, exists := newSpecMap[key]
						if !exists || !reflect.DeepEqual(oldValue, newValue) {
							fmt.Printf("Tribe '%s' Field '%s' has been modified. Old value: %v, New value: %v\n", tribeName, key, oldValue, newValue)
						}
					}
				}
			}
		},
		// Check for deleted squads
		DeleteFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)
			data, err := unstructuredObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Deleted object:", bytes.NewBuffer(data))
			var dataMap map[string]interface{}
			err = json.Unmarshal(data, &dataMap)
			if err != nil {
				log.Fatal(err)
			}

			tribeName := dataMap["metadata"].(map[string]interface{})["name"]

			fmt.Printf("tribe %v has been deleted\n", tribeName)
		},
	})

	squadInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		// Check for new squads
		AddFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)
			// Convert to metav1.Object to get CreationTimestamp
			squad, err := meta.Accessor(unstructuredObj)
			if err != nil {
				log.Fatal(err)
			}
			// Only process squads created after the informer started
			if squad.GetCreationTimestamp().Time.After(startTime) {
				squadInfo, err := unstructuredObj.MarshalJSON()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(bytes.NewBuffer(squadInfo))
				fmt.Printf("Namespace %v\n", unstructuredObj.GetNamespace())
				name := unstructuredObj.GetName()
				fmt.Printf("squad %v has been created\n", name)

				// Getting specific fields from the squad. If we add new fields we will need to add them here :(
				applications := unstructuredObj.Object["spec"].(map[string]interface{})["applications"]
				approvers := unstructuredObj.Object["spec"].(map[string]interface{})["approvers"]
				developers := unstructuredObj.Object["spec"].(map[string]interface{})["developers"]
				devops := unstructuredObj.Object["spec"].(map[string]interface{})["devops"]
				federatedDevops := unstructuredObj.Object["spec"].(map[string]interface{})["federatedDevops"]

				// As everything is an interface we need to do a type assertion to get the actual value
				approversSlice, ok1 := approvers.([]interface{})
				developersSlice, ok2 := developers.([]interface{})
				devopsSlice, ok3 := devops.([]interface{})
				federatedDevopsSlice, ok4 := federatedDevops.([]interface{})
				if !ok1 || !ok2 || !ok3 || !ok4 {
					fmt.Println("Type assertion failed")
					return
				}

				var users []string
				for _, approver := range approversSlice {
					if approverStr, ok := approver.(string); ok {
						users = append(users, approverStr)
					}
				}

				for _, developer := range developersSlice {
					if developerStr, ok := developer.(string); ok {
						users = append(users, developerStr)
					}
				}
				for _, devop := range devopsSlice {
					if devopStr, ok := devop.(string); ok {
						users = append(users, devopStr)
					}
				}
				for _, federatedDevop := range federatedDevopsSlice {
					if federatedDevopStr, ok := federatedDevop.(string); ok {
						users = append(users, federatedDevopStr)
					}
				}
				// in users we have the list of users that need to be added to the slack group
				squadCreatedOrchestrator(name, users)
			}
		},
		// Check for changes in squads
		UpdateFunc: func(oldObj, newObj interface{}) {
			unstructuredNewObj := newObj.(*unstructured.Unstructured)
			newData, err := unstructuredNewObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}

			unstructuredOldObj := oldObj.(*unstructured.Unstructured)
			oldData, err := unstructuredOldObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}

			// Checking specific fields of objects
			var dataMap map[string]interface{}
			err = json.Unmarshal(newData, &dataMap)
			if err != nil {
				log.Fatal(err)
			}

			var oldMap map[string]interface{}
			var newMap map[string]interface{}

			err = json.Unmarshal(oldData, &oldMap)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(newData, &newMap)
			if err != nil {
				log.Fatal(err)
			}

			oldSpec, oldExists := oldMap["spec"]
			newSpec, newExists := newMap["spec"]

			squadName := newMap["metadata"].(map[string]interface{})["name"]
			if oldExists && newExists {
				oldSpecMap, okOld := oldSpec.(map[string]interface{})
				newSpecMap, okNew := newSpec.(map[string]interface{})

				if okOld && okNew {
					for key, oldValue := range oldSpecMap {
						newValue, exists := newSpecMap[key]
						if !exists || !reflect.DeepEqual(oldValue, newValue) {
							orchestrate(squadName, key, oldValue.([]interface{}), newValue.([]interface{}))
						}
					}
				}
			}
		},
		// Check for deleted squads
		DeleteFunc: func(obj interface{}) {
			unstructuredObj := obj.(*unstructured.Unstructured)
			data, err := unstructuredObj.MarshalJSON()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Deleted object:", bytes.NewBuffer(data))
			var dataMap map[string]interface{}
			err = json.Unmarshal(data, &dataMap)
			if err != nil {
				log.Fatal(err)
			}

			squadName := dataMap["metadata"].(map[string]interface{})["name"]

			fmt.Printf("Squad %v has been deleted\n", squadName)
		},
	})
	factory.Start(stopCh)
	factory.WaitForCacheSync(stopCh)

	<-stopCh
}

// function to compare slices
func compareSlices(oldSlice, newSlice []interface{}) ([]interface{}, []interface{}) {
	newElements := make([]interface{}, 0)
	deletedElements := make([]interface{}, 0)

	// Find new elements
	for _, newValue := range newSlice {
		exists := false
		for _, oldValue := range oldSlice {
			if oldValue == newValue {
				exists = true
				break
			}
		}
		if !exists {
			newElements = append(newElements, newValue)
		}
	}

	// Find deleted elements
	for _, oldValue := range oldSlice {
		exists := false
		for _, newValue := range newSlice {
			if oldValue == newValue {
				exists = true
				break
			}
		}
		if !exists {
			deletedElements = append(deletedElements, oldValue)
		}
	}

	return newElements, deletedElements
}

// function to orchestrate external calls based on diferences
func orchestrate(squadName interface{}, key string, oldValue, newValue []interface{}) {
	added, deleted := compareSlices(oldValue, newValue)
	if len(added) > 0 {
		fmt.Printf("Squad '%s' has new '%s': %v\n", squadName, key, added)
		switch key {
		case "applications":
			// Create new service in opsgenie
			fmt.Printf("Creating new services %v in opsgenie for squad %s\n", added, squadName)
			for _, service := range added {
				teamId, err := getTeamIdFromName(squadName.(string))
				if err != nil {
					fmt.Println(err)
				}
				createService(service.(string), teamId)
			}
		case "federatedDevops":
			// Create new user in opsgenie
			fmt.Printf("Creating new user %v in opsgenie for squad %s\n", added, squadName)
			fmt.Printf("Adding new member to slack groups %v for squad %s\n", added, squadName)
		case "developers":
			// Add the user to slack groups
			fmt.Printf("Adding new member to slack groups %v for squad %s\n", added, squadName)
			for _, user := range added {
				//addUserToSlackGroup(user.(string), squadName.(string))
				err := addUserToSlackGroup(user.(string), squadName.(string))
				if err != nil {
					log.Fatal(err)
				}
			}
		default:
			//do nothing
		}
	}
	if len(deleted) > 0 {
		fmt.Printf("Squad '%s' has less '%s': %v\n", squadName, key, deleted)
		switch key {
		case "applications":
			// Delete service in opsgenie
			fmt.Printf("Deleting services %v in opsgenie for squad %s\n", deleted, squadName)
		case "federatedDevops":
			// Delete user in opsgenie
			fmt.Printf("Deleting user %v in opsgenie for squad %s\n", deleted, squadName)
			fmt.Printf("Deleting member to slack groups %v for squad %s\n", deleted, squadName)
		case "developers":
			// Add the user to slack groups
			fmt.Printf("Deleting member to slack groups %v for squad %s\n", deleted, squadName)
		default:
			//do nothing
		}
	}

}

// function to handle operations when a squad is created
func squadCreatedOrchestrator(squadName string, users []string) {
	// Create slack channels
	fmt.Printf("Creating slack channel %s\n", "monitoring-"+squadName+"-non-critical")
	fmt.Printf("Creating slack channel %s\n", "monitoring-"+squadName+"-critical")
	channelID1, err1 := createSlackChannel("monitoring-" + squadName + "-non-critical")
	channelID2, err2 := createSlackChannel("monitoring-" + squadName + "-critical")
	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	}
	fmt.Printf("Slack channel: %s successfully created\n", "monitoring-"+squadName+"-non-critical")
	fmt.Printf("Slack channel: %s successfully created\n", "monitoring-"+squadName+"-critical")

	var channelIDs []string = []string{channelID1, channelID2}
	var userIDs []string
	// getting the userIDs in slack from the users
	for _, user := range users {
		userID, err := getUserByEmail(user)
		if err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, userID)
	}

	fmt.Printf("Creating slack group %s\n", "azdir-testingautomation-"+squadName)
	groupID, err := createNewUserGroup("azdir-testingautomation-"+squadName, userIDs, channelIDs)
	if err != nil {
		log.Fatal(err)

	}
	fmt.Printf("Slack group: %s successfully created. groupID: %s\n", "azdir-testingautomation"+squadName, groupID)
}
