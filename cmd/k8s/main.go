package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// "k8s.io/client-go/rest"
)

func cronjobs(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Cronjobs \n")

	clientset := getClientSet()
	cronjobsClient := clientset.BatchV1beta1().CronJobs(apiv1.NamespaceAll)

	// List all running cronjobs
	// fmt.Printf("Listing cronjobs in all namespaces \n")
	cronlist, err := cronjobsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range cronlist.Items {
		fmt.Fprintf(w, " * %s %s %s %d \n", d.Name, d.Spec.Schedule, d.Status.LastScheduleTime, *d.Spec.FailedJobsHistoryLimit)
	}
}

func getClientSet() *kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}

func main() {

	fmt.Printf("Starting the server now on http://localhost:8090 \n")
	http.HandleFunc("/", cronjobs)

	http.ListenAndServe(":8090", nil)

}
