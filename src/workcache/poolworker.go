// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	//"fmt"
	"log"
	"time"
)

// Background generate jobs, with low priority.  Size is large.
var pregenerateJobsMaxSize int = 0
var pregenerateJobs chan WorkRequest = nil

func InitQueue() {
	pregenerateJobsMaxSize = ConfigPregenerationQueueSize()
	pregenerateJobs = make(chan WorkRequest, pregenerateJobsMaxSize)
}

func addPregenerateRequest(req WorkRequest) {
	if len(pregenerateJobs) >= pregenerateJobsMaxSize-2 {
		// queue is full, do not put any more (to avoid blocking)
		log.Printf("WARNING: Pregeneration queue is full, not enqueuing any more, %v\n", len(pregenerateJobs))
		return
	}
	pregenerateJobs <- req
}

func doProcess(name int) {
	for {
		// wait on queue, with periodical timeout
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		select {
		case preJob := <-pregenerateJobs:
			//log.Printf("Worker %v : pregenerate job", name)
			resp := getCachedWorkByAccountOrHash(preJob)
			if resp.Error != nil {
				log.Printf("WARNING: Could not process request, sleeping to slow queue, %v \n", resp.Error)
				time.Sleep(20 * time.Second)
			}
		case <-ticker.C:
			// timeout, idle loop
		}
	}
}

func startWorkers(backgroundWorkerCount int) {
	for i := 0; i < backgroundWorkerCount; i++ {
		go doProcess(i)
	}
	log.Printf("%v pool workers started\n", backgroundWorkerCount)
}

func StatusPregenerQueueSize() int { return len(pregenerateJobs) }
