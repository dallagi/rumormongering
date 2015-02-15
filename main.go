package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	peers = flag.Uint("peers", 5, "number of peers")

	k        = flag.Uint("k", 2, "k parameter")
	strategy = flag.String("strategy", "cf", "cf (counter feedback) or br (blind random)")

	deltaT = flag.Duration("delta-t", 100*time.Millisecond, "how often gossip messages are sent")

	tries = flag.Int("tries", 5, "number of times to run the simulation and average the results")
)

func main() {
	var (
		c = csv.NewWriter(os.Stdout)
	)
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	t := *tries

	c.Write([]string{"strategy", "deltaT", "k", "run", "residue", "traffic", "time_avg", "time_max"})
	for i := 0; i < t; i++ {
		env := NewTestEnv(*peers, *strategy, *k, *deltaT)
		r, t, ta, tm := env.GetStats()
		c.Write([]string{
			*strategy,
			deltaT.String(),
			fmt.Sprintf("%d", *k),
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%.4f", r),
			fmt.Sprintf("%.4f", t),
			fmt.Sprintf("%.4f", ta),
			fmt.Sprintf("%.4f", tm),
		})
		c.Flush()
		env.Stop()
	}
	err := c.Error()
	if err != nil {
		log.Println(err)
	}
}
