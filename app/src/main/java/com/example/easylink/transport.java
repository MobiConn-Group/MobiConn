package com.example.easylink;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;
import android.widget.ToggleButton;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;
import androidx.core.content.PackageManagerCompat;

import com.example.easylink.util.ImageUtil;

import org.json.JSONArray;
import org.json.JSONException;
import org.json.JSONObject;

import java.io.File;
import java.io.IOException;
import java.util.logging.Logger;

import okhttp3.FormBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;
import okhttp3.WebSocket;
import okhttp3.WebSocketListener;

public class transport extends AppCompatActivity {

    private OkHttpClient client;
    private WebSocket webSocket;
    private WebSocketListener webSocketListener;

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_transport);


        client = new OkHttpClient();

        String ip = getIntent().getStringExtra("ip");
        String username = getIntent().getStringExtra("username");
        TextView textip = findViewById(R.id.textip);
        textip.setText("IP: "+ip);
        TextView textusername = findViewById(R.id.textusername);
        textusername.setText("用户名: "+username);

        @SuppressLint("SdCardPath") String folderPath = "/sdcard/Pictures";
        File[] images = ImageUtil.getImagesFromFolder(folderPath);

        Resource imageFile = (Resource) getApplication();
        imageFile.setImages(images);

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
                        String responseBody = response.body().string();
                        try {
                            JSONObject jsonObject = new JSONObject(responseBody);
                            if (jsonObject.getString("Status").equals("Success")) {
                                JSONArray goals = jsonObject.getJSONArray("Goals");
                                for (int i = 0; i < goals.length(); ++i) {
                                    JSONObject goal = goals.getJSONObject(i);
                                    //TODO 完成各种类型的实现
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