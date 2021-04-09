// Copyright 2020 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2016-03-30/containerservice"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

type KubernetesGenerator struct {
	AzureService
}

func (g *KubernetesGenerator) listKubernetesServers() ([]terraformutils.Resource, error) {
	var resources []terraformutils.Resource
	ctx := context.Background()
	subscriptionID := g.Args["config"].(authentication.Config).SubscriptionID
	KubernetesClient := Kubernetes.NewClient(subscriptionID)

	KubernetesServersIterator, err := KubernetesClient.ListComplete(ctx)
	if err != nil {
		return nil, err
	}

	for KubernetesServersIterator.NotDone() {
		KubernetesServer := KubernetesServersIterator.Value()
		resources = append(resources, terraformutils.NewSimpleResource(
			*KubernetesServer.ID,
			*KubernetesServer.Name,
			"azurerm_kubernetes_cluster",
			g.ProviderName,
			[]string{}))

		if err := KubernetesServersIterator.Next(); err != nil {
			log.Println(err)
			break
		}
	}

	return resources, nil
}

func (g *KubernetesGenerator) InitResources() error {
	functions := []func() ([]terraformutils.Resource, error){
		g.listKubernetesServers,
	}

	for _, f := range functions {
		resources, err := f()
		if err != nil {
			return err
		}
		g.Resources = append(g.Resources, resources...)
	}

	return nil
}
