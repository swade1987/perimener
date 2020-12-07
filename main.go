package main

import (
	"context"
	"crypto/rand"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog/log"
	"github.com/swade1987/perimener/pkg/pods"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type config struct {
	DelaySeconds          int64  `env:"DELAY_SECONDS" envDefault:"0"`
	ExpectedReadyPodCount int    `env:"EXPECTED_READY_POD_COUNT,required"`
	Namespace             string `env:"NAMESPACE,required"`
	PodLabel              string `env:"POD_LABEL,required"`
	RandomWindowSeconds   int64  `env:"RANDOM_WINDOW_SECONDS" envDefault:"0"`
	SleepCount            int    `env:"SLEEP_COUNT" envDefault:"5"`
	UseLocalKubeConfig    bool   `env:"USE_LOCAL_KUBECONFIG" envDefault:"false"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal().Err(err).Msg("failed to run")
	}
	log.Info().Msg("ending perimener - gracefully exiting")
}

func run() error {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		return err
	}

	log.Info().
		Int64("delay_seconds", cfg.DelaySeconds).
		Int("expected_ready_pod_count", cfg.ExpectedReadyPodCount).
		Str("namespace", cfg.Namespace).
		Str("pod_label", cfg.PodLabel).
		Int64("random_window_seconds", cfg.RandomWindowSeconds).
		Int("sleep_count", cfg.SleepCount).
		Bool("use_local_kube_config", cfg.UseLocalKubeConfig).
		Msg("starting perimener")

	// Create an rest client not targeting a specific API version
	kCfg, err := constructKubernetesConfig(cfg.UseLocalKubeConfig)
	if err != nil {
		return err
	}
	client, err := kubernetes.NewForConfig(kCfg)
	if err != nil {
		return err
	}

	// Wait until we have X number of pods in a Ready state and then gracefully exit.
	failureCount := 0
	for {

		// Get the list of pods
		var podList, err = client.CoreV1().Pods(cfg.Namespace).List(context.TODO(), metav1.ListOptions{
			LabelSelector: cfg.PodLabel,
		})

		if err != nil {
			log.Info().
				Err(err).
				Str("namespace", cfg.Namespace).
				Str("pod_label", cfg.PodLabel).
				Msg("no pods currently available matching label selector in namespace")
		}

		currentReadyPods := pods.ReadyCount(podList)
		if currentReadyPods < cfg.ExpectedReadyPodCount {
			failureCount++
			log.Info().
				Int("current_ready_pods", currentReadyPods).
				Int("expected_ready_pod_count", cfg.ExpectedReadyPodCount).
				Int("failure_count", failureCount).
				Str("namespace", cfg.Namespace).
				Str("pod_label", cfg.PodLabel).
				Int("sleep_count", cfg.SleepCount).
				Msg("insufficient pods in a ready state; will retry after delay")
			time.Sleep(time.Duration(cfg.SleepCount) * time.Second)
			continue
		}

		log.Info().
			Str("namespace", cfg.Namespace).
			Str("pod_label", cfg.PodLabel).
			Int("current_ready_pods", currentReadyPods).
			Msg("sufficient pods in a ready state")

		if cfg.RandomWindowSeconds != 0 && failureCount > 1 {
			waitForSeconds(cfg.RandomWindowSeconds)
		}

		if cfg.DelaySeconds > 0 {
			log.Info().
				Str("namespace", cfg.Namespace).
				Str("pod_label", cfg.PodLabel).
				Int("current_ready_pods", currentReadyPods).
				Int64("delay_seconds", cfg.DelaySeconds).
				Msg("delaying start of containers")
			time.Sleep(time.Duration(cfg.DelaySeconds) * time.Second)
		}
		return nil
	}
}

func constructKubernetesConfig(useLocalKubeConfig bool) (*rest.Config, error) {
	if useLocalKubeConfig {
		kubeconfig := filepath.Join(os.Getenv("PWD"), "kubeconfig")
		log.Info().
			Str("kubeconfig_path", kubeconfig).
			Msg("using local kubeconfig")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	return clientcmd.BuildConfigFromFlags("", "")
}

func waitForSeconds(windowSeconds int64) {
	randomSeconds, _ := rand.Int(rand.Reader, big.NewInt(windowSeconds))

	log.Info().
		Int64("wait_duration_seconds", randomSeconds.Int64()).
		Int64("max_duration_seconds", windowSeconds).
		Msg("sleeping for random amount of time in [0, wait_duration_seconds)")
	time.Sleep(time.Duration(randomSeconds.Int64()) * time.Second)
}
