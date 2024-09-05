package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func authToClusterLocal() (*kubernetes.Clientset, error) {
	home, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfigOrDie(config), nil
}

func authToClusterFromCluster() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfigOrDie(config), nil
}

func authToCluster() *kubernetes.Clientset {
	// try from cluster
	client, err := authToClusterFromCluster()
	if err != nil {
		log.Println("Cant auth from cluster trying auth from local")
		client, err = authToClusterLocal()
		if err != nil {
			log.Fatalf("Failed to get Kubernetes client: %v", err)
		}
	}
	return client
}

func proxyToAllAPods(w http.ResponseWriter, r *http.Request, client *kubernetes.Clientset, config config) {
	ctx := context.TODO()
	opts := metav1.ListOptions{LabelSelector: config.serviceLabelSelector}
	pods, err := client.CoreV1().Pods(config.serviceNamespace).List(ctx, opts)
	if err != nil {
		log.Printf("Failed to list pods in namespace '%s': %v\n", config.serviceNamespace, err)
		http.Error(w, "Failed to list pods", http.StatusInternalServerError)
		return
	}

	log.Printf("Found %d pods in namespace '%s'\n", len(pods.Items), config.serviceNamespace)

	responses := make(chan string, len(pods.Items)) // Channel to collect responses
	var wg sync.WaitGroup

	for _, pod := range pods.Items {
		wg.Add(1)
		go func(podIP string) {
			defer wg.Done()
			targetURL := fmt.Sprintf("http://%s:%d%s", podIP, config.servicePort, r.URL.String())
			log.Printf("Sending request to %s\n", targetURL)
			resp, err := http.Get(targetURL)
			if err != nil {
				log.Printf("Error contacting pod %s: %v\n", podIP, err)
				responses <- fmt.Sprintf("Error contacting pod %s: %v\n", podIP, err)
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response from pod %s: %v\n", podIP, err)
				responses <- fmt.Sprintf("Error reading response from pod %s: %v\n", podIP, err)
				return
			}
			response := fmt.Sprintf("Response from pod %s (%s): %s", pod.Name, podIP, body)
			responses <- response
		}(pod.Status.PodIP)
	}

	wg.Wait()
	close(responses)

	// Write all responses back to the client
	for response := range responses {
		fmt.Fprintln(w, response)
	}
}

func test(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "OK")
}

func main() {
	client := authToCluster()
	config := initConfig()
	http.HandleFunc("/test", test)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { proxyToAllAPods(w, r, client, config) })
	log.Println("Starting proxy to all with port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
