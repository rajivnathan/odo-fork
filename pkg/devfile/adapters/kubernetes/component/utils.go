package component

import (
	"errors"

	"github.com/golang/glog"

	"github.com/openshift/odo/pkg/devfile/adapters/kubernetes/storage"
	devfileCommon "github.com/openshift/odo/pkg/devfile/versions/common"
	"github.com/openshift/odo/pkg/kclient"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func convertEnvs(vars []devfileCommon.DockerimageEnv) []corev1.EnvVar {
	kVars := []corev1.EnvVar{}
	for _, env := range vars {
		kVars = append(kVars, corev1.EnvVar{
			Name:  *env.Name,
			Value: *env.Value,
		})
	}
	return kVars
}

func getResourceReqs(comp devfileCommon.DevfileComponent) corev1.ResourceRequirements {
	reqs := corev1.ResourceRequirements{}
	limits := make(corev1.ResourceList)
	if comp.MemoryLimit != nil {
		memoryLimit, err := resource.ParseQuantity(*comp.MemoryLimit)
		if err == nil {
			limits[corev1.ResourceMemory] = memoryLimit
		}
		reqs.Limits = limits
	}
	return reqs
}

func componentExists(client kclient.Client, name string) bool {
	_, err := client.GetDeploymentByName(name)
	return err == nil
}

// createComponentStorage creates PVCs with the given list of volume names
func createComponentStorage(Client *kclient.Client, volumes []string, componentName string) (map[string]*corev1.PersistentVolumeClaim, error) {
	volumeNameToPVC := make(map[string]*corev1.PersistentVolumeClaim)

	for _, vol := range volumes {
		label := "component=" + componentName + ",storage-name=" + vol
		glog.V(3).Infof("Checking for PVC with name %v and label %v\n", vol, label)
		PVCs, err := Client.GetPVCsFromSelector(label)
		if err != nil {
			glog.V(0).Infof("Error occured while getting the PVC")
			err = errors.New("Unable to get the PVC: " + err.Error())
			return nil, err
		}
		if len(PVCs) == 1 {
			glog.V(3).Infof("Found an existing PVC with name %v and label %v\n", vol, label)
			existingPVC := &PVCs[0]
			volumeNameToPVC[vol] = existingPVC
		} else if len(PVCs) == 0 {
			glog.V(3).Infof("Creating a PVC with name %v and label %v\n", vol, label)
			createdPVC, err := storage.Create(Client, vol, componentName)
			volumeNameToPVC[vol] = createdPVC
			if err != nil {
				glog.V(0).Infof("Error creating the PVC: " + err.Error())
				err = errors.New("Error creating the PVC: " + err.Error())
				return nil, err
			}
		} else {
			err = errors.New("More than 1 PVC found with the label " + label + ": " + err.Error())
			return nil, err
		}
	}

	return volumeNameToPVC, nil
}
