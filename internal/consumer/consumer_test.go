/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package consumer_test

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	pvcexplorerv1alpha1 "github.com/pvc-explorer-operator/pvc-explorer/api/v1alpha1"
	"github.com/pvc-explorer-operator/pvc-explorer/internal/consumer"
)

const (
	testNS      = "default"
	testCron    = "my-cron"
	testPodName = "pod1"
	testDeploy  = "my-deploy"
)

func scheme() *runtime.Scheme {
	s := runtime.NewScheme()
	_ = corev1.AddToScheme(s)
	_ = appsv1.AddToScheme(s)
	_ = batchv1.AddToScheme(s)
	_ = pvcexplorerv1alpha1.AddToScheme(s)
	return s
}

func pvcVolume(claimName string, readOnly bool) corev1.Volume {
	return corev1.Volume{
		Name: claimName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: claimName,
				ReadOnly:  readOnly,
			},
		},
	}
}

func runningPod(pvcName string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: testPodName, Namespace: testNS},
		Spec: corev1.PodSpec{
			NodeName: "node1",
			Volumes:  []corev1.Volume{pvcVolume(pvcName, false)},
		},
		Status: corev1.PodStatus{Phase: corev1.PodRunning},
	}
}

func TestMountsPVC_Match(t *testing.T) {
	pod := runningPod("my-pvc")
	if !consumer.MountsPVC(pod, "my-pvc") {
		t.Fatal("expected MountsPVC to return true")
	}
}

func TestMountsPVC_NoMatch(t *testing.T) {
	pod := runningPod("other-pvc")
	if consumer.MountsPVC(pod, "my-pvc") {
		t.Fatal("expected MountsPVC to return false")
	}
}

func TestDetect_PodOwner_Bare(t *testing.T) {
	pod := runningPod("my-pvc")
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(pod).Build()

	consumers, err := consumer.Detect(context.Background(), c, testNS, "my-pvc")
	if err != nil {
		t.Fatal(err)
	}
	if len(consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(consumers))
	}
	if consumers[0].OwnerKind != "Pod" || consumers[0].OwnerName != testPodName {
		t.Errorf("unexpected owner: %s/%s", consumers[0].OwnerKind, consumers[0].OwnerName)
	}
}

func TestDetect_PodToDeploymentViaReplicaSet(t *testing.T) {
	deploy := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: testDeploy, Namespace: testNS}}
	rs := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-rs",
			Namespace: testNS,
			OwnerReferences: []metav1.OwnerReference{
				{Kind: "Deployment", Name: testDeploy, APIVersion: "apps/v1"},
			},
		},
	}
	pod := runningPod("my-pvc")
	pod.OwnerReferences = []metav1.OwnerReference{
		{Kind: "ReplicaSet", Name: "my-rs", APIVersion: "apps/v1"},
	}
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(deploy, rs, pod).Build()

	consumers, err := consumer.Detect(context.Background(), c, testNS, "my-pvc")
	if err != nil {
		t.Fatal(err)
	}
	if len(consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(consumers))
	}
	if consumers[0].OwnerKind != "Deployment" || consumers[0].OwnerName != testDeploy {
		t.Errorf("unexpected owner: %s/%s", consumers[0].OwnerKind, consumers[0].OwnerName)
	}
}

func TestDetect_PodToCronJobViaJob(t *testing.T) {
	cronJob := &batchv1.CronJob{ObjectMeta: metav1.ObjectMeta{Name: testCron, Namespace: testNS}}
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-job",
			Namespace: testNS,
			OwnerReferences: []metav1.OwnerReference{
				{Kind: "CronJob", Name: "my-cron", APIVersion: "batch/v1"},
			},
		},
	}
	pod := runningPod("my-pvc")
	pod.OwnerReferences = []metav1.OwnerReference{
		{Kind: "Job", Name: "my-job", APIVersion: "batch/v1"},
	}
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(cronJob, job, pod).Build()

	consumers, err := consumer.Detect(context.Background(), c, testNS, "my-pvc")
	if err != nil {
		t.Fatal(err)
	}
	if len(consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(consumers))
	}
	if consumers[0].OwnerKind != "CronJob" || consumers[0].OwnerName != testCron {
		t.Errorf("unexpected owner: %s/%s", consumers[0].OwnerKind, consumers[0].OwnerName)
	}
}

func TestDetect_SkipsTerminatedPods(t *testing.T) {
	pod := runningPod("my-pvc")
	pod.Status.Phase = corev1.PodSucceeded
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(pod).Build()

	consumers, err := consumer.Detect(context.Background(), c, testNS, "my-pvc")
	if err != nil {
		t.Fatal(err)
	}
	if len(consumers) != 0 {
		t.Fatalf("expected 0 consumers, got %d", len(consumers))
	}
}

func TestIndex_GetReturnsNilForUnknown(t *testing.T) {
	idx := consumer.NewIndex()
	result := idx.Get(testNS, "no-such-pvc")
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestIndex_SyncAndGet(t *testing.T) {
	pod := runningPod("my-pvc")
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(pod).Build()

	idx := consumer.NewIndex()
	if err := idx.Sync(context.Background(), c, testNS); err != nil {
		t.Fatal(err)
	}

	consumers := idx.Get(testNS, "my-pvc")
	if len(consumers) != 1 {
		t.Fatalf("expected 1 consumer, got %d", len(consumers))
	}
	if consumers[0].PodName != testPodName {
		t.Errorf("unexpected pod name: %s", consumers[0].PodName)
	}
}

func TestIndex_SyncClearsStaleEntries(t *testing.T) {
	pod := runningPod("my-pvc")
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(pod).Build()

	idx := consumer.NewIndex()
	_ = idx.Sync(context.Background(), c, testNS)

	c2 := fake.NewClientBuilder().WithScheme(scheme()).Build()
	if err := idx.Sync(context.Background(), c2, testNS); err != nil {
		t.Fatal(err)
	}

	consumers := idx.Get(testNS, "my-pvc")
	if len(consumers) != 0 {
		t.Fatalf("expected 0 consumers after stale sync, got %d", len(consumers))
	}
}

func TestIndex_GetReturnsCopy(t *testing.T) {
	pod := runningPod("my-pvc")
	c := fake.NewClientBuilder().WithScheme(scheme()).WithObjects(pod).Build()

	idx := consumer.NewIndex()
	_ = idx.Sync(context.Background(), c, testNS)

	consumers := idx.Get(testNS, "my-pvc")
	consumers[0].PodName = "mutated"

	consumers2 := idx.Get(testNS, "my-pvc")
	if consumers2[0].PodName == "mutated" {
		t.Fatal("Get should return a copy, not a reference")
	}
}
