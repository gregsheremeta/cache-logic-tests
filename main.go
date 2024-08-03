package main

import (
	"fmt"
	"reflect"
	"runtime"
)

const dashes string = "------------------------------------------------------------------------------"

type CacheUsage int

const (
	NoCache CacheUsage = iota
	UseCache
	NotSet
)

type algorithmToTest func(bool, CacheUsage, CacheUsage, CacheUsage) bool

// Helper function to get the name of a function
func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func main() {
	var algorithm algorithmToTest

	algorithm = shouldTryToUseCache
	runAllTestsUsingAlgorithm(algorithm)

	algorithm = shouldTryToUseCache_UsingOnlyBooleans
	runAllTestsUsingAlgorithm(algorithm)
}

func runAllTestsUsingAlgorithm(algorithm algorithmToTest) {

	fmt.Printf("\n%s\n%s\nrunning all tests using algorithm: %s\n%s\n%s\n", dashes, dashes, getFunctionName(algorithm), dashes, dashes)

	testDescription := "server allows caching, but it was not requested at any level"
	serverLevelForceDisabled := false
	runLevel := NotSet
	componentLevel := NotSet
	pipelineLevel := NotSet
	expected := false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server disables all caching, random settings at other levels"
	serverLevelForceDisabled = true
	runLevel = NotSet
	componentLevel = NotSet
	pipelineLevel = NotSet
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server disables all caching, random settings at other levels 2"
	serverLevelForceDisabled = true
	runLevel = UseCache
	componentLevel = NotSet
	pipelineLevel = UseCache
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server disables all caching, random settings at other levels 3"
	serverLevelForceDisabled = true
	runLevel = NotSet
	componentLevel = NoCache
	pipelineLevel = UseCache
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests check for when only one level is set to UseCache

	testDescription = "server allows caching, requested only at run level"
	serverLevelForceDisabled = false
	runLevel = UseCache
	componentLevel = NotSet
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, requested only at component level"
	serverLevelForceDisabled = false
	runLevel = NotSet
	componentLevel = UseCache
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, requested only at pipeline level"
	serverLevelForceDisabled = false
	runLevel = NotSet
	componentLevel = NotSet
	pipelineLevel = UseCache
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests check for when a lower precedence level is set to NoCache, but a higher precedence level is set to UseCache.
	// the higher precedence level should win

	testDescription = "server allows caching, disabled at pipeline level, but requested at run level (run level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = UseCache
	componentLevel = NotSet
	pipelineLevel = NoCache
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, disabled at component level, but requested at run level (run level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = UseCache
	componentLevel = NoCache
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, disabled at pipeline level, but requested at component level (component level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = NotSet
	componentLevel = UseCache
	pipelineLevel = NoCache
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests is the same as above, but the requests are just flipped.
	// so, check for when a lower precedence level is set to UseCache, but a higher precedence level is set to NoCache.
	// the higher precedence level should win

	testDescription = "server allows caching, requested at pipeline level, but disabled at run level (run level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = NoCache
	componentLevel = NotSet
	pipelineLevel = UseCache
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, requested at component level, but disabled at run level (run level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = NoCache
	componentLevel = UseCache
	pipelineLevel = NotSet
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server allows caching, requested at pipeline level, but disabled at component level (component level should take precedence)"
	serverLevelForceDisabled = false
	runLevel = NotSet
	componentLevel = NoCache
	pipelineLevel = UseCache
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel, expected)

	// Notice that the previous block all work in `shouldTryToUseCache` but not in `shouldTryToUseCache_UsingOnlyBooleans` !!
	// This is because the latter does not have the concept of NotSet.
	//
	// What happens is that any time a higher precendence level says "I definitely don't want cache",
	// the lower precedence levels are winning when they are true!
	// This happens because basically shouldTryToUseCache_UsingOnlyBooleans is a giant boolean OR !!
	// It will say "yes I want cache" if ANY of the levels say they want cache, which is not what we want.
	// The precedence is not respected when using only booleans.

}

func runTestAndPrintResults(algorithm algorithmToTest, testDescription string, serverLevelForceDisabled bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage, expected bool) {
	fmt.Printf("\n\nrunning test case [%s]...\n%s\n", testDescription, dashes)
	result := algorithm(serverLevelForceDisabled, runLevel, componentLevel, pipelineLevel)
	fmt.Printf("shouldTryToUseCache: %v \t expected: %v\n", result, expected)
	if result != expected {
		fmt.Println("!!!!!!!!!!!!!!!! TEST FAILED !!!!!!!!!!!!!!!!")
	} else {
		fmt.Println("TEST PASSED")
	}
}

///////////////////////////////////////////
// below are the various algorithms to test
///////////////////////////////////////////

// this algorithm uses switches with 3 values - UseCache, NoCache, NotSet
func shouldTryToUseCache(serverLevelForceDisabled bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage) bool {

	if serverLevelForceDisabled == true {
		fmt.Println("will NOT try to use cache because it was globally disabled at the server level")
		return false
	}

	// if we're here, serverLevelForceDisabled was false
	fmt.Println("the server-wide force caching disabled is not set. Will proceed to check the other levels")

	if runLevel == UseCache {
		fmt.Println("will try to use cache because it was requested at Run time")
		return true
	}

	if runLevel == NoCache {
		fmt.Println("will NOT try to use cache because it was requested not to at Run time")
		return false
	}

	// if we're here, run level was NotSet
	fmt.Println("no caching setting was found at the run level. Will now check the component level")

	if componentLevel == UseCache {
		fmt.Println("will try to use cache because it was requested at the component level")
		return true
	}

	if componentLevel == NoCache {
		fmt.Println("will NOT try to use cache because it was requested not to at the component {{component_name}} level")
		return false
	}

	// if we're here, run level was NotSet
	fmt.Println("no caching setting was found at the component level. Will now check the pipeline level")

	if pipelineLevel == UseCache {
		fmt.Println("will try to use cache because it was requested at the pipeline level")
		return true
	}

	if pipelineLevel == NoCache {
		fmt.Println("will NOT try to use cache because it was requested not to at the pipeline level")
		return false
	}

	// if we're here, pipeline level was NotSet
	fmt.Println("no caching setting was found at the pipeline level")

	// overall default case is here. Default is no cache!
	fmt.Println("skipped caching. It wasn't disabled globally, but it wasn't requested at any level")
	return false
}

// this algorithm uses booleans - no concept of NotSet
func shouldTryToUseCache_UsingOnlyBooleans(serverLevelForceDisabled bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage) bool {

	if serverLevelForceDisabled == true {
		fmt.Println("will NOT try to use cache because it was globally disabled at the server level")
		return false
	}

	// make the booleans
	runLevelBoolean := runLevel == UseCache
	componentLevelBoolean := componentLevel == UseCache
	pipelineLevelBoolean := pipelineLevel == UseCache

	// if we're here, serverLevelForceDisabled was false
	fmt.Println("the server-wide force caching disabled is not set. Will proceed to check the other levels")

	if runLevelBoolean {
		fmt.Println("will try to use cache because it was requested at Run time")
		return true
	}

	if componentLevelBoolean {
		fmt.Println("will try to use cache because it was requested at the component level")
		return true
	}

	if pipelineLevelBoolean {
		fmt.Println("will try to use cache because it was requested at the pipeline level")
		return true
	}

	// notice that the above 3 ifs could have been written as a single boolean OR:
	// if runLevelBoolean || componentLevelBoolean || pipelineLevelBoolean  ... return true
	// not what we want!

	// overall default case is here. Default is no cache!
	fmt.Println("skipped caching. It wasn't disabled globally, but it wasn't requested at any level")
	return false
}
