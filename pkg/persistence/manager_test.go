package persistence

import (
	"testing"
	"time"
)

func TestMemoryItemManager_Add(t *testing.T) {
	type args struct {
		item Item
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "Default", args: args{item: Item{
			Code:    "a-code",
			Content: "Content",
			Created: time.Now(),
			Timeout: 1 * time.Millisecond,
		}}},
		{name: "Auto timeout", args: args{item: Item{
			Code:    "a-code",
			Content: "Content",
			Created: time.Now(),
			Timeout: 100 * time.Millisecond,
		}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewMemoryItemManager()

			if _, ok := manager.Get(tt.args.item.Code); ok {
				t.Errorf("Add() got = %v, expected %v", ok, false)
			}

			manager.Add(tt.args.item)

			if it, ok := manager.Get(tt.args.item.Code); !ok || it.Code != tt.args.item.Code {
				t.Errorf("Add() got = %v, expected %v", it, tt.args.item)
			}

			time.Sleep(2 * tt.args.item.Timeout)

			if it, ok := manager.Get(tt.args.item.Code); ok {
				t.Errorf("Add() got = %v, should have time outed", it)
			}
		})
	}
}

func TestMemoryItemManager_Close(t *testing.T) {
	t.Run("Close manager", func(t *testing.T) {
		manager := NewMemoryItemManager()
		manager.Close()
		item := Item{
			Code: "a-code",
			Timeout: 10 * time.Millisecond,
		}
		manager.Add(item)
		time.Sleep(150 * time.Millisecond)

		if length := len(manager.items); length != 1 {
			t.Errorf("Close() got = %d items, should have received %d", length, 0)
		}
	})
}


func TestMemoryItemManager_Remove(t *testing.T) {
	t.Run("Remove item", func(t *testing.T) {
		manager := NewMemoryItemManager()
		item := Item{
			Code:    "a-code",
			Timeout: 10 * time.Millisecond,
		}
		manager.Add(item)
		manager.Remove(item.Code)

		if got, ok := manager.Get(item.Code); ok {
			t.Errorf("Remove() got = %v, should have received nil", got)
		}
	})
}