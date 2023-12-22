package com.example.easylink;

import android.app.Application;

import java.io.File;
import java.util.List;

public class Resource extends Application {
    private String token;
    private List<File> images;
    private String ip;
    private String username;
    private List<String> uploadFiles;

    public String getToken() {
        return token;
    }

    public void setToken(String token) {
        this.token = token;
    }

    public List<File> getImages() {
        return images;
    }

    public void setImages(List<File> images) {
        this.images = images;
    }

    public String getIp() {
        return ip;
    }

    public void setIp(String ip) {
        this.ip = ip;
    }

    public String getUsername() {
        return username;
    }

    public void setUsername(String username) {
        this.username = username;
    }

    public List<String> getUploadFiles() {
        return uploadFiles;
    }

    public void setUploadFiles(List<String> uploadFiles) {
        this.uploadFiles = uploadFiles;
    }

}