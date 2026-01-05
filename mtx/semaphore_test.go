package mtx

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	ctxT := t.Context()
	ts := NewTimeoutSemaphore(6)

	wg := sync.WaitGroup{}

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctxT, time.Second*10)
			defer cancel()

			//nolint: gosec // only for test
			err := ts.AcquireWithAutoRelease(ctx, 2, time.Second*10, uint64(i), func(uint64) {
				t.Logf("%06d release timeout, auto release\n", i)
			})

			if err != nil {
				t.Logf("%06d acquire failed\n", i)

				return
			}

			t.Logf("%06d acquire success\n", i)
			time.Sleep(time.Second * time.Duration(rand.Int31n(20))) //nolint: gosec // only for test
			ts.Release(uint64(i), 2)                                 //nolint: gosec // only for test
			t.Logf("%06d release success\n", i)
		}(i)
	}

	wg.Wait()
}

func TestSemaphore2(t *testing.T) {
	ctxT := t.Context()
	ts := NewTimeoutSemaphore(6)

	wg := sync.WaitGroup{}

	for i := range 10 {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ctx, cancel := context.WithTimeout(ctxT, time.Second*10)
			defer cancel()

			// //nolint: gosec // only for test
			err := ts.AcquireWithAutoRelease(ctx, 2, time.Second*4, uint64(i), func(uint64) {
				t.Logf("%06d release timeout, auto release\n", i)
			})

			if err != nil {
				t.Logf("%06d acquire failed\n", i)

				return
			}

			t.Logf("%06d acquire success\n", i)
			time.Sleep(time.Second * 5)
			ts.Release(uint64(i), 2) //nolint: gosec // only for test
			t.Logf("%06d release success\n", i)
		}(i)
	}

	wg.Wait()
}
