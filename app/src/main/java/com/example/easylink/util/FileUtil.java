package com.example.easylink.util;

import android.content.Intent;
import android.util.Log;
import android.view.View;

import java.io.File;
import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class FileUtil {
    // 使用OkHttp上传文件
    public static void uploadFile(File file, String url) {
        // 创建OkHttpClient实例
        OkHttpClient client = new OkHttpClient();
        // file是要上传的文件
        RequestBody fileBody = RequestBody.create(MediaType.parse("multipart/form-data"), file);
        // 不仅可以支持传文件，还可以在传文件的同时，传参数
        MultipartBody requestBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("Files", file.getName(), fileBody)
                .build();
        Request request = new Request.Builder()
                .url(url)
                .post(requestBody)
                .build();
        Call call = client.newCall(request);
        call.enqueue(new Callback() {
            @Override
            public void onResponse(Call call, Response response) throws IOException {
                // 文件上传成功
                if (response.isSuccessful()) {
                    Log.i("uploadFile", "onResponse: " + response.body().string());
                } else {
                    Log.i("uploadFile", "onResponse: " + response.message());
                }
            }

            @Override
            public void onFailure(Call call, IOException e) {
                // 文件上传失败
                Log.i("uploadFile", "onFailure: " + e.getMessage());
            }
        });
    }
}
