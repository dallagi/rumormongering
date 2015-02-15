package main

import (
	"math/rand"
	"sync"
	"sync/atomic"
)

type Request struct {
	Value  Update
	Answer chan bool
}

type Transport interface {
	// Tick is called by the node on every tick
	Tick()
	// Push selects a random peer and pushes a message to it
	Push(u Update) (bool, error)
	// Requests returns the channel where Request object are sent when they reach
	// the transport
	Requests() <-chan Request

	// Stop stops the transport
	Stop()
}

type SimulationTransport struct {
	nodeID    int
	nodes     []*Node
	nodeChans []chan Request

	packets uint64

	tick         uint64
	firstReceive map[Update]uint64
	mu           sync.Mutex
}

func (s *SimulationTransport) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tick++
}

func (s *SimulationTransport) randomPeer() chan Request {
	rnd := s.nodeID
	// get a random node until it's different from myself
	for rnd == s.nodeID {
		rnd = rand.Intn(len(s.nodes))
	}
	return s.nodeChans[rnd]
}

func (s *SimulationTransport) Push(u Update) (bool, error) {
	atomic.AddUint64(&s.packets, 1)
	peer := s.randomPeer()
	answer := make(chan bool)
	peer <- Request{
		Value:  u,
		Answer: answer,
	}
	return <-answer, nil
}

func (s *SimulationTransport) Requests() <-chan Request {
	ch := s.nodeChans[s.nodeID]
	outCh := make(chan Request)
	go func() {
		for request := range ch {
			s.mu.Lock()
			if _, ok := s.firstReceive[request.Value]; !ok {
				s.firstReceive[request.Value] = s.tick
			}
			s.mu.Unlock()
			outCh <- request
		}
	}()
	return outCh
}

func (s *SimulationTransport) Stop() {
	close(s.nodeChans[s.nodeID])
}
