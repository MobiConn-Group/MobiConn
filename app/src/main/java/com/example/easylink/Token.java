package com.example.easylink;

import android.app.Application;

public class Token extends Application {
    private String token;

    public String getToken() {
        return token;
    }

    public void setToken(String _token) {
        token = _token;
    }
}