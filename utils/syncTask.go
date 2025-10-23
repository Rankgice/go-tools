package utils

import "sync"

type TaskFunc func() error

// SyncTask 执行异步任务
func SyncTask(taskFuncList ...TaskFunc) *SyncTaskType {
	t := SyncTaskType{&sync.WaitGroup{}, false, make(chan error, 1), &sync.Mutex{}}
	t.AddSyncTask(taskFuncList...)
	return &t
}

type SyncTaskType struct {
	wg       *sync.WaitGroup
	hasError bool
	errChan  chan error
	mu       *sync.Mutex
}

// Wait 等待异步任务完成
func (s *SyncTaskType) Wait() error {
	waitChan := make(chan error, 1)
	go func() {
		s.wg.Wait()
		waitChan <- nil
		close(waitChan)
	}()
	select {
	case err := <-s.errChan:
		return err
	case <-waitChan:
		s.mu.Lock()
		defer s.mu.Unlock()
		if s.hasError {
			return <-s.errChan
		}
		return nil
	}
}

// AddSyncTask 继续添加异步任务
func (s *SyncTaskType) AddSyncTask(taskFuncList ...TaskFunc) {
	s.wg.Add(len(taskFuncList))
	for _, task := range taskFuncList {
		go func(wg *sync.WaitGroup, tf TaskFunc) {
			defer wg.Done()
			err := tf()
			if err != nil {
				s.mu.Lock()
				defer s.mu.Unlock()
				if !s.hasError {
					s.hasError = true
					s.errChan <- err
					close(s.errChan)
				}
			}
		}(s.wg, task)
	}
}
