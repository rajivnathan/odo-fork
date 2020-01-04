package component

import (
	"strings"

	"github.com/redhat-developer/odo-fork/pkg/idp"
	"github.com/redhat-developer/odo-fork/pkg/kclient"
	corev1 "k8s.io/api/core/v1"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

// Task contains information required to run a command on a container
type Task struct {
	Name      string
	Type      string
	Container Container
	Command   []string
}

type Container struct {
	Name           string
	Image          string
	VolumeMappings []idp.VolumeMapping
	RuntimePorts   idp.RuntimePorts
}

func runTasks(Client *kclient.Client, tasks []idp.SpecTask, pod *corev1.Pod) (err error) {

	for _, task := range tasks {

		if len(task.Command) > 0 {
			err = executeCmd(Client, task.Command, pod)
			if err != nil {
				glog.V(0).Infof("Error occured while executing command %s in the pod %s: %s\n", strings.Join(task.Command, " "), pod.Name, err)
				err = errors.New("Unable to exec command " + strings.Join(task.Command, " ") + " in the runtime container: " + err.Error())
				return
			}
		}
	}

	// noTimeout := int64(0)

	// taskContainerInfo, err := devPack.GetTaskContainerInfo(task)
	// if err != nil {
	// 	glog.V(0).Infof("Error occured while getting the Task Container Info for task " + task.Name)
	// 	err = errors.New("Error occured while getting the Task Container Info for task " + task.Name + ": " + err.Error())
	// 	return err
	// }

	// containerImage := taskContainerInfo.Image
	// var containerName, trimmedNamespacedKubernetesObject, srcDestination string

	// if len(namespacedKubernetesObject) > 40 {
	// 	trimmedNamespacedKubernetesObject = namespacedKubernetesObject[:40]
	// } else {
	// 	trimmedNamespacedKubernetesObject = namespacedKubernetesObject
	// }

	// if task.Type == idp.RuntimeTask {
	// 	containerName = trimmedNamespacedKubernetesObject + "-runtime"
	// } else if task.Type == idp.SharedTask {
	// 	containerName = trimmedNamespacedKubernetesObject + task.Container
	// 	if len(containerName) > 63 {
	// 		containerName = containerName[:63]
	// 	}
	// }

	// var pvcClaimName, mountPath, subPath []string
	// var cmpPVC []*corev1.PersistentVolumeClaim
	// for _, vm := range taskContainerInfo.VolumeMappings {
	// 	cmpPVC = append(cmpPVC, idpPVC[vm.VolumeName])
	// 	pvcClaimName = append(pvcClaimName, idpPVC[vm.VolumeName].Name)
	// 	mountPath = append(mountPath, vm.ContainerPath)
	// 	subPath = append(subPath, vm.SubPath)
	// }

	// if len(task.SourceMapping.DestPath) > 0 {
	// 	srcDestination = task.SourceMapping.DestPath
	// }

	// BuildTaskInstance := BuildTask{
	// 	// Name:           containerName,
	// 	Image:         task.Container.Image,
	// 	ContainerName: task.Container.Name,
	// 	// Namespace:      namespace,
	// 	MountPath:      mountPath,
	// 	SubPath:        subPath,
	// 	Command:        task.Command,
	// 	SrcDestination: srcDestination,
	// }
	// BuildTaskInstance.Labels = map[string]string{
	// 	"app": BuildTaskInstance.Name,
	// }

	// var watchOptions metav1.ListOptions
	// if task.Type == idp.SharedTask {
	// 	watchOptions = metav1.ListOptions{
	// 		LabelSelector:  "app=" + BuildTaskInstance.Name,
	// 		TimeoutSeconds: &timeout,
	// 	}
	// } else if task.Type == idp.RuntimeTask {
	// 	watchOptions = metav1.ListOptions{
	// 		LabelSelector:  "app=" + namespacedKubernetesObject + ",deployment=" + namespacedKubernetesObject,
	// 		TimeoutSeconds: &timeout,
	// 	}
	// 	BuildTaskInstance.Ports = runtimePorts
	// }

	// glog.V(0).Infof("Checking if " + task.Type + " Container has already been deployed...\n")

	// foundTaskContainer := false
	// po, _ := Client.WaitAndGetPod(watchOptions, corev1.PodRunning, "Checking to see if a "+task.Type+" Container has already been deployed")
	// if po != nil {
	// 	glog.V(0).Infof("Running pod found: %s...\n\n", po.Name)
	// 	BuildTaskInstance.PodName = po.Name
	// 	foundTaskContainer = true
	// }

	// if !foundTaskContainer {
	// 	glog.V(0).Info("===============================")
	// 	glog.V(0).Info("Creating a " + task.Type + " Container")

	// 	if task.Type == idp.SharedTask {
	// 		s := log.Spinner("Creating pod")
	// 		defer s.End(false)
	// 		_, err := Client.CreatePod(BuildTaskInstance.Name, BuildTaskInstance.ContainerName, BuildTaskInstance.Image, BuildTaskInstance.ServiceAccountName, BuildTaskInstance.Labels, BuildTaskInstance.PVCName, BuildTaskInstance.MountPath, BuildTaskInstance.SubPath, BuildTaskInstance.Privileged)
	// 		if err != nil {
	// 			glog.V(0).Info("Failed to create a pod: " + err.Error())
	// 			err = errors.New("Failed to create a pod " + BuildTaskInstance.Name)
	// 			return err
	// 		}
	// 		s.End(true)
	// 	} else if task.Type == idp.RuntimeTask {
	// 		s := log.Spinner("Creating component")
	// 		defer s.End(false)
	// 		if err = BuildTaskInstance.CreateComponent(Client, componentConfig, cmpPVC); err != nil {
	// 			err = errors.New("Unable to create component deployment: " + err.Error())
	// 			return err
	// 		}
	// 		s.End(true)
	// 	}

	// 	glog.V(0).Info("Successfully created a " + task.Type + " Container")
	// 	glog.V(0).Info("===============================")
	// }

	// watchOptions.TimeoutSeconds = &noTimeout

	// Only sync project to the Container if a Source Mapping is provided
	// if len(srcDestination) > 0 {
	// 	err = syncToRunningContainer(Client, watchOptions, cwd, BuildTaskInstance.SrcDestination, []string{})
	// 	if err != nil {
	// 		glog.V(0).Infof("Error occured while syncing project to the %s Container: %s\n", task.Type, err)
	// 		err = errors.New("Unable to sync to the pod: " + err.Error())
	// 		return err
	// 	}
	// }

	// Only sync scripts to the Container if a Source Mapping is provided
	// if len(task.RepoMappings) > 0 {
	// 	for _, rm := range task.RepoMappings {
	// 		idpYamlDir, _ := filepath.Split(cwd + idp.IDPYamlPath)
	// 		sourcePath := idpYamlDir + rm.SrcPath
	// 		destinationPath := rm.DestPath
	// 		sourceDir, _ := filepath.Split(sourcePath)

	// 		err = syncToRunningContainer(Client, watchOptions, sourceDir, destinationPath, []string{sourcePath})
	// 		if err != nil {
	// 			glog.V(0).Infof("Error occured while syncing scripts to the %s Container: %s\n", task.Type, err)
	// 			err = errors.New("Unable to sync to the pod: " + err.Error())
	// 			return err
	// 		}
	// 	}
	// }
	return
}

func executeCmd(client *kclient.Client, cmd []string, pod *corev1.Pod) error {

	// Execute the tasks in the specified Container
	for _, task := range cmd {
		command := []string{"/bin/sh", "-c", task}

		glog.V(0).Infof("Executing command %s in the pod %s", command, pod.Name)

		// err := client.ExecCMDInContainer(podName, "", command, os.Stdout, os.Stdout, nil, false)
		// if err != nil {
		// 	glog.V(0).Infof("Error occured while executing command %s in the pod %s: %s\n", strings.Join(command, " "), podName, err)
		// 	err = errors.New("Unable to exec command " + strings.Join(command, " ") + " in the runtime container: " + err.Error())
		// 	return err
		// }
	}

	return nil
}
