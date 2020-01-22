package component

import (
	"os"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/redhat-developer/odo-fork/pkg/config"
	"github.com/redhat-developer/odo-fork/pkg/devfile"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	"github.com/redhat-developer/odo-fork/pkg/log"
	"github.com/redhat-developer/odo-fork/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TaskExecDevfile is the Build Task or the Runtime Task execution implementation of the IDP
func TaskExecDevfile(Client *kclient.Client, componentConfig config.LocalConfigInfo, fullBuild bool, devfile *devfile.Devfile) error {
	namespace := Client.Namespace
	cmpName := componentConfig.GetName()
	appName := componentConfig.GetApplication()

	// Namespace the component
	namespacedKubernetesObject, err := util.NamespaceKubernetesObject(cmpName, appName)

	glog.V(0).Infof("Namespace: %s\n", namespace)

	serviceAccountName := "default"
	glog.V(0).Infof("Service Account: %s\n", serviceAccountName)

	// cwd is the project root dir, where udo command will run
	cwd, err := os.Getwd()
	if err != nil {
		err = errors.New("Unable to get the cwd" + err.Error())
		return err
	}
	glog.V(0).Infof("CWD: %s\n", cwd)

	timeout := int64(10)

	// Check if the component exists, otherwise create one
	glog.V(0).Infof("Checking if the Component has already been deployed...\n")

	foundComponent := false
	watchOptions := metav1.ListOptions{
		LabelSelector:  "app=" + namespacedKubernetesObject + ",deployment=" + namespacedKubernetesObject,
		TimeoutSeconds: &timeout,
	}
	po, _ := Client.WaitAndGetPod(watchOptions, corev1.PodRunning, "Checking to see if a Component has already been deployed")
	if po != nil {
		glog.V(0).Infof("Running pod found: %s...\n\n", po.Name)
		// BuildTaskInstance.PodName = po.Name
		foundComponent = true
	}

	// Create component if it doesn't exist
	if !foundComponent {
		glog.V(0).Info("===============================")
		glog.V(0).Info("Creating the Component")

		s := log.Spinner("Creating component")
		defer s.End(false)
		// if err = BuildTaskInstance.CreateComponent(Client, componentConfig, cmpPVC); err != nil {
		// 	err = errors.New("Unable to create component deployment: " + err.Error())
		// 	return err
		// }
		labels := map[string]string{
			"app":        namespacedKubernetesObject,
			"deployment": namespacedKubernetesObject,
		}
		if po, err = createComponentFromDevfile(devfile, namespacedKubernetesObject, namespace, serviceAccountName, labels); err != nil {
			err = errors.New("Unable to create component deployment: " + err.Error())
			return err
		}

		if _, err = Client.CreateDeployment(po); err != nil {
			err = errors.New("Unable to create component deployment: " + err.Error())
			return err
		}
		s.End(true)

		glog.V(0).Info("Successfully created the component")
		glog.V(0).Info("===============================")
	}

	// Execute task on component
	// runActions(Client, actions, po)

	return nil
}

// Create the component based on all containers referenced in the IDP, we will have a single fat pod with all containers
func createComponentFromDevfile(devfile *devfile.Devfile, componentName, namespace, serviceAccount string, labels map[string]string) (*corev1.Pod, error) {

	// Get all the possible tasks that can be run
	// var tasks []idp.SpecTask = devPack.Spec.Tasks
	// containerNames := make(map[string]struct{})
	// var exists = struct{}{}

	// // Find the container for each task and add it to the set
	// for _, task := range tasks {
	// 	if _, ok := containerNames[task.Container]; !ok {
	// 		containerNames[task.Container] = exists
	// 	}
	// }

	// Get a container reference for each container in the set
	containers := []corev1.Container{}

	for _, component := range devfile.Components {
		if component.Type == "dockerimage" && component.Alias != nil {
			glog.V(0).Info("Component image: ", component.Image)
			k8container := kclient.GenerateContainerSpec(*component.Alias, *component.Image, true)
			containers = append(containers, k8container)
		}
	}

	if len(containers) == 0 {
		return nil, errors.New("No containers defined")
	}

	// Create a pod that includes all of the containers
	po, err := kclient.GeneratePodSpec("fatpod", namespace, serviceAccount, labels, containers, []string{}, []string{}, []string{})

	return po, err

}
