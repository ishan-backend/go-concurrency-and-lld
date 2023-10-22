**1. Factory (Simple & Abstract)**
* https://refactoring.guru/design-patterns/factory-method
* https://refactoring.guru/design-patterns/factory-method/go/example
* https://refactoring.guru/design-patterns/factory-comparison
  * *Simple Factory* / *Factory Method*:
    * **example.go/func createGun** or **interface_factory.go/**
    * The Factory Method is a creational design pattern that provides an interface for creating objects but allows subclasses to alter the type of an object that will be created.

  * *Abstract Factory* (https://refactoring.guru/design-patterns/abstract-factory):
    * **factory_generator.go** / **example_abstract_factory.go**
    * https://refactoring.guru/design-patterns/abstract-factory/go/example
    * The Abstract Factory is a creational design pattern that allows producing families of related or dependent objects without specifying their concrete classes.
    * Many designs start by using Factory Method (less complicated and more customizable via subclasses) and evolve toward Abstract Factory, Prototype, or Builder (more flexible, but more complicated).

**Notes**:
 - Builder focuses on constructing complex objects step by step. Abstract Factory specializes in creating families of related objects. Abstract Factory returns the product immediately, whereas Builder lets you run some additional construction steps before fetching the product.

 - Abstract Factory classes are often based on a set of Factory Methods, but you can also use Prototype to compose the methods on these classes.

 - Abstract Factory can serve as an alternative to Facade when you only want to hide the way the subsystem objects are created from the client code.

 - You can use Abstract Factory along with Bridge. This pairing is useful when some abstractions defined by Bridge can only work with specific implementations. In this case, Abstract Factory can encapsulate these relations and hide the complexity from the client code.

 - Abstract Factories, Builders and Prototypes can all be implemented as Singletons.

**2. Builder**:


**5. Object Pool**:
* https://nuculabs.dev/p/object-pool-pattern
* https://pkg.go.dev/sync#Pool
* used in situations when constructing objects is a costly operation, for example building an HTTPClient or DatabaseClient object can take some time.
* By having a pool of resources, the resources are requested from the pool when needed and then returned when not needed so they can be reused later.
* you'll just grab an instance instead of constructing it again from scratch.
* In Go we have sync.pool which implements the Object Pool pattern for us, we just need to give it a New function that returns a pointer.
* To track how many active instances, we have of the object Resource, we use the counter variable.
