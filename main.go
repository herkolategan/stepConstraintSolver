package main

import "fmt"

type (
	step struct {
		id            string        // unique id (name) of the step
		condition     stepPredicate // composable condition to be met before the step can be added to the path
		postCondition stepPredicate // composable condition to be met on complete path before the path is accepted
	}

	stepPredicate func(path) bool // stepPredicate checks if a condition holds true for a given path
	path          []*step

	// steps is a helper struct used during evaluation
	steps struct {
		path
		limit int
	}
)

const limit = 10000

func main() {
	var (
		s          steps
		start      = s.newStep("start     ").atMost(1).atLeast(1)
		upgradeOne = s.newStep("upgradeOne").atMost(2).atLeast(1)
		upgradeTwo = s.newStep("upgradeTwo").atMost(1).atLeast(1)
		finalize   = s.newStep("finalize  ").atMost(1).atLeast(0)
		mixIn1     = s.newStep("mix1      ").atMost(1).atLeast(1)
		mixIn2     = s.newStep("mix2      ").atMost(1).atLeast(1)
	)

	upgradeOne.
		after(start).
		before(upgradeTwo)

	upgradeTwo.
		after(upgradeOne).
		before(finalize)

	mixIn1.
		after(upgradeOne)

	mixIn2.
		after(mixIn1).
		before(finalize)

	finalize.
		after(upgradeTwo)

	// Evaluate all possible execution paths (a randomizer can then be used to choose one)
	paths := s.evaluate(limit)

	// Print results
	printPaths(&paths)
}

// Implementation Example

// evaluate all possible execution paths given a set of steps
// optionally a rng can then be used to select one arbitrarily
func evaluate(steps *steps, path path, paths *[]path) {
	// avoid a stack overflow if the conditions result in infinite possibilities
	if steps.limit <= 0 {
		return
	}
	steps.limit--

	// loop through steps to check if a new path can be evaluated based on conditions for each step
	for _, s := range steps.path {
		nextPath := append([]*step(nil), path...)
		nextPath = append(nextPath, s)
		if s.condition == nil || s.condition(nextPath) {
			evaluate(steps, nextPath, paths)
		}
	}

	// confirm post conditions are met before adding the path as a possible execution
	postConditionsMet := true
	for _, s := range steps.path {
		if s.postCondition != nil && !s.postCondition(path) {
			postConditionsMet = false
			break
		}
	}
	if postConditionsMet {
		*paths = append(*paths, path)
	}
}

// count the number of times a step is present in a path
func (p *path) count(s *step) int {
	count := 0
	for _, e := range *p {
		if e.id == s.id {
			count++
		}
	}
	return count
}

func (s *steps) evaluate(limit int) []path {
	var paths []path
	s.limit = limit
	evaluate(s, nil, &paths)
	return paths
}

// newStep creates a step and adds it to the set of steps that will be evaluated
func (s *steps) newStep(name string) *step {
	newStep := &step{name, nil, nil}
	s.path = append(s.path, newStep)
	return newStep
}

// addCondition to a target condition predicate (the target is either the in-path condition or the post path condition)
func (s *step) addCondition(target *stepPredicate, pred stepPredicate) *step {
	if *target == nil {
		*target = pred
		return s
	}
	originalCondition := *target
	*target = func(path path) bool {
		return originalCondition(path) && pred(path)
	}
	return s
}

// after adds a condition that a step can only execute after a given step
func (s *step) after(afterStep *step) *step {
	condition := func(p path) bool {
		if p.count(afterStep) >= 1 {
			return true
		}
		return false
	}
	return s.addCondition(&s.condition, condition)
}

// before adds a condition that a step can only execute before a given step
func (s *step) before(beforeStep *step) *step {
	condition := func(p path) bool {
		if p.count(beforeStep) == 0 {
			return true
		}
		return false
	}
	return s.addCondition(&s.condition, condition)
}

// atMost adds a condition that a step can occur at most for the given count
func (s *step) atMost(count int) *step {
	condition := func(p path) bool {
		if p.count(s) <= count {
			return true
		}
		return false
	}
	return s.addCondition(&s.condition, condition)
}

// atLeast adds a post condition that a step should occur at least for the given count
func (s *step) atLeast(count int) *step {
	condition := func(p path) bool {
		if p.count(s) >= count {
			return true
		}
		return false
	}
	return s.addCondition(&s.postCondition, condition)
}

func printPaths(paths *[]path) {
	for _, p := range *paths {
		for _, s := range p {
			fmt.Printf("%s ", s.id)
		}
		fmt.Println()
	}
}
