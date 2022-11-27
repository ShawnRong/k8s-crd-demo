package utils

import (
	"bytes"
	"html/template"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	v1alpha1 "shawnrong.github.io/examplecontroller/api/v1alpha1"
)

func parseTemplate(templateName string, website *v1alpha1.Website) []byte {
	tmpl, err := template.ParseFiles("controllers/template/" + templateName + ".yml")
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, website)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewDeployment(website *v1alpha1.Website) *appsv1.Deployment {
	d := &appsv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", website), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewService(website *v1alpha1.Website) *corev1.Service {
	s := &corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", website), s)
	if err != nil {
		panic(err)
	}
	return s
}
