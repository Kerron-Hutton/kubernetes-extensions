package v1alpha1

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppType string

const (
	Frontend AppType = "frontend"
	Backend  AppType = "backend"
)

func constructService(appSet *ApplicationSet, appType AppType) *corev1.Service {
	var port int32

	if appType == Frontend {
		port = appSet.Spec.Frontend.Port
	} else {
		port = appSet.Spec.Backend.Port
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-svc", appSet.Name, appType),
			Namespace: appSet.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "application-set-operator",
				"app.kubernetes.io/component":  string(appType),
				"app.kubernetes.io/name":       appSet.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/component": string(appType),
				"app.kubernetes.io/name":      appSet.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     int32(port),
				},
			},
		},
	}
}

func constructDeployment(appSet *ApplicationSet, appType AppType) *appsv1.Deployment {
	var numOfReplicas, port int32
	var image string

	if appType == Frontend {
		numOfReplicas = appSet.Spec.Frontend.Replicas
		image = appSet.Spec.Frontend.Image
		port = appSet.Spec.Frontend.Port
	} else {
		numOfReplicas = appSet.Spec.Backend.Replicas
		image = appSet.Spec.Frontend.Image
		port = appSet.Spec.Backend.Port
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-%s-deploy", appSet.Name, appType),
			Namespace: appSet.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "application-set-operator",
				"app.kubernetes.io/component":  string(appType),
				"app.kubernetes.io/name":       appSet.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &numOfReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/component": string(appType),
					"app.kubernetes.io/name":      appSet.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/component": string(appType),
						"app.kubernetes.io/name":      appSet.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: image,
							Name:  appSet.Name,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (appSet *ApplicationSet) ConstructAppSet() []client.Object {
	var objects []client.Object

	frontendDeployment := constructDeployment(appSet, Frontend)
	frontendService := constructService(appSet, Frontend)

	objects = append(objects, frontendDeployment)
	objects = append(objects, frontendService)

	backendDeployment := constructDeployment(appSet, Backend)
	backendService := constructService(appSet, Backend)

	objects = append(objects, backendDeployment)
	objects = append(objects, backendService)

	return objects
}
