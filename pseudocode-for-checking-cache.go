
if cache_is_enabled {
	// try to hit the cache
	err, got_a_cache_hit, cached_results := check_cache()
	if err != nil {
		log(“tried to use cache, but hit an error! Will treat as a cache miss and proceed with podSpecPatch.”)
		// fall through to pod spec patch code below
	}
	else if got_a_cache_hit {
		log(“got a cache hit! Will not do podSpecPatch. Will reuse artifacts from previous result.”)

		// do something with the cache result
		mark_this_task_as_using_cache()
		save_artifacts_for_this_task(cached_results)

		return // we are done here! Do NOT fall through to pod spec patch
	}
	else {
		log(“got a cache miss. Will proceed with podSpecPatch.”)
		// fall through to pod spec patch code below
	}
}

// if you got to here, you either didn’t even try to use cache, or you missed, or you got an error while trying (which we treat as a miss)
// podSpecPatch time! Time to fire up a pod

prepare pod spec patch ...
