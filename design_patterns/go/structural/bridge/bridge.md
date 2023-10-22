Example of Cartesian Product Complexity explosion

e.g
1. Make common ThreadScheduler.
2. It can be pre-emptive or co-operative.
3. It can be for Windows or UNIX based systems

- 2x2 = 4 possibilities here, for 3x3 = 9 different structures which is un-manageable.
- Bridge tries to avoid this complexity explosion scenario.

```go
type ThreadScheduler struct {
	*IPlatformScheduler
}
// IPlatformScheduler is implemented by UnixScheduler and WindowsScheduler
// Co-operativeThreadScheduler and Pre-emptiveThreadScheduler implement ThreadScheduler.

Seperation of Heirarchy - instead of one deep tree, we try to have shallow trees or flat lists of implementation
```

https://refactoring.guru/design-patterns/bridge/go/example
* Bridge tries to decouple some abstraction from implementation
* We try to create parallel trees/Hierarchies - thus prevent compl explosion
* One hierarchy that makes of another rather than vector product kind of approach
* Bridge pattern = stronger form of encapsulation