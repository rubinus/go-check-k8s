package phases

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rubinus/go-check-k8s/util/kobe"
	"github.com/spf13/viper"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"time"
)

const (
	PhaseInterval             = 5 * time.Second
	DefaultPhaseTimeoutMinute = 5
)

func RunPlaybookAndGetResult(ctx context.Context, resCh chan *kobe.AllResult, b kobe.Interface, playbookName, tag string, writer io.Writer) error {
	taskId, err := b.RunPlaybook(playbookName, tag)

	var result kobe.Result
	if err != nil {
		return err
	}

	// 读取 ansible 执行日志
	if writer != nil {
		go func() {
			err = b.Watch(writer, taskId)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	timeout := viper.GetInt("job.timeout")
	if timeout < DefaultPhaseTimeoutMinute {
		timeout = DefaultPhaseTimeoutMinute
	}
	err = wait.PollUntilContextTimeout(ctx, PhaseInterval, time.Duration(timeout)*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		res, err := b.GetResult(taskId)
		if err != nil {
			return true, err
		}
		//fmt.Println("结果:", res.Content)
		if res.Finished {
			if res.Success {
				result, err = kobe.ParseResult(res.Content)
				if err != nil {
					return true, err
				}
				allResult := &kobe.AllResult{
					Result: result,
					TaskId: taskId,
				}
				resCh <- allResult

			} else {
				if res.Content != "" {
					result, err = kobe.ParseResult(res.Content)
					if err != nil {
						return true, err
					}
					result.GatherFailedInfo()
					if result.HostFailedInfo != nil && len(result.HostFailedInfo) > 0 {
						by, _ := json.Marshal(&result.HostFailedInfo)
						return true, errors.New(string(by))
					}
				}
			}
			return true, nil
		}
		return false, nil
	})
	return err
}

func WaitForDeployRunning(ctx context.Context, namespace string, deploymentName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 2*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		d, err := kubeClient.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WaitForDaemonSetRunning(ctx context.Context, namespace string, daemonsetName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 2*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		d, err := kubeClient.AppsV1().DaemonSets(namespace).Get(ctx, daemonsetName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.NumberReady > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WaitForStatefulSetsRunning(ctx context.Context, namespace string, statefulSetsName string, kubeClient *kubernetes.Clientset) error {
	kubeClient.CoreV1()
	err := wait.PollUntilContextTimeout(ctx, 5*time.Second, 2*time.Minute, true, func(ctx context.Context) (done bool, err error) {
		d, err := kubeClient.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetsName, metav1.GetOptions{})
		if err != nil {
			return true, err
		}
		if d.Status.ReadyReplicas > 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		return err
	}
	return nil
}
