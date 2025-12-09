package types

import "sync"

type Stopper struct {
    stopReceiver chan struct{} 
    childs       []chan struct{}
    childsMutex  sync.Mutex
    wg           sync.WaitGroup
}

func NewStopper() *Stopper {
    return &Stopper{
        stopReceiver: make(chan struct{}), 
        childs:       make([]chan struct{}, 0),
    }
}

// remove chan from childs
func removeChild(channels []chan struct{}, target chan struct{}) []chan struct{} {
    for i, ch := range channels {
        if ch == target {
            channels[i] = channels[len(channels)-1]
            return channels[:len(channels)-1]
        }
    }
    return channels 
}

func (s *Stopper) Go(f func(stopper *Stopper)) {
    childStopper := NewStopper()
    childCh := childStopper.stopReceiver

    // add child
    s.childsMutex.Lock()
    s.childs = append(s.childs, childCh)
    s.childsMutex.Unlock()
    
    s.wg.Add(1)
    go func() {
        defer func () {
			// remove child
            s.childsMutex.Lock()
            s.childs = removeChild(s.childs, childCh)
            s.childsMutex.Unlock()
            
            s.wg.Done()
        }()

        f(childStopper)
    }()
}

func (s *Stopper) WaitForStopRequest() chan struct{} {
    return s.stopReceiver
}

func (s *Stopper) StopChilds() {
    s.childsMutex.Lock()
	defer s.childsMutex.Unlock()

    for _, channel := range s.childs {
        channel <- struct{}{}
    }

    s.wg.Wait()
}