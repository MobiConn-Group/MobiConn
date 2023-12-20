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
    public static void album_file(Context context, File file, String url) throws IOException {
        // 创建OkHttpClient实例
        OkHttpClient client = new OkHttpClient();
        // file是要上传的文件
        RequestBody fileBody = RequestBody.create(MediaType.parse("multipart/form-data"), file);
        // 不仅可以支持传文件，还可以在传文件的同时，传参数
        MultipartBody requestBody = new MultipartBody.Builder()
                .setType(MultipartBody.FORM)
                .addFormDataPart("Files", "image", fileBody)
                .build();
        Request request = new Request.Builder().url(url).post(requestBody).build();
        Response response = client.newCall(request).execute();
    }

    public static void album(Context context, String url, Resource resource) {
        List<File> album = resource.getImages();
        for (File image : album) {
            try {
                album_file(context, image, url);
            } catch (IOException e) {
                e.printStackTrace();
            }
        }
    }
}
