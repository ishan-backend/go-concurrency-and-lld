package com.company;

/*
 * Observer Pattern: Need ?
 * Youtube: channel -> many subscribers, Notifications is needed.
 * Big youtube object -> inside many channels, inside channel -> many subscriber. On uploading video to channel, notif's sent to subscribers. Push based model.
 * So subscribers are observers to a subject {channel}.
 *
 *
 *
 * */
public class Subscriber {
    private String name;
    private Channel channel;

    public Subscriber(String name) {
        this.name = name;
    }

    public void Update() {
        System.out.println("Video uploaded"); // notification to subcriber to a channel
    }

    public void subscribeChannel(Channel ch) { // subs should know which channel its subcribing to.
        channel = ch;
    }
}
