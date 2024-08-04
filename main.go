package main

import (
	"fmt"
	"reflect"
	"runtime"
)

const dashes string = "------------------------------------------------------------------------------"

type CacheUsage int

const (
	CacheDisabled CacheUsage = iota
	CacheEnabled
	NotSet
)

type algorithmToTest func(bool, CacheUsage, CacheUsage, CacheUsage) bool

// Helper function to get the name of a function
func getFunctionName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func main() {
	var algorithm algorithmToTest

	algorithm = isCacheEnabled
	runAllTestsUsingAlgorithm(algorithm)

	algorithm = isCacheEnabled_UsingOnlyBooleans
	runAllTestsUsingAlgorithm(algorithm)
}

func runAllTestsUsingAlgorithm(algorithm algorithmToTest) {

	fmt.Printf("\n%s\n%s\nrunning all tests using algorithm: %s\n%s\n%s\n", dashes, dashes, getFunctionName(algorithm), dashes, dashes)

	testDescription := "server supports caching, but it was not enabled at any level"
	serverLevelCachingSupported := true
	runLevel := NotSet
	componentLevel := NotSet
	pipelineLevel := NotSet
	expected := false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server set all caching unsupported, random settings at other levels"
	serverLevelCachingSupported = false
	runLevel = NotSet
	componentLevel = NotSet
	pipelineLevel = NotSet
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server set all caching unsupported, random settings at other levels 2"
	serverLevelCachingSupported = false
	runLevel = CacheEnabled
	componentLevel = NotSet
	pipelineLevel = CacheEnabled
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server set all caching unsupported, random settings at other levels 3"
	serverLevelCachingSupported = false
	runLevel = NotSet
	componentLevel = CacheDisabled
	pipelineLevel = CacheEnabled
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests check for when only one level is set to CacheEnabled

	testDescription = "server supports caching, enabled only at run level"
	serverLevelCachingSupported = true
	runLevel = CacheEnabled
	componentLevel = NotSet
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, enabled only at component level"
	serverLevelCachingSupported = true
	runLevel = NotSet
	componentLevel = CacheEnabled
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, enabled only at pipeline level"
	serverLevelCachingSupported = true
	runLevel = NotSet
	componentLevel = NotSet
	pipelineLevel = CacheEnabled
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests check for when a lower precedence level is set to CacheDisabled, but a higher precedence level is set to CacheEnabled.
	// the higher precedence level should win

	testDescription = "server supports caching, disabled at pipeline level, but enabled at run level (run level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = CacheEnabled
	componentLevel = NotSet
	pipelineLevel = CacheDisabled
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, disabled at component level, but enabled at run level (run level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = CacheEnabled
	componentLevel = CacheDisabled
	pipelineLevel = NotSet
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, disabled at pipeline level, but enabled at component level (component level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = NotSet
	componentLevel = CacheEnabled
	pipelineLevel = CacheDisabled
	expected = true
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	// the following block of tests is the same as above, but the requests are just flipped.
	// so, check for when a lower precedence level is set to CacheEnabled, but a higher precedence level is set to CacheDisabled.
	// the higher precedence level should win

	testDescription = "server supports caching, enabled at pipeline level, but disabled at run level (run level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = CacheDisabled
	componentLevel = NotSet
	pipelineLevel = CacheEnabled
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, enabled at component level, but disabled at run level (run level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = CacheDisabled
	componentLevel = CacheEnabled
	pipelineLevel = NotSet
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	testDescription = "server supports caching, enabled at pipeline level, but disabled at component level (component level should take precedence)"
	serverLevelCachingSupported = true
	runLevel = NotSet
	componentLevel = CacheDisabled
	pipelineLevel = CacheEnabled
	expected = false
	runTestAndPrintResults(algorithm, testDescription, serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel, expected)

	// Notice that the previous block all work in `isCacheEnabled` but not in `isCacheEnabled_UsingOnlyBooleans` !!
	// This is because the latter does not have the concept of NotSet.
	//
	// What happens is that any time a higher precendence level says "I definitely don't want cache",
	// the lower precedence levels are winning when they are true!
	// This happens because basically isCacheEnabled_UsingOnlyBooleans is a giant boolean OR !!
	// It will say "yes I want cache" if ANY of the levels say they want cache, which is not what we want.
	// The precedence is not respected when using only booleans.

}

func runTestAndPrintResults(algorithm algorithmToTest, testDescription string, serverLevelCachingSupported bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage, expected bool) {
	fmt.Printf("\n\nrunning test case [%s]...\n%s\n", testDescription, dashes)
	result := algorithm(serverLevelCachingSupported, runLevel, componentLevel, pipelineLevel)
	fmt.Printf("isCacheEnabled: %v \t expected: %v\n", result, expected)
	if result != expected {
		fmt.Println("!!!!!!!!!!!!!!!! TEST FAILED !!!!!!!!!!!!!!!!")
	} else {
		fmt.Println("TEST PASSED")
	}
}

///////////////////////////////////////////
// below are the various algorithms to test
///////////////////////////////////////////

// this algorithm uses switches with 3 values - CacheEnabled, CacheDisabled, NotSet
func isCacheEnabled(serverLevelCachingSupported bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage) bool {

	if !serverLevelCachingSupported {
		fmt.Println("caching is disabled for this component (will NOT try to use cache) because it was set to unsupported (globally disabled) at the server level")
		return false
	}

	// if we're here, serverLevelCachingSupported was true
	fmt.Println("the server-wide caching supported setting is true. Caching is not globally unsupported. Will proceed to check the other levels to see if caching is enabled in any of those")

	if runLevel == CacheEnabled {
		fmt.Println("will try to use cache because it was enabled at Run time")
		return true
	}

	if runLevel == CacheDisabled {
		fmt.Println("will NOT try to use cache because it was disabled at Run time")
		return false
	}

	// if we're here, run level was NotSet
	fmt.Println("no caching setting was found at the run level. Will now check the component level")

	if componentLevel == CacheEnabled {
		fmt.Println("will try to use cache because it was enabled at the component level")
		return true
	}

	if componentLevel == CacheDisabled {
		fmt.Println("will NOT try to use cache because it was disabled at the component level")
		return false
	}

	// if we're here, run level was NotSet
	fmt.Println("no caching setting was found at the component level. Will now check the pipeline level")

	if pipelineLevel == CacheEnabled {
		fmt.Println("will try to use cache because it was enabled at the pipeline level")
		return true
	}

	if pipelineLevel == CacheDisabled {
		fmt.Println("will NOT try to use cache because it was disabled at the pipeline level")
		return false
	}

	// if we're here, pipeline level was NotSet
	fmt.Println("no caching setting was found at the pipeline level")

	// overall default case is here. Default is no cache!
	fmt.Println("skipped caching. It wasn't globally unsupported (force-disabled) by an admin, but it wasn't enabled at any level")
	return false
}

// this algorithm uses booleans - no concept of NotSet
func isCacheEnabled_UsingOnlyBooleans(serverLevelCachingSupported bool, runLevel CacheUsage, componentLevel CacheUsage, pipelineLevel CacheUsage) bool {

	if !serverLevelCachingSupported {
		fmt.Println("caching is disabled for this component (will NOT try to use cache) because it was set to unsupported (globally disabled) at the server level")
		return false
	}

	// if we're here, serverLevelCachingSupported was true
	fmt.Println("the server-wide caching supported setting is true. Caching is not globally unsupported. Will proceed to check the other levels to see if caching is enabled in any of those")

	// make the booleans
	runLevelBoolean := runLevel == CacheEnabled
	componentLevelBoolean := componentLevel == CacheEnabled
	pipelineLevelBoolean := pipelineLevel == CacheEnabled

	if runLevelBoolean {
		fmt.Println("will try to use cache because it was enabled at Run time")
		return true
	}

	if componentLevelBoolean {
		fmt.Println("will try to use cache because it was enabled at the component level")
		return true
	}

	if pipelineLevelBoolean {
		fmt.Println("will try to use cache because it was enabled at the pipeline level")
		return true
	}

	// notice that the above 3 ifs could have been written as a single boolean OR:
	// if runLevelBoolean || componentLevelBoolean || pipelineLevelBoolean  ... return true
	// not what we want!

	// overall default case is here. Default is no cache!
	fmt.Println("skipped caching. It wasn't globally unsupported (force-disabled) by an admin, but it wasn't enabled at any level")
	return false
}
