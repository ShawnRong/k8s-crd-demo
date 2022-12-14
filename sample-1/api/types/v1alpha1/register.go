package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "k8s-crd-demo.shawnrong.github.com"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{
	Group:   GroupName,
	Version: GroupVersion,
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnowTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnowTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Project{}, &ProjectList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
