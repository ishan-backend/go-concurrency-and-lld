S : Single responsibility principle <br>
O : Open/Closed principle <br>
L : Liskov substitution principle <br>
I : Interface segregation principle <br>
D : Dependency inversion principle <br>

#### Symptoms of rotting code design:
1. Rigidity -> Code to be difficult to change even if the change is small
2. Fragility -> Code to break whenever a new change is introduced in the system
3. Immobility -> Code being not reusable
4. Viscosity -> Hacking rather than finding a solution that preserve the design when it comes to change

#### Single responsibility principle:
* two separate aspects of a problem need to be handled by different modules.

#### Open/Closed principle:
* a module should be open for extensions, but closed for modification

#### Liskov substitution principle:
* Derived methods should expect no more and provide no less
* In an object-oriented class-based language, the concept of the Liskov substitution principle is that a user of a base class should be able to function properly of all derived classes.

#### Interface segregation principle:
* Many client specific interfaces are better than one general purpose interface
* each interface must be defined in a way so that it provides the exact and full set of functionalities needed by at least one of the client

#### Dependency inversion principle:


Read: https://s8sg.medium.com/solid-principle-in-go-e1a624290346