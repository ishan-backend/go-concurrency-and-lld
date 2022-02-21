package com.company;

public class Main {

    public static void main(String[] args) {
	    // Singleton Pattern: these would return the same object which is static.
        Abc obj1 = Abc.getInstance();
        Abc obj2 = Abc.getInstance();

        // Observer pattern:
        Channel gfg = new Channel();
        Subscriber s1 = new Subscriber("Ishan Pandey");
        Subscriber s2 = new Subscriber("Ritul Pandey");
        // Subscribe to channel
        gfg.subcribe(s1);
        gfg.subcribe(s2);

        // subscriber knows that its subscribed
        s1.subscribeChannel(gfg);
        s2.subscribeChannel(gfg);

        // upload video
        gfg.upload("JAVA LLD");

        // unsubscribe
        gfg.unSubcribe(s2);
    }
}

/*
* Singleton pattern (concept):
* 1. Don't use public constructor to create instance of the class i.e. Abc obj1 = new Abc();
*     Infact, create a private constructor {this blocks from calling public constructor in main}.
*
* 2. Create an object of class inside the class itself, so that it returns the same instance always.
*
* 3. Define a static method that returns the object of class and name it as getInstance()
* */
class Abc {
    static Abc obj1 = new Abc();
    private Abc () {

    }
    public static Abc getInstance() {
        return obj1;
    }
}

