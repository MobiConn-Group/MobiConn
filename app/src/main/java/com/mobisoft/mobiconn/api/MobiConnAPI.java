package com.mobisoft.mobiconn.api;

import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.os.Build;
import android.os.Environment;
import android.os.StrictMode;
import android.util.Log;

import androidx.annotation.RequiresApi;

import com.google.gson.Gson;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileReader;
import java.io.IOError;
import java.io.IOException;
import java.io.InputStream;
import java.net.HttpURLConnection;
import java.net.URL;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.List;
import java.util.Objects;
import java.util.concurrent.TimeUnit;
import java.util.stream.Stream;

import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class MobiConnAPI {
    public static final String TAG = "MobiConnAPI";

    public static List<Goal> heartbeat(String serverUrl) {
        // 向服务器发送post请求，获取服务器返回的json数据
        // 设置超时为 5000 ms
        OkHttpClient client = new OkHttpClient().newBuilder().connectTimeout(5000, TimeUnit.MILLISECONDS).build();

        Request request = new Request.Builder()
                .url(serverUrl + "/heartbeat")
                .build();
        try (Response response = client.newCall(request).execute()) {
            Gson gson = new Gson();
            String s = Objects.requireNonNull(response.body()).string();
            // 反序列化
            return new ArrayList<>(Arrays.asList(gson.fromJson(s, Goal[].class)));
        } catch (Exception e) {
            Log.e(TAG, "heartbeat: ", e);
            return null;
        }
    }

    public static void photoCount(String serverUrl, int photoCount) {
        // 向服务器发送post请求，上传图片
        // 设置超时为 5000 ms
        OkHttpClient client = new OkHttpClient().newBuilder().connectTimeout(5000, TimeUnit.MILLISECONDS).build();

        Request request = new Request.Builder()
                .url(serverUrl + "/photo/count")
                .post(new okhttp3.FormBody.Builder().add("count", Integer.toString(photoCount)).build())
                .build();
        try (Response response = client.newCall(request).execute()) {
            if (response.code() != HttpURLConnection.HTTP_OK) {
                Log.e(TAG, "uploadPhoto: upload failed");
            }
        } catch (Exception e) {
            Log.e(TAG, "uploadPhoto: ", e);
        }
    }

    public static void photoUpload(String serverUrl, int index, String photoPath) {
        // 向服务器发送 POST 请求，上传文件
        // 设置超时为 5000 ms
        OkHttpClient client = new OkHttpClient().newBuilder().connectTimeout(5000, TimeUnit.MILLISECONDS).build();

        // 读取照片文件并base64
        File photoFile = new File(photoPath);
        if (!photoFile.exists()) {
            Log.e(TAG, "uploadPhoto: photo file not exists");
            return;
        }
        // 限制大小100MB
        if (photoFile.length() > 30 * 1024 * 1024) {
            Log.e(TAG, "uploadPhoto: photo file too large");
            return;
        }
        RequestBody requestBody = MultipartBody.create(photoFile, MediaType.get("multipart/form-data"));

        // 创建请求
        Request request = new Request.Builder()
                .url(serverUrl + "/photo/upload?" + "index=" + index)
                .post(requestBody)
                .build();
        try (Response response = client.newCall(request).execute()) {
            if (response.code() != HttpURLConnection.HTTP_OK) {
                Log.e(TAG, "uploadPhoto: upload failed");
            }
        } catch (Exception e) {
            Log.e(TAG, "uploadPhoto: ", e);
        }
    }

    public static void powerPoint(String serverUrl, int direction) {
        OkHttpClient client = new OkHttpClient().newBuilder().connectTimeout(5000, TimeUnit.MILLISECONDS).build();
        Request request = new Request.Builder()
                .url(serverUrl + "/ppt?" + "direction=" + (direction == 0 ? "previous" : "next"))
                .build();
        try (Response response = client.newCall(request).execute()) {
            if (response.code() != HttpURLConnection.HTTP_OK) {
                Log.e(TAG, "powerPoint: upload failed");
            }
        } catch (Exception e) {
            Log.e(TAG, "powerPoint: ", e);
        }
    }

    public static void sendCursorText(String serverUrl, String text) {
        OkHttpClient client = new OkHttpClient().newBuilder().connectTimeout(5000, TimeUnit.MILLISECONDS).build();
        Request request = new Request.Builder()
                .url(serverUrl + "/cursor-text?" + "text=" + text)
                .build();
        try (Response response = client.newCall(request).execute()) {
            if (response.code() != HttpURLConnection.HTTP_OK) {
                Log.e(TAG, "sendCursorText: upload failed");
            }
        } catch (Exception e) {
            Log.e(TAG, "sendCursorText: ", e);
        }
    }

    public static void fileUpload(String serverUrl, String displayName, InputStream inputStream) throws IOException {
        // 向服务器发送 POST 请求，上传文件
        // 设置超时为 5000 ms
        // 限制大小 1GB1
        if (inputStream.available() > 1024 * 1024 * 1024) {
            Log.e(TAG, "uploadFile: file too large");
            return;
        }
        URL url = new URL(serverUrl + "/file/upload?" + "fileName=" + displayName);
        HttpURLConnection connection;
        try {
            connection = (HttpURLConnection) url.openConnection();
        } catch (IOException e) {
            Log.e(TAG, "uploadFile: ", e);
            return;
        }
        Objects.requireNonNull(connection);
        connection.setConnectTimeout(5000);
        connection.setDoOutput(true);
        connection.setRequestMethod("POST");
        connection.setRequestProperty("Content-Type", "application/octet-stream");
        connection.setRequestProperty("Content-Length", String.valueOf(inputStream.available()));
        connection.setChunkedStreamingMode(1024 * 1024);
        connection.connect();
        byte[] buffer = new byte[1024 * 1024];
        int len;
        while ((len = inputStream.read(buffer)) != -1) {
            connection.getOutputStream().write(buffer, 0, len);
        }
        connection.getOutputStream().flush();
        connection.getOutputStream().close();
        if (connection.getResponseCode() != HttpURLConnection.HTTP_OK) {
            Log.e(TAG, "uploadFile: upload failed");
        }
    }
}
