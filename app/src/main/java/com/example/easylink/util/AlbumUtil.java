package com.example.easylink.util;

import android.content.Context;
import android.widget.Toast;

import com.example.easylink.Resource;

import java.io.File;
import java.io.IOException;
import java.util.List;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class AlbumUtil {
    public static void album(String url, Resource resource) {
        List<File> album = resource.getImages();
        for (File image : album) {
            FileUtil.uploadFile(image, url);
        }
    }
}
