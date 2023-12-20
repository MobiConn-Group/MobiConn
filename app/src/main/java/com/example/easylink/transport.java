package com.example.easylink;

import android.annotation.SuppressLint;
import android.content.ContentResolver;
import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.net.Uri;
import android.os.Bundle;
import android.os.Environment;
import android.provider.MediaStore;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;
import android.widget.Toast;
import android.widget.ToggleButton;

import androidx.annotation.NonNull;
import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import androidx.core.content.PackageManagerCompat;

import com.example.easylink.util.AlbumUtil;
import com.example.easylink.util.ImageUtil;
import com.example.easylink.util.MediaUtil;
import com.example.easylink.util.VibrateUtil;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.util.ArrayList;
import java.util.List;
import java.util.Timer;
import java.util.TimerTask;
import java.util.logging.Logger;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.FormBody;
import okhttp3.MediaType;
import okhttp3.MultipartBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;

public class transport extends AppCompatActivity {



    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_transport);

        Resource resource =  (Resource)getApplication();

        String ip = resource.getIp();
        String username = resource.getUsername();
        TextView textIp = findViewById(R.id.textip);
        TextView textUsername = findViewById(R.id.textusername);
        textIp.setText("IP: "+ip);
        textUsername.setText("用户名: "+username);


        runServer(ip);

        Button textTransButton = findViewById(R.id.textTransButton);
        textTransButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Intent intent = new Intent();
                intent.setClass(transport.this, text_transport.class);
                startActivity(intent);

            }
        });

        Button fileTransBottom = findViewById(R.id.fileTransButton);
        fileTransBottom.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Intent intent = new Intent();
                intent.setClass(transport.this, file_transport.class);
                startActivity(intent);

            }
        });

        Button logout = findViewById(R.id.logout);
        logout.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                transport.this.finish();
            }
        });
    }

    private void runServer(String ip) {
        String url = "http://" + ip + ":25236";
        // 创建OkHttpClient实例
        OkHttpClient client = new OkHttpClient();

        // 创建一个线程
        Thread thread = new Thread(() -> {
            while (true) {
                try {
                    // 创建GET请求
                    Resource token =  (Resource)getApplication();
                    Request request = new Request.Builder()
                            .url(url) // 替换为你的服务器地址
                            .header("Token", token.getToken())
                            .build();

                    // 发送请求并获取响应
                    Response response = client.newCall(request).execute();

                    // 处理响应
                    if (response.isSuccessful()) {
                        // 响应成功
                        assert response.body() != null;
                        String responseBody = response.body().string();
                        try {
                            JSONObject jsonObject = new JSONObject(responseBody);
                            if (jsonObject.getString("Status").equals("Success")) {
                                JSONArray goals = jsonObject.getJSONArray("Goals");
                                for (int i = 0; i < goals.length(); ++i) {
                                    JSONObject goal = goals.getJSONObject(i);
                                    //TODO 完成各种类型的实现
                                    String action = goal.getString("Action");
                                    switch (action) {
                                        case "download" : {

                                        }
                                        case "ring" : {
                                            MediaUtil.ring(this);
                                        }
                                        case "vibrate" : {
                                            VibrateUtil.vibrate(this, 5000);
                                        }
                                        case "upload" : {

                                        }
                                        case "album" : {
                                            AlbumUtil.album(url, (Resource)getApplication());
                                        }
                                        case "photo" : {

                                        }
                                    }
                                }
                            } else {
                                Logger.getGlobal().warning("Get Fault");
                            }
                        } catch (JSONException e) {
                            e.printStackTrace();
                        }
                        System.out.println("Response: " + responseBody);
                    } else {
                        // 响应失败
                        System.out.println("Request failed");
                    }

                    // 等待100ms
                    Thread.sleep(100);
                } catch (Exception e) {
                    e.printStackTrace();
                }
            }
        });

        // 启动线程
        thread.start();
    }


}