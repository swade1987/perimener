package pods

import (
	v1 "k8s.io/api/core/v1"
	"log"
)

func ReadyCount(podList *v1.PodList) (readyPods int) {
	for _, pod := range podList.Items {

		// A deletion timestamp indicates that a pod is terminating. Do not count this pod.
		if pod.ObjectMeta.DeletionTimestamp != nil {
			continue
		}

		var noOfContainers = len(pod.Status.ContainerStatuses)
		var readyContainers = 0

		for _, status := range pod.Status.ContainerStatuses {
			if status.Ready == true {
				readyContainers++
			} else {
				log.Printf("Container %s is not Ready.", status.Name)
			}
		}

		if noOfContainers == readyContainers {
			readyPods++
		}
	}

	return
}
