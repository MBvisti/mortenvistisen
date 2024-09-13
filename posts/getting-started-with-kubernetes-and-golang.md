# How to Get Status and Version of Kubernetes Nodes Using Golang

After years of wrestling with complex front-end frameworks, I've found peace in the simplicity of Go. It's a language that doesn't try to be clever. It just works. Today, we're going to use Go to talk to Kubernetes. We'll learn how to get the status of Kubernetes nodes and their versions. This is useful stuff if you're building tools to manage clusters.

## Getting Started

First, we need some tools. Install Go if you haven't already. You'll also need a Kubernetes cluster. If you don't have one, set up minikube on your machine. It's a simple way to get a cluster running locally.

We'll use the Kubernetes Go client library. It's called client-go. It's the official way to talk to Kubernetes from Go code. Here's how we set it up:

```go
go mod init k8s-node-status
go get k8s.io/client-go@latest
```

This creates a new Go module and adds client-go to it.

## Creating a Kubernetes Client

To talk to Kubernetes, we need a client. The client is our messenger. It carries our questions to the cluster and brings back answers. Here's how we create one:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Get the kubeconfig file path
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeconfig := filepath.Join(home, ".kube", "config")

    // Build the config from the kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        fmt.Printf("Error building kubeconfig: %v\n", err)
        os.Exit(1)
    }

    // Create the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating clientset: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Successfully created Kubernetes client")
}
```

This code does a few things. It finds your kubeconfig file. That file tells the client how to connect to your cluster. Then it creates a config object from that file. Finally, it uses that config to create a clientset. The clientset is our main tool for talking to Kubernetes.

## Retrieving Kubernetes Node Status

Now that we have a client, let's use it. We'll start by getting the status of our nodes. The status tells us if a node is ready to run pods. Here's how we do it:

```go
import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getNodeStatus(clientset *kubernetes.Clientset) {
    // List all nodes
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing nodes: %v\n", err)
        return
    }

    // Print status of each node
    for _, node := range nodes.Items {
        fmt.Printf("Node: %s\n", node.Name)
        for _, condition := range node.Status.Conditions {
            if condition.Type == "Ready" {
                fmt.Printf("  Status: %v\n", condition.Status)
                fmt.Printf("  Last Heartbeat Time: %v\n", condition.LastHeartbeatTime)
            }
        }
    }
}
```

This function does the heavy lifting. It asks Kubernetes for a list of all nodes. Then it loops through each node. For each node, it finds the "Ready" condition. This condition tells us if the node is ready to accept pods. It prints the status and the last time the node checked in.

## Getting Kubernetes Node Version

Knowing the version of your nodes is important. It helps with upgrades and troubleshooting. Here's how we can get that information:

```go
func getNodeVersion(clientset *kubernetes.Clientset) {
    // List all nodes
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing nodes: %v\n", err)
        return
    }

    // Print version of each node
    for _, node := range nodes.Items {
        fmt.Printf("Node: %s\n", node.Name)
        fmt.Printf("  Kubernetes Version: %s\n", node.Status.NodeInfo.KubeletVersion)
    }
}
```

This function is similar to the last one. It gets a list of all nodes. But instead of looking at the status, it looks at the NodeInfo. The NodeInfo contains the version of Kubernetes running on that node.

## Putting It All Together

Now let's use these functions in our main program. Here's how our complete program looks:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // Get the kubeconfig file path
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error getting user home dir: %v\n", err)
        os.Exit(1)
    }
    kubeconfig := filepath.Join(home, ".kube", "config")

    // Build the config from the kubeconfig file
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        fmt.Printf("Error building kubeconfig: %v\n", err)
        os.Exit(1)
    }

    // Create the clientset
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        fmt.Printf("Error creating clientset: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("Node Status:")
    getNodeStatus(clientset)

    fmt.Println("\nNode Versions:")
    getNodeVersion(clientset)
}

func getNodeStatus(clientset *kubernetes.Clientset) {
    // List all nodes
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing nodes: %v\n", err)
        return
    }

    // Print status of each node
    for _, node := range nodes.Items {
        fmt.Printf("Node: %s\n", node.Name)
        for _, condition := range node.Status.Conditions {
            if condition.Type == "Ready" {
                fmt.Printf("  Status: %v\n", condition.Status)
                fmt.Printf("  Last Heartbeat Time: %v\n", condition.LastHeartbeatTime)
            }
        }
    }
}

func getNodeVersion(clientset *kubernetes.Clientset) {
    // List all nodes
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Error listing nodes: %v\n", err)
        return
    }

    // Print version of each node
    for _, node := range nodes.Items {
        fmt.Printf("Node: %s\n", node.Name)
        fmt.Printf("  Kubernetes Version: %s\n", node.Status.NodeInfo.KubeletVersion)
    }
}
```

This program creates a Kubernetes client. Then it uses that client to get the status and version of all nodes in the cluster. It's simple, but powerful.

## Running the Code

To run this code, save it as `main.go`. Then run:

```
go run main.go
```

You should see output like this:

```
Node Status:
Node: minikube
  Status: True
  Last Heartbeat Time: 2023-09-13T12:34:56Z

Node Versions:
Node: minikube
  Kubernetes Version: v1.26.3
```

This output tells us that we have one node named "minikube". It's ready to accept pods. The last time it checked in was at 12:34:56 on September 13, 2023. It's running Kubernetes version 1.26.3.

## Advanced Topics

Our code works, but it's basic. In a real-world scenario, you might want to do more. Here are some ideas:

1. **Using context for timeouts**: Right now, our code will wait forever for a response from Kubernetes. That's not ideal. We can use context to set a timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
nodes, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
```

This code will wait up to 10 seconds for a response. If it doesn't get one, it will give up.

2. **Implementing pagination**: If you have a large cluster, you might not want to get all nodes at once. You can use pagination:

```go
continueToken := ""
for {
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
        Limit:    100,
        Continue: continueToken,
    })
    if err != nil {
        fmt.Printf("Error listing nodes: %v\n", err)
        return
    }

    // Process nodes...

    continueToken = nodes.Continue
    if continueToken == "" {
        break
    }
}
```

This code gets nodes in batches of 100. It keeps going until there are no more nodes.

3. **Adding filters**: Maybe you only want to see nodes that aren't ready. You can use field selectors for that:

```go
nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{
    FieldSelector: "spec.unschedulable=true",
})
```

This code only returns nodes that are marked as unschedulable.

## Wrapping Up

We've learned how to get the status of Kubernetes nodes using Golang. We've also seen how to get the version of nodes in Kubernetes using Golang. These are powerful tools. With them, you can build applications that monitor and manage your Kubernetes clusters.

Remember, the Kubernetes API is vast. We've only scratched the surface. But the patterns we've used here apply to other resources too. You can use similar code to work with pods, services, and more.

Kubernetes is complex. But with Go and client-go, we can tame that complexity. We can build tools that make our lives easier. And isn't that what programming is all about?
