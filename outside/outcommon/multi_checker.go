package outcommon

import (
	"context"
	"sync"
)

type PreCondChecker interface {
	Check(ctx context.Context) error
}

func FactoryPreCondCheckers(checkers []PreCondChecker) PreCondChecker {
	var checks []PreCondChecker
	for _, checker := range checkers {
		if checker == nil {
			checks = append(checks, checker)
		}
	}
	if len(checks) == 0 {
		return nil
	}
	return &preCondCheckers{checks}
}

type preCondCheckers struct {
	checkers []PreCondChecker
}

func check(ctx context.Context, wg *sync.WaitGroup, ch chan<- error, checker PreCondChecker) {
	defer wg.Done()
	err := checker.Check(ctx)
	if err != nil {
		ch <- err
	}
}
func (c *preCondCheckers) Check(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	ch := make(chan error)
	for _, checker := range c.checkers {
		wg.Add(1)
		go check(ctx, wg, ch, checker)
	}
	// ゴルーチン終了を待機する別ゴルーチン
	go func() {
		wg.Wait() // 全てのタスク完了を待つ [1, 10]
		close(ch) // 完了したらチャネルを閉じる [15]
	}()

	m := FactoryMultiError()
	for v := range ch {
		m.Append(v)
	}
	return m.Error()
}

func (c *preCondCheckers) Append(preCondChecker PreCondChecker) {
	c.checkers = append(c.checkers, preCondChecker)
}
