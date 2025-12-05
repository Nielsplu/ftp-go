package server

import "sync"

type Stopper struct {
	stopReceiver chan struct{} 
	childs       []chan struct{}
	wg           sync.WaitGroup
}

func NewStopper() *Stopper {
	return &Stopper{
		stopReceiver: make(chan struct{}, 1),
		childs: 	  make([]chan struct{}, 0),
	}
}

func (s *Stopper) Go(f func(stopper *Stopper)) {
	childStopper := NewStopper()
	s.childs = append(s.childs, childStopper.stopReceiver)
	
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f(childStopper)
	}()
}

func (s *Stopper) Wait() chan struct{} {
	return s.stopReceiver
}

func (s *Stopper) Stop() {
	for _, channel := range s.childs {
		channel<-struct{}{}
	}
	s.wg.Wait()
}