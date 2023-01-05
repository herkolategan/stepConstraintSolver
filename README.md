# stepConstraintSolver

Generates a path of all possible permutations of a set of steps based on some arbitrary constraints.

## Example Use Case

1. Add a step
2. Add contraints to the step
3. Add another step with relational constraints to the previous step
4. Generate all possibilities

Examples of constraints:
- before: this step should happen before a specified step
- after: this step should happen after a specified step
- at least: the step should happen (be present) at least once
- at most: the step should not happen more than a specified amount.

## Code Example

```go
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
```

Outputs:
```
start      upgradeOne upgradeOne upgradeTwo mix1       mix2       finalize
start      upgradeOne upgradeOne upgradeTwo mix1       mix2
start      upgradeOne upgradeOne mix1       upgradeTwo mix2       finalize
start      upgradeOne upgradeOne mix1       upgradeTwo mix2
start      upgradeOne upgradeOne mix1       mix2       upgradeTwo finalize
start      upgradeOne upgradeOne mix1       mix2       upgradeTwo
start      upgradeOne upgradeTwo mix1       mix2       finalize
start      upgradeOne upgradeTwo mix1       mix2
start      upgradeOne mix1       upgradeOne upgradeTwo mix2       finalize
start      upgradeOne mix1       upgradeOne upgradeTwo mix2
start      upgradeOne mix1       upgradeOne mix2       upgradeTwo finalize
start      upgradeOne mix1       upgradeOne mix2       upgradeTwo
start      upgradeOne mix1       upgradeTwo mix2       finalize
start      upgradeOne mix1       upgradeTwo mix2
start      upgradeOne mix1       mix2       upgradeOne upgradeTwo finalize
start      upgradeOne mix1       mix2       upgradeOne upgradeTwo
start      upgradeOne mix1       mix2       upgradeTwo finalize
start      upgradeOne mix1       mix2       upgradeTwo
```
