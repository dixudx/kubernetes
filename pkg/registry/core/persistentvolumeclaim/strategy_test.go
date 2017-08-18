/*
Copyright 2015 The Kubernetes Authors.

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

package persistentvolumeclaim

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"
)

func TestPersistentVolumeClaimToSelectableFields(t *testing.T) {
	expectedStr := "metadata.name=foo,name=foo"
	pvc := api.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	pvcFieldsSet := PersistentVolumeClaimToSelectableFields(&pvc)
	if pvcFieldsSet.String() != expectedStr {
		t.Errorf("unexpected fieldSelector %q for PersistentVolumeClaim", pvcFieldsSet.String())
	}

	testcases := []struct {
		ExpectedKey   string
		ExpectedValue string
	}{
		{
			ExpectedKey:   "metadata.name",
			ExpectedValue: "foo",
		},
		{
			ExpectedKey:   "name",
			ExpectedValue: "foo",
		},
	}

	for _, tc := range testcases {
		if !pvcFieldsSet.Has(tc.ExpectedKey) {
			t.Errorf("missing PersistentVolumeClaim fieldSelector %q", tc.ExpectedKey)
		}
		if pvcFieldsSet.Get(tc.ExpectedKey) != tc.ExpectedValue {
			t.Errorf("PersistentVolumeClaim filedSelector %q has got unexpected value %q", tc.ExpectedKey, pvcFieldsSet.Get(tc.ExpectedKey))
		}
	}
}
