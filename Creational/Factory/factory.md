Used for wholesale object creation (non-piecewise), unlike builder pattern, in a single invocation.
Used at places where object creation logic becomes too convoluted

In Golang, Factory functions also called Constructors - common approach in Go
When there are several factory functions, you group them and make a factory.

This factory can be passed as a param to some other method.

---------------------

There are two approaches:

1. Interface Factory -  step 1: create a new employee step 2: modify that employee
2. Factory Generator - we try to acheive both in one statement
    a) Functional approach
        - we have factories (func) stored in variables, and we can pass these variables into functions - functional programming
    b) Structural approach
        - 
