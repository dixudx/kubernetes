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

package secret

import (
	"reflect"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/kubernetes/pkg/api"

	// install all api groups for testing
	_ "k8s.io/kubernetes/pkg/api/testapi"
)

func TestExportSecret(t *testing.T) {
	tests := []struct {
		objIn     runtime.Object
		objOut    runtime.Object
		exact     bool
		expectErr bool
	}{
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			objOut: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			},
			exact: true,
		},
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
				},
				Type: api.SecretTypeServiceAccountToken,
			},
			expectErr: true,
		},
		{
			objIn: &api.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "bar",
					Annotations: map[string]string{
						api.ServiceAccountUIDKey: "true",
					},
				},
			},
			expectErr: true,
		},
		{
			objIn:     &api.Pod{},
			expectErr: true,
		},
	}

	for _, test := range tests {
		err := Strategy.Export(genericapirequest.NewContext(), test.objIn, test.exact)
		if err != nil {
			if !test.expectErr {
				t.Errorf("unexpected error: %v", err)
			}
			continue
		}
		if test.expectErr {
			t.Error("unexpected non-error")
			continue
		}
		if !reflect.DeepEqual(test.objIn, test.objOut) {
			t.Errorf("expected:\n%v\nsaw:\n%v\n", test.objOut, test.objIn)
		}
	}
}

func TestSecretToSelectableFields(t *testing.T) {
	expectedStr := "metadata.name=foo,type=type1"
	secret := api.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
		Type: api.SecretType("type1"),
	}

	secretFieldsSet := SecretToSelectableFields(&secret)
	if secretFieldsSet.String() != expectedStr {
		t.Errorf("unexpected fieldSelector %q for Secret", secretFieldsSet.String())
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
			ExpectedKey:   "type",
			ExpectedValue: "type1",
		},
	}

	for _, tc := range testcases {
		if !secretFieldsSet.Has(tc.ExpectedKey) {
			t.Errorf("missing Secret fieldSelector %q", tc.ExpectedKey)
		}
		if secretFieldsSet.Get(tc.ExpectedKey) != tc.ExpectedValue {
			t.Errorf("Secret filedSelector %q has got unexpected value %q", tc.ExpectedKey, secretFieldsSet.Get(tc.ExpectedKey))
		}
	}
}
