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

package node

import (
	"testing"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"

	// install all api groups for testing
	_ "k8s.io/kubernetes/pkg/api/testapi"
)

func TestMatchNode(t *testing.T) {
	testFieldMap := map[bool][]fields.Set{
		true: {
			{"metadata.name": "foo"},
		},
		false: {
			{"foo": "bar"},
		},
	}

	for expectedResult, fieldSet := range testFieldMap {
		for _, field := range fieldSet {
			m := MatchNode(labels.Everything(), field.AsSelector())
			_, matchesSingle := m.MatchesSingle()
			if e, a := expectedResult, matchesSingle; e != a {
				t.Errorf("%+v: expected %v, got %v", fieldSet, e, a)
			}
		}
	}
}

func TestNodeToSelectableFields(t *testing.T) {
	expectedStr := "metadata.name=foo,spec.unschedulable=false"
	node := api.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
		Spec: api.NodeSpec{
			Unschedulable: false,
		},
	}

	nodeFieldsSet := NodeToSelectableFields(&node)
	if nodeFieldsSet.String() != expectedStr {
		t.Errorf("unexpected fieldSelector %q for Node", nodeFieldsSet.String())
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
			ExpectedKey:   "spec.unschedulable",
			ExpectedValue: "false",
		},
	}

	for _, tc := range testcases {
		if !nodeFieldsSet.Has(tc.ExpectedKey) {
			t.Errorf("missing Node fieldSelector %q", tc.ExpectedKey)
		}
		if nodeFieldsSet.Get(tc.ExpectedKey) != tc.ExpectedValue {
			t.Errorf("Node filedSelector %q has got unexpected value %q", tc.ExpectedKey, nodeFieldsSet.Get(tc.ExpectedKey))
		}
	}
}
