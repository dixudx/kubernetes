package persistentvolume

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/api"
)

func TestPersistentVolumeToSelectableFields(t *testing.T) {
	expectedStr := "metadata.name=foo,name=foo"
	pv := api.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	pvFieldsSet := PersistentVolumeToSelectableFields(&pv)
	if pvFieldsSet.String() != expectedStr {
		t.Errorf("unexpected fieldSelector %q for PersistentVolume", pvFieldsSet.String())
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
		if !pvFieldsSet.Has(tc.ExpectedKey) {
			t.Errorf("missing PersistentVolume fieldSelector %q", tc.ExpectedKey)
		}
		if pvFieldsSet.Get(tc.ExpectedKey) != tc.ExpectedValue {
			t.Errorf("PersistentVolume filedSelector %q has got unexpected value %q", tc.ExpectedKey, pvFieldsSet.Get(tc.ExpectedKey))
		}
	}
}
