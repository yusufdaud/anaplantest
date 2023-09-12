package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"

	testgroupv1 "github.com/yusufdaud/anaplantest/api/v1"
)

// YusufReconciler reconciles a Yusuf object
type YusufReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=testgroup.daud.com,resources=yusufs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=testgroup.daud.com,resources=yusufs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=testgroup.daud.com,resources=yusufs/finalizers,verbs=update
func (r *YusufReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("reconciling yusuf custom resource")

	var yusuf testgroupv1.Yusuf

	if err := r.Get(ctx, req.NamespacedName, &yusuf); err != nil {
		log.Error(err, "unable to fetch Yusuf")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if yusuf.Status.Happy {
		log.Info("yusuf is already happy")
		return ctrl.Result{}, nil
	}

	cl, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		log.Error(err, "failed to create client")
		return ctrl.Result{}, err
	}

	log.Info("creating pod")
	// Using a typed object.
	err = createPod(yusuf, err, cl)
	if err != nil {
		log.Error(err, "failed to create pod")
		return ctrl.Result{}, err
	}

	// Update Yusuf happy status
	yusuf.Status.Happy = true
	if err := r.Status().Update(ctx, &yusuf); err != nil {
		log.Error(err, "unable to update yusuf's happy status", "status", true)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func createPod(y testgroupv1.Yusuf, err error, cl client.Client) error {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      y.Name,
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: y.APIVersion,
					Kind:       y.Kind,
					Name:       y.Name,
					UID:        y.UID,
				},
			},
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image: "nginx",
					Name:  y.Name,
				},
			},
		},
	}

	err = cl.Create(context.Background(), pod)
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *YusufReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testgroupv1.Yusuf{}).
		Complete(r)
}
