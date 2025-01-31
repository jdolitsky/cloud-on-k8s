// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package kb

import (
	"fmt"
	"testing"

	commonv1beta1 "github.com/elastic/cloud-on-k8s/pkg/apis/common/v1beta1"
	"github.com/elastic/cloud-on-k8s/pkg/apis/kibana/v1beta1"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/annotation"
	"github.com/elastic/cloud-on-k8s/pkg/controller/common/events"
	"github.com/elastic/cloud-on-k8s/pkg/utils/k8s"
	"github.com/elastic/cloud-on-k8s/test/e2e/test"
	"github.com/elastic/cloud-on-k8s/test/e2e/test/elasticsearch"
	"github.com/elastic/cloud-on-k8s/test/e2e/test/kibana"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

// TestCrossNSAssociation tests associating Elasticsearch and Kibana running in different namespaces.
func TestCrossNSAssociation(t *testing.T) {
	// This test currently does not work in the E2E environment because each namespace has a dedicated
	// controller (see https://github.com/elastic/cloud-on-k8s/issues/1438)
	if !test.Ctx().Local {
		t.SkipNow()
	}

	esNamespace := test.Ctx().ManagedNamespace(0)
	kbNamespace := test.Ctx().ManagedNamespace(1)
	name := "test-cross-ns-assoc"

	esBuilder := elasticsearch.NewBuilder(name).
		WithNamespace(esNamespace).
		WithESMasterDataNodes(1, elasticsearch.DefaultResources).
		WithRestrictedSecurityContext()
	kbBuilder := kibana.NewBuilder(name).
		WithNamespace(kbNamespace).
		WithElasticsearchRef(esBuilder.Ref()).
		WithNodeCount(1).
		WithRestrictedSecurityContext()

	builders := []test.Builder{esBuilder, kbBuilder}
	test.RunMutations(t, builders, builders)
}

func TestKibanaAssociationWithNonExistentES(t *testing.T) {
	name := "test-kb-assoc-non-existent-es"
	kbBuilder := kibana.NewBuilder(name).
		WithNodeCount(1)

	k := test.NewK8sClientOrFatal()
	steps := test.StepList{}
	steps = steps.WithSteps(kbBuilder.InitTestSteps(k))
	steps = steps.WithSteps(kbBuilder.CreationTestSteps(k))
	steps = steps.WithStep(test.Step{
		Name: "Non existent backend should generate event",
		Test: test.Eventually(func() error {
			eventList, err := k.GetEvents(test.EventListOptions(kbBuilder.Kibana.Namespace, kbBuilder.Kibana.Name)...)
			if err != nil {
				return err
			}

			for _, evt := range eventList {
				if evt.Type == corev1.EventTypeWarning && evt.Reason == events.EventAssociationError {
					return nil
				}
			}

			return fmt.Errorf("event did not fire: %s", events.EventAssociationError)
		}),
	})
	steps = steps.WithSteps(kbBuilder.DeletionTestSteps(k))

	steps.RunSequential(t)
}

func TestKibanaAssociationWhenReferencedESDisappears(t *testing.T) {
	name := "test-kb-del-referenced-es"
	esBuilder := elasticsearch.NewBuilder(name).
		WithESMasterDataNodes(1, elasticsearch.DefaultResources)
	kbBuilder := kibana.NewBuilder(name).
		WithElasticsearchRef(esBuilder.Ref()).
		WithNodeCount(1)

	failureSteps := func(k *test.K8sClient) test.StepList {
		return test.StepList{
			test.Step{
				Name: "Updating to invalid Elasticsearch reference should succeed",
				Test: func(t *testing.T) {
					var kb v1beta1.Kibana
					require.NoError(t, k.Client.Get(k8s.ExtractNamespacedName(&kbBuilder.Kibana), &kb))
					kb.Spec.ElasticsearchRef.Namespace = "xxxx"
					require.NoError(t, k.Client.Update(&kb))
				},
			},
			test.Step{
				Name: "Lost Elasticsearch association should generate events",
				Test: test.Eventually(func() error {
					eventList, err := k.GetEvents(test.EventListOptions(kbBuilder.Kibana.Namespace, kbBuilder.Kibana.Name)...)
					if err != nil {
						return err
					}

					assocEstablishedEventSeen := false
					assocLostEventSeen := false
					noBackendEventSeen := false

					for _, evt := range eventList {
						switch {
						case evt.Type == corev1.EventTypeNormal && evt.Reason == events.EventAssociationStatusChange:
							prevStatus, currStatus := annotation.ExtractAssociationStatus(evt.ObjectMeta)
							if prevStatus == commonv1beta1.AssociationEstablished && currStatus != prevStatus {
								assocLostEventSeen = true
							}

							if currStatus == commonv1beta1.AssociationEstablished {
								assocEstablishedEventSeen = true
							}
						case evt.Type == corev1.EventTypeWarning && evt.Reason == events.EventAssociationError:
							noBackendEventSeen = true
						}
					}

					if assocEstablishedEventSeen && assocLostEventSeen && noBackendEventSeen {
						return nil
					}

					return fmt.Errorf("expected events did not fire: assocEstablished=%v assocLost=%v noBackend=%v",
						assocEstablishedEventSeen, assocLostEventSeen, noBackendEventSeen)
				}),
			},
		}
	}

	test.RunUnrecoverableFailureScenario(t, failureSteps, kbBuilder, esBuilder)
}
