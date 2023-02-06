package dtrack

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFetchAll(t *testing.T) {
	var wantItems []int
	for i := 0; i < 468; i++ {
		wantItems = append(wantItems, i)
	}
	gotItems, err := FetchAll(func(po PageOptions) (p Page[int], err error) {
		for i := 0; i < po.PageSize; i++ {
			idx := (po.PageSize * (po.PageNumber - 1)) + i
			if idx >= len(wantItems) {
				break
			}
			p.Items = append(p.Items, wantItems[idx])
		}
		p.TotalCount = len(wantItems)
		return p, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if diff := cmp.Diff(wantItems, gotItems); diff != "" {
		t.Errorf("unexpected items:\n%s", diff)
	}
}

func TestFetchAll_PageFetchFuncErr(t *testing.T) {
	var testErr = errors.New("test error")
	if _, err := FetchAll(
		func(po PageOptions) (p Page[int], err error) {
			return p, testErr
		},
	); !errors.Is(err, testErr) {
		t.Errorf("expected err but got nil")
	}
}

func TestForEach(t *testing.T) {
	var (
		wantItems []int
		gotItems  []int
	)
	for i := 0; i < 468; i++ {
		wantItems = append(wantItems, i)
	}
	if err := ForEach(
		func(po PageOptions) (p Page[int], err error) {
			for i := 0; i < po.PageSize; i++ {
				idx := (po.PageSize * (po.PageNumber - 1)) + i
				if idx >= len(wantItems) {
					break
				}
				p.Items = append(p.Items, wantItems[idx])
			}
			p.TotalCount = len(wantItems)
			return p, nil
		},
		func(item int) error {
			gotItems = append(gotItems, item)
			return nil
		},
	); err != nil {
		t.Errorf("unexpected error calling ForEach: %s", err)
	}
	if diff := cmp.Diff(wantItems, gotItems); diff != "" {
		t.Errorf("unexpected items:\n%s", diff)
	}
}

func TestForEach_PageFetchFuncErr(t *testing.T) {
	var testErr = errors.New("test error")
	if err := ForEach(
		func(po PageOptions) (p Page[int], err error) {
			return p, testErr
		},
		func(item int) error {
			return nil
		},
	); !errors.Is(err, testErr) {
		t.Errorf("expected error from pageFetchFunc but got nil")
	}
}

func TestForEach_HandlerFuncErr(t *testing.T) {
	var testErr = errors.New("test error")
	if err := ForEach(
		func(po PageOptions) (p Page[int], err error) {
			p.Items = []int{0, 1, 2, 3}
			p.TotalCount = len(p.Items)
			return p, nil
		},
		func(item int) error {
			return testErr
		},
	); !errors.Is(err, testErr) {
		t.Errorf("expected error from handlerFunc but got nil")
	}
}
