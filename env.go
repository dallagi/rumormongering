package main

import (
	"math/rand"
	"sync"
	"time"
)

type TestEnv struct {
	Nodes      []*Node
	transports []*SimulationTransport

	Messages int

	wg *sync.WaitGroup
}

func (t *TestEnv) inject(update string) {
	node := t.Nodes[rand.Intn(len(t.Nodes))]
	node.AddUpdate(Update(update))
}

func (t *TestEnv) residue(update string) int {
	count := 0
	toCount := Update(update)
	for _, n := range t.Nodes {
		count++
		n.MapUpdates(func(u Update) {
			if u == toCount {
				count--
			}
		})
	}
	return count
}

func (t *TestEnv) Stop() {
	for _, n := range t.Nodes {
		n.Stop()
	}

	t.wg.Wait()
}

func (t *TestEnv) isStopped() bool {
	infective := true
	for _, n := range t.Nodes {
		n.MapUpdates(func(u Update) {
			if n.Strategy.Status(u) == Infective {
				infective = false
			}
		})
	}

	return infective
}

func (t *TestEnv) packets() uint64 {
	var r uint64
	for _, t := range t.transports {
		r += t.packets
	}
	return r
}

func (t *TestEnv) times(update string) (float64, float64) {
	var (
		u           = Update(update)
		a           = new(avg)
		max float64 = -1
	)
	for _, t := range t.transports {
		t.mu.Lock()
		if v, ok := t.firstReceive[u]; ok {
			f := float64(v)
			a.Add(f)
			if f > max {
				max = f
			}
		}
		t.mu.Unlock()
	}
	return a.Get(), max
}

// GetStats returns residue, traffic, t_avg, and t_last for the current run of the environment
func (t *TestEnv) GetStats() (float64, float64, float64, float64) {
	const testUpdate = "test"
	t.inject(testUpdate)
	// TODO: think about isStopped
	for !t.isStopped() {
		time.Sleep(100 * time.Millisecond)
	}

	nodes := float64(len(t.Nodes))

	// residue is the number of nodes that did not receive the update
	// divided by the number of peers
	residue := float64(t.residue(testUpdate)) / nodes
	// traffic is the total number of updates sent divided by the number of peers
	traffic := float64(t.packets()) / nodes
	// tAvg, tMax are the mean and maximum tick time at which the update was received
	tAvg, tMax := t.times(testUpdate)
	return residue, traffic, tAvg, tMax
}

func pickStrategy(st string, k uint) LossOfInterest {
	switch st {
	case "cf":
		return NewCounterFeedback(k)
	default:
		return NewBlindRandom(k)
	}
}

func NewTestEnv(peers uint, strategy string, k uint, deltaT time.Duration) *TestEnv {
	var (
		wg = new(sync.WaitGroup)
		e  = &TestEnv{
			Nodes:      make([]*Node, 0, peers),
			transports: make([]*SimulationTransport, 0, peers),

			wg: wg,
		}
		nodeChans = make([]chan Request, 0, peers)
	)

	for i := 0; i < int(peers); i++ {
		newNode := &Node{
			ID:       i,
			DeltaT:   deltaT,
			Strategy: pickStrategy(strategy, k),

			Updates: make([]Update, 0),

			wg:   wg,
			stop: make(chan struct{}),
		}

		e.Nodes = append(e.Nodes, newNode)
		nodeChans = append(nodeChans, make(chan Request))
	}

	for i := 0; i < int(peers); i++ {
		node := e.Nodes[i]
		transport := &SimulationTransport{
			nodeID:    i,
			nodes:     e.Nodes,
			nodeChans: nodeChans,

			firstReceive: make(map[Update]uint64),
		}
		e.transports = append(e.transports, transport)
		node.transport = transport
		go node.Run()
	}
	return e
}
