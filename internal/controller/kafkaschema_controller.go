package controller

import (
	"context"
	"strconv"
	"time"

	"incubly.oss/kafka-schema-operator/api/v1beta1"
	"incubly.oss/kafka-schema-operator/internal/schemareg"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// KafkaSchemaReconciler reconciles a KafkaSchema object
type KafkaSchemaReconciler struct {
	RequeueDelay         time.Duration
	DefaultCleanupPolicy v1beta1.CleanupPolicy
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kafka.incubly.oss,resources=kafkaschemas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kafka.incubly.oss,resources=kafkaschemas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kafka.incubly.oss,resources=kafkaschemas/finalizers,verbs=update

const finalizer = "kafka.incubly.oss/finalizer"

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// This controller will try to apply state of KafkaSchema resources on
// Schema Registry/Registries.
// It is unidirectional (i.e. it only synchronizes state from Kube to Registry,
// without trying to create KafkaSchema resources for discovered Registry
// schemas and subjects that don't have their representation in Kube).
//
// During reconciliation, it will synchronize subject, schema and compatibility mode.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *KafkaSchemaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	res := &v1beta1.KafkaSchema{}
	err := r.Get(ctx, req.NamespacedName, res)
	if err != nil {
		if errors.IsNotFound(err) {
			log.FromContext(ctx).Info("KafkaSchema CR not found. I can't do anything about it...")
			// finish without reconciliation loop - CR deleted(?)
			return ctrl.Result{}, nil
		}
		log.FromContext(ctx).Error(err, "Failed to get KafkaSchema CR")
		// not reflected in resource status, hopefully it's quite unlikely
		return ctrl.Result{}, err
	}

	logger := log.FromContext(
		ctx,
		"resource::name", res.Name,
		"resource::generation", strconv.FormatInt(res.Generation, 10),
		"resource::uid", res.UID,
	)

	if res.Generation == res.Status.ObservedGeneration {
		logger.Info(".statusObservedGeneration matches resource generation. No reconciliation needed")
		return ctrl.Result{Requeue: r.RequeueDelay >= 0, RequeueAfter: r.RequeueDelay}, nil
	}

	spec := res.Spec

	srClient, err := schemareg.NewClient(&spec.SchemaRegistry, logger)

	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.SchemaRegistryClient,
			"Failed to instantiate Schema Registry Client")
	}

	subjectName, err := resolveSubjectName(&spec)
	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.NameStrategy,
			"Failed to resolve subject name")
	}

	if res.GetDeletionTimestamp().IsZero() {
		res.Status.SchemaRegistryUrl = srClient.BaseUrl.String()
		res.Status.Subject = subjectName
		err := r.Status().Update(ctx, res)
		if err != nil {
			logger.Error(err, "failed to update resource status")
			return ctrl.Result{}, err
		}
		logger.Info("starting reconciliation for KafkaSchema CR")
		return r.reconcileResource(ctx, res, srClient, logger)
	} else {
		logger.Info("starting deletion for KafkaSchema CR")
		return r.deleteResource(ctx, res, srClient, logger)
	}
}

func (r *KafkaSchemaReconciler) reconcileResource(
	ctx context.Context,
	res *v1beta1.KafkaSchema,
	srClient *schemareg.SrClient,
	logger logr.Logger) (ctrl.Result, error) {

	subjectName := res.Status.Subject
	spec := res.Spec

	if controllerutil.AddFinalizer(res, finalizer) {
		err := r.Update(ctx, res)
		if err != nil {
			return r.logError(logger, err, ctx, res,
				v1beta1.ResourceUpdate,
				"Failed to add finalizer to KafkaSchema CR")
		}
	}

	maybeNormalizedSchema, err := GetMaybeNormalizedSchema(spec.Data)

	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.NormalizeSchema,
			"Failed to normalize schema")
	}

	schemaId, err := srClient.RegisterSchema(
		subjectName,
		schemareg.RegisterSchemaReq{
			Schema:     maybeNormalizedSchema,
			SchemaType: spec.Data.Format,
		},
	)
	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.RegisterSchema,
			"Failed to register schema in registry")
	}
	res.Status.SchemaId = schemaId

	compatibility := spec.Data.Compatibility
	if len(compatibility) > 0 {
		err = srClient.SetCompatibilityMode(
			subjectName,
			schemareg.SetCompatibilityModeReq{
				Compatibility: compatibility,
			})

		if err != nil {
			return r.logError(logger, err, ctx, res,
				v1beta1.SetCompatibilityMode,
				"Failed to update schema compatibility mode")
		}
	}

	res.SetReadyReason(v1beta1.Complete, "Reconciliation complete")
	res.Status.Healthy = true
	res.Status.Status = "True"
	res.Status.ObservedGeneration = res.Generation

	// force-update (without if on SetReadyReason) to apply all changes
	if err := r.Status().Update(ctx, res); err != nil {
		/*
			not reflected in resource status:
			failure on resource update might fail on resource update ;)
		*/
		logger.Error(err, "Failed to update status successful reconciliation")
		return ctrl.Result{}, err
	}

	logger.Info("KafkaSchema CR successfully reconciled")
	/*
		delay<0 - don't requeue (Requeue: false)
		delay=0 - exponential backoff (Requeue: true, RequeueAfter: 0)
		delay>0 - static interval (Requeue: true, RequeueAfter should be respected)
	*/
	return ctrl.Result{Requeue: r.RequeueDelay >= 0, RequeueAfter: r.RequeueDelay}, nil
}

func (r *KafkaSchemaReconciler) deleteResource(
	ctx context.Context,
	res *v1beta1.KafkaSchema,
	srClient *schemareg.SrClient,
	logger logr.Logger) (ctrl.Result, error) {

	// deleting / cleaning up resource
	err := performCleanup(res, srClient)
	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.Cleanup,
			"Failed to perform schema registry cleanup")
	}

	// deleting CR
	controllerutil.RemoveFinalizer(res, finalizer)
	err = r.Update(ctx, res)
	if err != nil {
		return r.logError(logger, err, ctx, res,
			v1beta1.Cleanup,
			"Failed to delete KafkaSchema CR")
	}
	logger.Info("KafkaSchema CR successfully deleted")

	// finish without reconciliation loop - CR deleted
	return ctrl.Result{}, nil
}

func (r *KafkaSchemaReconciler) logError(
	logger logr.Logger,
	err error,
	ctx context.Context,
	res *v1beta1.KafkaSchema,
	reason v1beta1.ReadyReason,
	msg string,
) (ctrl.Result, error) {
	logger.Error(err, msg)
	res.Status.Healthy = false
	res.Status.Status = "False"
	res.Status.ObservedGeneration = res.Generation

	if res.SetReadyReason(reason, msg) {
		// ignoring the update error - we should return the actual root cause instead
		_ = r.Status().Update(ctx, res)
	}

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *KafkaSchemaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1beta1.KafkaSchema{}).
		Complete(r)
}
