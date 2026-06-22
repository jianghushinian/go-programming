package watcher

import (
	"k8s.io/client-go/kubernetes"

	"cronx/internal/pkg/store"
)

type Config struct {
	Store store.IStore

	Clientset kubernetes.Interface
}
