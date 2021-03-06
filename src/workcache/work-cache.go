// Copyright © 2019-2020 catenocrypt.  See LICENSE file for license information.

package workcache

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/catenocrypt/nano-work-cache/rpcclient"
)

type CacheEntry struct {
	hash       string
	work       string
	difficulty uint64
	multiplier float64
	account    string
	// valid, computing
	status       string
	timeComputed int64 // unix time
	timeAdded    int64
}

var (
	// The cache, key is hash
	workCache map[string]CacheEntry = map[string]CacheEntry{}
	// Mutex to protect write and enumeration
	workCacheLock = &sync.Mutex{}
	// Time of last addition to cache
	cacheUpdateTime int64 = 0
)

func CacheUpdateTime() int64 { return cacheUpdateTime }

// Add a work result to the cache.  Account is optional (may be empty).
func addToCache(e rpcclient.WorkResponse, account string, timeComputed int64) {
	addToCacheInternal(CacheEntry{
		e.Hash,
		e.Work,
		e.Difficulty,
		e.Multiplier,
		account,
		"valid",
		timeComputed,
		0,
	})
}

// Mark in the cache that work request has started
func addToCacheStart(hash string) {
	addToCacheInternal(CacheEntry{
		hash,
		"",
		0,
		0,
		"",
		"computing",
		0,
		0,
	})
}

func addToCacheInternal(e CacheEntry) {
	if len(e.hash) == 0 {
		// empty key, omit
		return
	}
	workCacheLock.Lock()
	now := time.Now().Unix()
	e.timeAdded = now
	workCache[e.hash] = e
	cacheUpdateTime = now
	workCacheLock.Unlock()
}

func getFromCache(hash string) (CacheEntry, bool) {
	workCacheLock.Lock()
	e, ok := workCache[hash]
	workCacheLock.Unlock()
	if !ok {
		// not in cache
		return e, false
	}
	// found in cache
	return e, true
}

func cacheIsValid(e CacheEntry) bool {
	if e.status == "valid" {
		return true
	}
	return false
}

// Note: difficulty may be missing (0)
func cacheDiffIsOK(e CacheEntry, diff uint64) bool {
	if diff != 0 && e.difficulty != 0 && e.difficulty < diff {
		// but diff is smaller
		return false
	}
	// diff is OK (larger or equal)
	return true
}

// StatusCacheSize Return the current number of entries in the cache
func StatusCacheSize() int {
	return len(workCache)
}

func padString(val string) string {
	if len(val) == 0 {
		return "_"
	}
	return val
}

// Convert an entry to a single-line string representation
func entryToString(entry CacheEntry) string {
	if len(entry.hash) == 0 {
		return ""
	}
	return fmt.Sprintf("%v %v %x %v %v %v %v %v", padString(entry.hash), padString(entry.work), entry.difficulty, entry.multiplier,
		padString(entry.account), padString(entry.status), entry.timeComputed, entry.timeAdded)
}

// Fill cache entry from a single-line string represenation (parse it), see entryToString.
// Returns true on success.
func entryLoadFromString(line string, entry *CacheEntry) bool {
	tokens := strings.Split(line, " ")
	if len(tokens) < 2 {
		// mimium hash and work values are needed; this is too short
		return false
	}
	entry.hash = tokens[0]
	entry.work = tokens[1]
	if len(tokens) >= 8 {
		diff, _ := strconv.ParseUint(tokens[2], 16, 64)
		entry.difficulty = diff
		multip, _ := strconv.ParseFloat(tokens[3], 64)
		entry.multiplier = multip
		entry.account = tokens[4]
		entry.status = tokens[5]
		timeComputed, _ := strconv.ParseInt(tokens[6], 10, 64)
		timeAdded, _ := strconv.ParseInt(tokens[7], 10, 64)
		entry.timeComputed = timeComputed
		entry.timeAdded = timeAdded
	}
	return true
}

func RemoveOldEntries(cutoffAgeDays float64) {
	workCacheLock.Lock()
	oldSize := len(workCache)
	var newCache map[string]CacheEntry = make(map[string]CacheEntry, oldSize)
	now := time.Now().Unix()
	for key, entry := range workCache {
		ageDay := float64(now-entry.timeComputed) / float64(3600*24)
		if ageDay <= cutoffAgeDays {
			newCache[key] = entry
		}
	}
	newSize := len(newCache)
	if newSize != oldSize {
		workCache = newCache
		cacheUpdateTime = now
		log.Println("Cache: Removed old entries, size reduced from", oldSize, "to", newSize, "(cutoff", cutoffAgeDays, "days )")
	}
	workCacheLock.Unlock()
}
