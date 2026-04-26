package controllers

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	examplev1 "pod-counter-crd/api/v1"
)

type PodCounterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Clientset *kubernetes.Clientset
}

// +kubebuilder:rbac:groups=example.com,resources=podcounters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.com,resources=podcounters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

func (r *PodCounterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Get the PodCounter resource
	podCounter := &examplev1.PodCounter{}
	if err := r.Get(ctx, req.NamespacedName, podCounter); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("PodCounter resource not found, ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch PodCounter")
		return ctrl.Result{}, err
	}

	// Get all namespaces with label "monitor-pods=true"
	namespaces := &corev1.NamespaceList{}
	if err := r.List(ctx, namespaces, client.MatchingLabels{"monitor-pods": "true"}); err != nil {
		log.Error(err, "unable to list namespaces")
		return ctrl.Result{}, err
	}

	monitoredNamespaces := []string{}
	podCounts := make(map[string]int32)

	// Count pods in each labeled namespace
	for _, ns := range namespaces.Items {
		nsName := ns.Name
		pods := &corev1.PodList{}
		if err := r.List(ctx, pods, client.InNamespace(nsName)); err != nil {
			log.Error(err, "unable to list pods", "namespace", nsName)
			continue
		}

		podCount := int32(len(pods.Items))
		monitoredNamespaces = append(monitoredNamespaces, nsName)
		podCounts[nsName] = podCount

		log.Info("Pod count", "namespace", nsName, "count", podCount)
	}

	// Update status
	podCounter.Status.MonitoredNamespaces = monitoredNamespaces
	podCounter.Status.PodCounts = podCounts
	podCounter.Status.LastChecked = time.Now().Format(time.RFC3339)

	if err := r.Status().Update(ctx, podCounter); err != nil {
		log.Error(err, "unable to update PodCounter status")
		return ctrl.Result{}, err
	}

	log.Info("PodCounter reconciled successfully", "monitored_namespaces", len(monitoredNamespaces))

	// Requeue after interval
	interval := time.Duration(podCounter.Spec.Interval) * time.Second
	if interval == 0 {
		interval = 10 * time.Second
	}

	return ctrl.Result{RequeueAfter: interval}, nil
}

func (r *PodCounterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.PodCounter{}).
		Watches(&corev1.Pod{}, nil).
		Watches(&corev1.Namespace{}, nil).
		Complete(r)
}
