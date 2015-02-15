package main

import (
	"log"
	"sync"
	"time"
)

type Node struct {
	ID       int
	DeltaT   time.Duration
	Strategy LossOfInterest

	Updates   []Update
	updatesMu sync.RWMutex

	wg   *sync.WaitGroup
	stop chan struct{}

	transport Transport
}

func (n *Node) hasUpdate(u Update) bool {
	n.updatesMu.RLock()
	defer n.updatesMu.RUnlock()
	for _, update := range n.Updates {
		if update == u {
			return true
		}
	}

	return false
}

func (n *Node) AddUpdate(u Update) {
	n.updatesMu.Lock()
	defer n.updatesMu.Unlock()
	n.Updates = append(n.Updates, u)
}

func (n *Node) MapUpdates(f func(u Update)) {
	n.updatesMu.RLock()
	defer n.updatesMu.RUnlock()
	for _, u := range n.Updates {
		f(u)
	}
}

func (n *Node) Stop() {
	n.transport.Stop()
	close(n.stop)
}

func (n *Node) Run() {
	n.wg.Add(1)
	defer n.wg.Done()

	t := time.NewTicker(n.DeltaT)
	defer t.Stop()

	reqCh := n.transport.Requests()

	for {
		select {
		case <-t.C:
			n.transport.Tick()
			n.MapUpdates(func(u Update) {
				if n.Strategy.Status(u) != Infective {
					return
				}
				go func() {
					hasRumor, err := n.transport.Push(u)
					if err != nil {
						log.Println(err)
						return
					}
					n.Strategy.Feedback(u, hasRumor)
				}()
			})
		case req := <-reqCh:
			//log.Printf("[%d] got update %s\n", n.ID, req.Value)
			has := n.hasUpdate(req.Value)

			// send feedback
			req.Answer <- has
			if !has {
				n.AddUpdate(req.Value)
			}
		case <-n.stop:
			return
		}
	}
}
