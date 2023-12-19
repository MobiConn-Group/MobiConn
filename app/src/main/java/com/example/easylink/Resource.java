package com.example.easylink;

import android.app.Application;

import java.io.File;

public class Resource extends Application {
    private String token;
    private File[] images;

    public String getToken() {
        return token;
    }

    public void setToken(String _token) {
        token = _token;
    }

    public File[] getImages() {
        return images;
    }

    public void setImages(File[] _images) {
        images = _images;
    }

}