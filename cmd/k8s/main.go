package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	// "k8s.io/client-go/rest"
)

type PageVariables struct {
	Date string
	Time string
}

const (
	defaultDateFormat = "02-01-2006"
	defaultTimeFormat = "16:04:05"
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

	now := time.Now()
	homePageVars := PageVariables{
		Date: now.Format(defaultDateFormat),
		Time: now.Format(defaultTimeFormat),
	}
	t, err := template.ParseFiles("ui/public/index.html")
	if err != nil {
		log.Fatal("template parsing error: ", err)
	}

	err = t.Execute(w, homePageVars)
	if err != nil {
		log.Fatal("template executing error: ", err)
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

	fmt.Printf("Starting the server now on http://localhost:8080 \n")
	http.HandleFunc("/", cronjobs)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
