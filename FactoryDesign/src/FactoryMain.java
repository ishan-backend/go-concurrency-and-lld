
public class FactoryMain {

    public static void main(String p[]) {

        // Client code: since Windows is exposed.
        OS obj1 = new Android();
        obj1.spec();

        OS obj2 = new IOS();
        obj2.spec();

        OS obj3 = new Windows();
        obj3.spec();

        /*
        * If you want to add new OS and if you are using Windows class directly here, it means client knows this, And might need to recompile whole package.
        *
        * Solution: Factory Pattern.
        * 1. Create a new class as OSFactory -> gives object of type OS with input as type of OS you need.
        * 2.
        * */

        // You can change internal OSFactory class but not the client application.
        OSFactory osf = new OSFactory(); // create instance of new wrapper class
        OS obj4 = osf.getInstance("Open");
        obj4.spec();
    }
}
