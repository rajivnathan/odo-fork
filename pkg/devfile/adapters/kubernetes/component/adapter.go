package component

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"

	"github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/openshift/odo/pkg/devfile/adapters/kubernetes/utils"
	"github.com/openshift/odo/pkg/kclient"
)

// New instantiantes a component adapter
func New(adapterContext common.AdapterContext, client kclient.Client) Adapter {
	return Adapter{
		Client:         client,
		AdapterContext: adapterContext,
	}
}

// Adapter is a component adapter implementation for Kubernetes
type Adapter struct {
	Client kclient.Client
	common.AdapterContext
}

// Initialize initilizes the component from the devfile
func (a Adapter) Initialize() (*corev1.PodTemplateSpec, error) {
	componentName := a.ComponentName

	labels := map[string]string{
		"component": componentName,
	}

	containers := utils.GetContainers(a.Devfile)
	if len(containers) == 0 {
		return nil, fmt.Errorf("No valid components found in the devfile")
	}

	objectMeta := kclient.CreateObjectMeta(componentName, a.Client.Namespace, labels, nil)
	podTemplateSpec := kclient.GeneratePodTemplateSpec(objectMeta, containers)
	return podTemplateSpec, nil
}

// Start updates the component if a matching component exists or creates one if it doesn't exist
func (a Adapter) Start(podTemplateSpec *corev1.PodTemplateSpec) (err error) {
	componentName := a.ComponentName

	// 	componentAliasToVolumes := utils.GetVolumes(a.Devfile)

	// 	// Get a list of all the unique volume names
	// 	var uniqueVolumes []string
	// 	processedVolumes := make(map[string]bool)
	// 	for _, volumes := range componentAliasToVolumes {
	// 		for _, vol := range volumes {
	// 			if _, ok := processedVolumes[*vol.Name]; !ok {
	// 				processedVolumes[*vol.Name] = true
	// 				uniqueVolumes = append(uniqueVolumes, *vol.Name)
	// 			}
	// 		}
	// 	}

	// 	// createComponentStorage creates PVC from the unique Devfile volumes and returns a map of volume name to the PVC created
	// 	volumeNameToPVC, err := utils.CreateComponentStorage(&a.Client, uniqueVolumes, componentName)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Add PVC to the podTemplateSpec
	// 	err = kclient.AddPVCAndVolumeMount(podTemplateSpec, volumeNameToPVC, componentAliasToVolumes)
	// 	if err != nil {
	// 		return err
	// 	}

	deploymentSpec := kclient.GenerateDeploymentSpec(*podTemplateSpec)

	glog.V(3).Infof("Creating deployment %v", deploymentSpec.Template.GetName())
	glog.V(3).Infof("The component name is %v", componentName)

	if utils.ComponentExists(a.Client, componentName) {
		glog.V(3).Info("The component already exists, attempting to update it")
		_, err = a.Client.UpdateDeployment(*deploymentSpec)
		if err != nil {
			return err
		}
		glog.V(3).Infof("Successfully updated component %v", componentName)
	} else {
		_, err = a.Client.CreateDeployment(*deploymentSpec)
		if err != nil {
			return err
		}
		glog.V(3).Infof("Successfully created component %v", componentName)
	}

	podSelector := fmt.Sprintf("component=%s", componentName)
	watchOptions := metav1.ListOptions{
		LabelSelector: podSelector,
	}

	_, err = a.Client.WaitAndGetPod(watchOptions, corev1.PodRunning, "Waiting for component to start")
	return err
}
