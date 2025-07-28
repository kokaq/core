package tests

import (
	"bytes"
	"testing"

	"github.com/kokaq/core/queue"
)

func TestNewPage(t *testing.T) {
	p := queue.NewPage()
	if p == nil {
		t.Fatal("NewPage returned nil")
	}
	if p.GetIndex() != 0 {
		t.Errorf("expected index 0, got %d", p.GetIndex())
	}
	if p.GetData() != nil {
		t.Errorf("expected data nil, got %v", p.GetData())
	}
}

func TestSetData(t *testing.T) {
	p := queue.NewPage()
	data := []byte{1, 2, 3}
	p.SetData(5, data)

	if p.GetIndex() != 5 {
		t.Errorf("expected index 5, got %d", p.GetIndex())
	}
	if !bytes.Equal(p.GetData(), data) {
		t.Errorf("expected data %v, got %v", data, p.GetData())
	}
}

func TestSetDataOverwrite(t *testing.T) {
	p := queue.NewPage()
	data1 := []byte{1, 2}
	data2 := []byte{3, 4, 5}
	p.SetData(1, data1)
	p.SetData(2, data2)

	if p.GetIndex() != 2 {
		t.Errorf("expected index 2, got %d", p.GetIndex())
	}
	if !bytes.Equal(p.GetData(), data2) {
		t.Errorf("expected data %v, got %v", data2, p.GetData())
	}
}

func TestSetDataNil(t *testing.T) {
	p := queue.NewPage()
	p.SetData(10, nil)
	if p.GetIndex() != 10 {
		t.Errorf("expected index 10, got %d", p.GetIndex())
	}
	if p.GetData() != nil {
		t.Errorf("expected data nil, got %v", p.GetData())
	}
}
