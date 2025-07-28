package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kokaq/core/queue"
)

func setupTestDir(t *testing.T) string {
	dir := filepath.Join(os.TempDir(), "namespace_test")
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("failed to cleanup test dir: %v", err)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}
	return dir
}

func TestNewNamespace_CreatesDirectory(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	config := queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 1}
	ns := queue.NewNamespace(dir, config)

	expectedDir := filepath.Join(dir, "testns-1")
	if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist", expectedDir)
	}
	if ns.Name != "testns" || ns.Id != 1 {
		t.Errorf("namespace fields not set correctly")
	}
}

func TestNamespace_AddAndGetQueue(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	config := queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 2}
	ns := queue.NewNamespace(dir, config)

	qConfig := &queue.QueueConfiguration{QueueId: 42, QueueName: "q42"}
	q, err := ns.AddQueue(qConfig)
	if err != nil {
		t.Fatalf("AddQueue failed: %v", err)
	}
	got, err := ns.GetQueue(42)
	if err != nil {
		t.Fatalf("GetQueue failed: %v", err)
	}
	if got != q {
		t.Errorf("GetQueue did not return the correct queue")
	}
}

func TestNamespace_GetQueue_NotFound(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	ns := queue.NewNamespace(dir, queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 3})
	_, err := ns.GetQueue(999)
	if err == nil {
		t.Errorf("expected error for non-existent queue")
	}
}

func TestNamespace_LoadQueue(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	ns := queue.NewNamespace(dir, queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 4})
	qConfig := &queue.QueueConfiguration{QueueId: 7, QueueName: "q7"}
	q1, err := ns.LoadQueue(qConfig)
	if err != nil {
		t.Fatalf("LoadQueue failed: %v", err)
	}
	q2, err := ns.LoadQueue(qConfig)
	if err != nil {
		t.Fatalf("LoadQueue failed: %v", err)
	}
	if q1 != q2 {
		t.Errorf("LoadQueue should return the same queue instance")
	}
}

func TestNamespace_ClearQueue(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	ns := queue.NewNamespace(dir, queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 5})
	qConfig := &queue.QueueConfiguration{QueueId: 8, QueueName: "q8"}
	q, _ := ns.AddQueue(qConfig)
	if err := q.Clear(); err != nil {
		t.Errorf("ClearQueue failed: %v", err)
	}
}

func TestNamespace_DeleteQueue(t *testing.T) {
	dir := setupTestDir(t)
	defer os.RemoveAll(dir)

	ns := queue.NewNamespace(dir, queue.NamespaceConfig{NamespaceName: "testns", NamespaceId: 6})
	qConfig := &queue.QueueConfiguration{QueueId: 9, QueueName: "q9"}
	ns.AddQueue(qConfig)
	if err := ns.DeleteQueue(9); err != nil {
		t.Errorf("DeleteQueue failed: %v", err)
	}
}
