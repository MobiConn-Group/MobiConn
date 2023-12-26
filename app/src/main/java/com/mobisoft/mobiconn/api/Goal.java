package com.mobisoft.mobiconn.api;

public class Goal {
    private final String action;

    private final String information;

    public String getAction() {
        return action;
    }

    public String getInformation() {
        return information;
    }

    public Goal(String action, String information) {
        this.action = action;
        this.information = information;
    }
}
