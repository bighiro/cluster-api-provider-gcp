/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGCPMachineTemplate_ValidateCreate(t *testing.T) {
	g := NewWithT(t)
	confidentialComputeEnabled := ConfidentialComputePolicyEnabled
	onHostMaintenanceTerminate := HostMaintenancePolicyTerminate
	onHostMaintenanceMigrate := HostMaintenancePolicyMigrate
	tests := []struct {
		name     string
		template *GCPMachineTemplate
		wantErr  bool
	}{
		{
			name: "GCPMachineTemplate with OnHostMaintenance set to Terminate - valid",
			template: &GCPMachineTemplate{
				Spec: GCPMachineTemplateSpec{
					Template: GCPMachineTemplateResource{
						Spec: GCPMachineSpec{
							InstanceType:      "n2d-standard-4",
							OnHostMaintenance: &onHostMaintenanceTerminate,
						}},
				},
			},
			wantErr: false,
		},
		{
			name: "GCPMachineTemplate with ConfidentialCompute enabled and OnHostMaintenance set to Terminate - valid",
			template: &GCPMachineTemplate{
				Spec: GCPMachineTemplateSpec{
					Template: GCPMachineTemplateResource{
						Spec: GCPMachineSpec{
							InstanceType:        "n2d-standard-4",
							ConfidentialCompute: &confidentialComputeEnabled,
							OnHostMaintenance:   &onHostMaintenanceTerminate,
						}},
				},
			},
			wantErr: false,
		},
		{
			name: "GCPMachineTemplate with ConfidentialCompute enabled and OnHostMaintenance set to Migrate - invalid",
			template: &GCPMachineTemplate{
				Spec: GCPMachineTemplateSpec{
					Template: GCPMachineTemplateResource{
						Spec: GCPMachineSpec{
							InstanceType:        "n2d-standard-4",
							ConfidentialCompute: &confidentialComputeEnabled,
							OnHostMaintenance:   &onHostMaintenanceMigrate,
						}},
				},
			},
			wantErr: true,
		},
		{
			name: "GCPMachineTemplate with ConfidentialCompute enabled and default OnHostMaintenance (Migrate) - invalid",
			template: &GCPMachineTemplate{
				Spec: GCPMachineTemplateSpec{
					Template: GCPMachineTemplateResource{
						Spec: GCPMachineSpec{
							InstanceType:        "n2d-standard-4",
							ConfidentialCompute: &confidentialComputeEnabled,
						}},
				},
			},
			wantErr: true,
		},
		{
			name: "GCPMachineTemplate with ConfidentialCompute enabled and unsupported instance type - invalid",
			template: &GCPMachineTemplate{
				Spec: GCPMachineTemplateSpec{
					Template: GCPMachineTemplateResource{
						Spec: GCPMachineSpec{
							InstanceType:        "e2-standard-4",
							ConfidentialCompute: &confidentialComputeEnabled,
							OnHostMaintenance:   &onHostMaintenanceTerminate,
						}},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			err := test.template.ValidateCreate()
			if test.wantErr {
				g.Expect(err).To(HaveOccurred())
			} else {
				g.Expect(err).NotTo(HaveOccurred())
			}
		})
	}
}
