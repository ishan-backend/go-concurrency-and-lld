https://refactoring.guru/design-patterns/structural-patterns

**1. Adapter**:
* Adapter is a structural design pattern that allows objects with incompatible interfaces to collaborate.


**2. Bridge**:
* Bridge is a structural design pattern that lets you split a large class or a set of closely related classes into two separate hierarchies—abstraction and implementation—which can be developed independently of each other.
* Bridge prevents = "Cartesian Product Complexity Explosion"

**3. Composite**:
* Objects can use other objects via composition
* Some scalar and composed objects need similar behaviour
* Composite pattern allows us to treat both types of objects uniformaly
* Iteration is supported via iterator design pattern

**4. Decorator**:
* Adding behaviour without altering the type itself.
* Augment existing object with additional functionality.
* Want to stick by OCP
* Want to keep new enhancements separate (SRP)
* Need to be able to interact with original structs/data and methods using embedding
* https://refactoring.guru/design-patterns/decorator
* How to combine functionalities of several structures within a single structure?
  * e.g. Dragon = Bird + Lizard
  * In Go, how can you support it?
    * Aggregation of structs having same fields is a problem - Ambiguous reference
    * Even with getters and setters, few changed properties can lead to inconsistent behaviour
  * So we need to build proper decorator around Lizard and Bird types
  * Vishnu = Man + Lion ; multiple_aggregation_improvised.go
  * decorator:
    * There are certain things you lose also while you use decorator
    * on circle.Resize() on ordinary shape
    * On coloredShape over circle, we cannot get Resize() is not available, since we did not aggregate it
    * We cannot add Resize() to interface because other shapes might not support resize
  * Decorators can be composed i.e. decorators on top of decorators
  