// Copyright (c) 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package v1alpha1

import (
	"fmt"

	"github.com/gardener/gardener/pkg/apis/seedmanagement"
	"github.com/gardener/gardener/pkg/apis/seedmanagement/helper"
	configv1alpha1 "github.com/gardener/gardener/pkg/gardenlet/apis/config/v1alpha1"

	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	if err := scheme.AddFieldLabelConversionFunc(
		SchemeGroupVersion.WithKind("ManagedSeed"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace", seedmanagement.ManagedSeedShootName:
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported: %s", label)
			}
		},
	); err != nil {
		return err
	}

	return nil
}

func Convert_v1alpha1_Gardenlet_To_seedmanagement_Gardenlet(in *Gardenlet, out *seedmanagement.Gardenlet, s conversion.Scope) error {
	if in.Config.Object == nil {
		cfg, err := helper.DecodeGardenletConfigurationFromBytes(in.Config.Raw, false)
		if err != nil {
			return err
		}
		in.Config.Object = cfg
	}
	return autoConvert_v1alpha1_Gardenlet_To_seedmanagement_Gardenlet(in, out, s)
}

func Convert_seedmanagement_Gardenlet_To_v1alpha1_Gardenlet(in *seedmanagement.Gardenlet, out *Gardenlet, s conversion.Scope) error {
	if err := autoConvert_seedmanagement_Gardenlet_To_v1alpha1_Gardenlet(in, out, s); err != nil {
		return err
	}
	if out.Config.Raw == nil {
		cfg, ok := out.Config.Object.(*configv1alpha1.GardenletConfiguration)
		if !ok {
			return fmt.Errorf("unknown gardenlet config object type")
		}
		raw, err := helper.EncodeGardenletConfigurationToBytes(cfg)
		if err != nil {
			return err
		}
		out.Config.Raw = raw
	}
	return nil
}
