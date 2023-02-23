package pods

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestReadyCount(t *testing.T) {

	var (
		now = metav1.Now()

		podWithDeletionTimestamp = v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				DeletionTimestamp: &now,
			},
		}

		podWithReadyContainers = v1.Pod{
			Status: v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{Ready: true},
					{Ready: true},
				},
			},
		}

		podWithNonReadyContainers = v1.Pod{
			Status: v1.PodStatus{
				ContainerStatuses: []v1.ContainerStatus{
					{Ready: false},
					{Ready: true},
				},
			},
		}

		inputReadyCount = []v1.PodList{

			{
				Items: []v1.Pod{
					podWithDeletionTimestamp,
					podWithReadyContainers,
					podWithReadyContainers,
				},
			},

			{
				Items: []v1.Pod{
					podWithReadyContainers,
					podWithReadyContainers,
				},
			},

			{
				Items: []v1.Pod{
					podWithNonReadyContainers,
					podWithReadyContainers,
				},
			},

			{
				Items: []v1.Pod{
					podWithNonReadyContainers,
				},
			},

			{
				Items: []v1.Pod{
					podWithDeletionTimestamp,
					podWithNonReadyContainers,
				},
			},

			{
				Items: []v1.Pod{
					podWithDeletionTimestamp,
				},
			},

			{
				Items: []v1.Pod{},
			},
		}

		outputReadyCount = []int{2, 2, 1, 0, 0, 0, 0}
	)

	for i, podList := range inputReadyCount {
		t.Run("", func(t *testing.T) {
                        // The input variables mutate with potentially unwanted effects 
                        podlist := podList
                        index := i
			t.Parallel()
			got := ReadyCount(&podlist)
			want := outputReadyCount[index]
			if got != want {
				t.Errorf("Got %q, want %q", got, want)
			}
		})
	}
}
