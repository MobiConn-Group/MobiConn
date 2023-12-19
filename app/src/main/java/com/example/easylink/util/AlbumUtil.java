package com.example.easylink.util;

import android.content.Context;
import android.widget.Toast;

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

public class AlbumUtil {
    private void album_file(Context context, File file) throws IOException {
        // 创建OkHttpClient实例
        OkHttpClient client = new OkHttpClient();
        // file是要上传的文件
        RequestBody fileBody = RequestBody.create(MediaType.parse("multipart/form-data"), file);
        // 不仅可以支持传文件，还可以在传文件的同时，传参数
        MultipartBody requestBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("参数名1", "参数1")
                .addFormDataPart("接收文件的参数名", "文件名", fileBody)
                .build();
        Request request = new Request.Builder().url("请求的url").post(requestBody).build();
        Response response = client.newCall(request).execute();
    }
}
