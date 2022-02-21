package com.company;

import java.util.ArrayList;
import java.util.List;

// Channel should know about subscribers
public class Channel {
    private List<Subscriber> subscribers = new ArrayList<>();
    private String title;

    public void subcribe(Subscriber s) {
        subscribers.add(s);
    }

    public void unSubcribe(Subscriber s) {
        subscribers.remove(s);
    }

    // Inform everyone
    public void notifySubscribers() {
        for(Subscriber s: subscribers) {
            s.Update();
        }
    }

    // Upload video
    public void upload(String title) {
        this.title = title;
        notifySubscribers();
    }
}
