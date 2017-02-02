package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"unicode"

	"github.com/ghodss/yaml"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/batch"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {

	config := &restclient.Config{
		Host:     "http://localhost:8001",
		Username: "admin",
		Password: "ZG2CZa6Mp3DbnDij",
		Insecure: true,
	}

	err := restclient.SetKubernetesDefaults(config)
	if err != nil {
		fmt.Println(err)
	}

	if err = restclient.LoadTLSFiles(config); err != nil {
		fmt.Println(err)
	}

	conn, err := client.New(config)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	pods, err := conn.Pods(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	fmt.Printf("Number of pods: %v \n", len(pods.Items))

	for _, pod := range pods.Items {
		fmt.Println(pod.ObjectMeta.Name)
	}

	scheduler := conn.Batch().Jobs(api.NamespaceDefault)

	// Read in a file and convert it from YAML to JSON, then schedule it!
	yamlJob, err := ioutil.ReadFile("batch.yaml")
	if err != nil {
		fmt.Println(err)
	}
	jsonJob, err := ToJSON(yamlJob)
	if err != nil {
		fmt.Println(err)
	}

	job := batch.Job{}
	if err = json.Unmarshal(jsonJob, &job); err != nil {
		fmt.Println(err)
	}
	fmt.Printf("API Version: %v \n", job.TypeMeta.APIVersion)
	fmt.Printf("Kind: %v \n", job.TypeMeta.Kind)
	fmt.Println(string(jsonJob))

	_, err = scheduler.Create(&job)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("vim-go")

}

// ToJSON converts a single YAML document into a JSON document
// or returns an error. If the document appears to be JSON the
// YAML decoding path is not used.
func ToJSON(data []byte) ([]byte, error) {
	if hasJSONPrefix(data) {
		return data, nil
	}
	return yaml.YAMLToJSON(data)
}

// hasJSONPrefix returns true if the provided buffer appears to start with
// a JSON open brace.
func hasJSONPrefix(buf []byte) bool {
	var jsonPrefix = []byte("{")
	return hasPrefix(buf, jsonPrefix)
}

// Return true if the first non-whitespace bytes in buf is prefix.
func hasPrefix(buf []byte, prefix []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	return bytes.HasPrefix(trim, prefix)
}
