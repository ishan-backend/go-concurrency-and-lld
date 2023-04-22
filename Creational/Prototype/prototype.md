An existing design (partially or fully constructed) is called prototype.
We make a copy/clone of prototype and customize/use it
    - deep copy
        - while copying pointers you have to be careful, you dont need to give pointer -> pointer -> object, but rather the same pointer

We can go ahead a step further, and create a PROTOTYPE FACTORY:
    - a kind of factory which serves prototypes and maybe even lets you customize the prototype at the point of creation.


-------------------------------

Serializing Constructs - are meant to be smart
If you give serialiser a struct, its going to figure out if string - to serialise, else a pointer, then the value of object will be serialised.
i.e. deep copy performed for you automatically.

-> we will have different type of DeepCopy() method - using encoders and decoders to save all data and read from it
