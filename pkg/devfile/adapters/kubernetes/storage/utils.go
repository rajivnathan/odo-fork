package storage

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/openshift/odo/pkg/kclient"
	"github.com/openshift/odo/pkg/util"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/klog/glog"
)

const (
	size = "1Gi"
)

// Create adds storage to given component of given application
func Create(Client *kclient.Client, name, componentName string) (*corev1.PersistentVolumeClaim, error) {

	labels := map[string]string{
		"component":    componentName,
		"storage-name": name,
	}

	quantity, err := resource.ParseQuantity(size)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to parse size: %v", size)
	}

	randomChars := util.GenerateRandomString(4)
	namespaceKubernetesObject, err := util.NamespaceOpenShiftObject(name, componentName)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to create namespaced name")
	}
	namespaceKubernetesObject = fmt.Sprintf("%v-%v", namespaceKubernetesObject, randomChars)

	objectMeta := kclient.CreateObjectMeta(namespaceKubernetesObject, Client.Namespace, labels, nil)
	pvcSpec := kclient.GeneratePVCSpec(quantity)

	// Create PVC
	glog.V(3).Infof("Creating a PVC with name %v\n", namespaceKubernetesObject)
	pvc, err := Client.CreatePVC(objectMeta, *pvcSpec)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create PVC")
	}
	return pvc, nil
}
