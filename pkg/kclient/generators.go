package kclient

import (

	// api resource types
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateContainer creates a container spec that can be used when creating a pod
func GenerateContainer(name, image string, isPrivileged bool, command, args []string, envVars []corev1.EnvVar, resourceReqs corev1.ResourceRequirements) *corev1.Container {
	container := &corev1.Container{
		Name:            name,
		Image:           image,
		ImagePullPolicy: corev1.PullAlways,
		Resources:       resourceReqs,
		Command:         command,
		Args:            args,
		Env:             envVars,
	}

	if isPrivileged {
		container.SecurityContext = &corev1.SecurityContext{
			Privileged: &isPrivileged,
		}
	}

	return container
}

// GeneratePodTemplateSpec creates a pod template spec that can be used to create a deployment spec
func GeneratePodTemplateSpec(podName, namespace string, labels map[string]string, containers []corev1.Container) *corev1.PodTemplateSpec {
	podTemplateSpec := &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: containers,
		},
	}

	return podTemplateSpec
}

// GenerateDeploymentSpec creates a deployment spec
func GenerateDeploymentSpec(podTemplateSpec corev1.PodTemplateSpec) *appsv1.DeploymentSpec {
	labels := podTemplateSpec.ObjectMeta.Labels
	deploymentSpec := &appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: podTemplateSpec,
	}

	return deploymentSpec
}
