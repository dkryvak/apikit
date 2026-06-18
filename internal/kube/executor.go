package kube

import (
	"apikit/internal/config"
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

const (
	jobName     = "apikit-job"
	podLabel    = "app=apikit"
	containerNm = "apikit"
)

// Executor runs commands inside an apikit pod in a given namespace/cluster.
type Executor struct {
	clientset *kubernetes.Clientset
	restCfg   *rest.Config
	namespace string
}

// NewExecutor builds a client-go executor for the given kube context (may be
// empty to use the current context) and namespace.
func NewExecutor(kubeContext, namespace string) (*Executor, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	if strings.TrimSpace(kubeContext) != "" {
		overrides.CurrentContext = kubeContext
	}

	clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
	restCfg, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build kube config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build kube client: %w", err)
	}

	return &Executor{clientset: clientset, restCfg: restCfg, namespace: namespace}, nil
}

// EnsurePod returns the name of a ready apikit pod, reusing an existing one or
// creating the Job if none is present.
func (e *Executor) EnsurePod(ctx context.Context) (string, error) {
	if name, ok := e.findReadyPod(ctx); ok {
		return name, nil
	}

	if err := e.createJob(ctx); err != nil && !apierrors.IsAlreadyExists(err) {
		return "", fmt.Errorf("failed to create apikit job: %w", err)
	}

	return e.waitForReadyPod(ctx, time.Duration(config.PodReadyTimeoutSeconds())*time.Second)
}

func (e *Executor) findReadyPod(ctx context.Context) (string, bool) {
	pods, err := e.clientset.CoreV1().Pods(e.namespace).List(ctx, metav1.ListOptions{LabelSelector: podLabel})
	if err != nil {
		return "", false
	}
	for _, p := range pods.Items {
		if p.Status.Phase != corev1.PodRunning {
			continue
		}
		for _, c := range p.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
				return p.Name, true
			}
		}
	}
	return "", false
}

func (e *Executor) createJob(ctx context.Context) error {
	img, err := config.Image()
	if err != nil {
		return err
	}

	lifetime := config.PodLifetimeSeconds()
	backoff := int32(0)
	ttl := int32(config.JobTTLSeconds())
	maxLife := int64(lifetime)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{Name: jobName, Namespace: e.namespace},
		Spec: batchv1.JobSpec{
			BackoffLimit:            &backoff,
			TTLSecondsAfterFinished: &ttl,
			ActiveDeadlineSeconds:   &maxLife,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "apikit"}},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{{
						Name:    containerNm,
						Image:   img,
						Command: []string{"sleep"},
						Args:    []string{strconv.Itoa(lifetime)},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("250m"),
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceCPU:    resource.MustParse("250m"),
								corev1.ResourceMemory: resource.MustParse("100Mi"),
							},
						},
					}},
				},
			},
		},
	}

	_, err = e.clientset.BatchV1().Jobs(e.namespace).Create(ctx, job, metav1.CreateOptions{})
	return err
}

func (e *Executor) waitForReadyPod(ctx context.Context, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	for {
		if name, ok := e.findReadyPod(ctx); ok {
			return name, nil
		}
		if time.Now().After(deadline) {
			return "", fmt.Errorf("timed out waiting for apikit pod to become ready in %s", e.namespace)
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}
}

// Exec runs args inside the pod, optionally feeding stdin, and returns the
// captured stdout and stderr.
func (e *Executor) Exec(ctx context.Context, pod string, args []string, stdin []byte) (stdout, stderr []byte, err error) {
	req := e.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod).
		Namespace(e.namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerNm,
			Command:   args,
			Stdin:     len(stdin) > 0,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(e.restCfg, "POST", req.URL())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init exec: %w", err)
	}

	var outBuf, errBuf bytes.Buffer
	streamOpts := remotecommand.StreamOptions{
		Stdout: &outBuf,
		Stderr: &errBuf,
	}
	if len(stdin) > 0 {
		streamOpts.Stdin = bytes.NewReader(stdin)
	}

	if err := executor.StreamWithContext(ctx, streamOpts); err != nil {
		return outBuf.Bytes(), errBuf.Bytes(), fmt.Errorf("exec failed: %w", err)
	}

	return outBuf.Bytes(), errBuf.Bytes(), nil
}
