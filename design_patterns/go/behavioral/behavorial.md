https://refactoring.guru/design-patterns/behavioral-patterns
* Behavioral design patterns are concerned with algorithms and the assignment of responsibilities between objects.

**Top 3 Behavioural Design Patterns**:
1. **Observer**
    * Sometimes we need to be informed when certain things happen.
      * When objects field changes
      * When object does something
      * Some external event happens,and we want our system to handle it
    * We want to generate events from object when something happens with it, and we want to listen to these events/notified when these occur.
    * Two participants:
      * Observable/Subject: entity generating events
      * Observers: wishes to be informed about events
      * **general_construct_to_support_notification.go** - as soon as a person gets ill, a doctor is called for
    * Property observers: **property_observers.go**
      * If subject's property changes, notify the observers.
      * You would need getter() and setter() for that property of subject/observable
      * Observer Mgmt encapsulates subject
      * There are certain challenges with using Change Notification on getter and setter: **property_dependencies.go**
        * A property can be just read only GetIfOkayForInjection, made from some other attributes from subject -> i.e. age; age can only be changed in setters
        * If you have lots of properties dependent on age property, or five properties (func) dependent on one property, then its nightmare
          * We dont define dependency properties inside definition of setters.
          * Instead we build high level framework - map where all the dependencies of all the properties are catalogued.
          * And then you iterate through this map, and perform notifications in a regularised way.
        
2. **Chain Of Responsibility**
    * Mirrored a lot in real-world
    * examples
      * Unethical behaviour by an employee, who takes blame?
        1. Employee
        2. Manager
        3. CEO
      * You click on graphical element on a form, who gets to handle the event?
        1. Button handles it, it stops further processing
        2. Group box - button does nothing, group box underlying handles event
        3. Underlying window - button, underlying group box do nothing, window handles it
      * CCG computer game
        1. Creature has attack and defense values.
        2. Can get booster cards - how to apply these values to creature?
    * Chain of Responsibility: chain of components, who all get a chance to process a command or a query; optionally having default processing implementation and ability to terminate processing chain (every element of chain can have).
    * Linked List of method invocations
    * Command Query Separation:
      * Command: Mutable - action or query which changes value
      * Query: Immutable - asking for information

3. **Strategy**
   * Many algorithms can be broken down into high level and low level parts
   * High level algorithms can be reused for different purposes. Low level parts are specific to different use-case.
   * Tea:
     * Pour water, boil water, add milk, add sugar - High level - skeleton
     * Pour teabag - for tea - beverage specific strategy
     * Pour coffee - for coffee - beverage specific strategy
   * So, Strategy = divide an algorithm to its 'skeleton' and its concrete implementation steps which can be varied/switched at runtime.
   
4. **Mediator**: [Todo]


