package com.example.easylink;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.ToggleButton;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;

import java.io.IOException;

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

        String username = getIntent().getStringExtra("username");
        String password = getIntent().getStringExtra("password");
        String url = getIntent().getStringExtra("url");

        //构造请求体，使用FormBody携带键值对参数
        RequestBody requestBody = new FormBody.Builder()
                .add("username", username)
                .add("password", password)
                .build();

        //构造请求对象，使用POST方法
        Request request = new Request.Builder()
                .url(url)
                .post(requestBody)
                .build();

        link_webserver(request);

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

    private void link_webserver(Request request) {
        WebSocketListener listener = new WebSocketListener() {
            @Override
            public void onOpen(WebSocket webSocket, Response response) {
                System.out.println("WebSocket connection opened");

            }

            @Override
            public void onMessage(WebSocket webSocket, String text) {
                System.out.println("Received message: " + text);
                // 在这里处理接收到的消息
            }

            @Override
            public void onClosed(WebSocket webSocket, int code, String reason) {
                System.out.println("WebSocket connection closed");
            }

            @Override
            public void onFailure(WebSocket webSocket, Throwable t, Response response) {
                System.out.println("WebSocket connection failed");
            }
        };

        WebSocket webSocket = client.newWebSocket(request, listener);

        // 可选：发送消息到服务器
        webSocket.send("Hello, server!");
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        // 关闭WebSocket连接
        webSocket.cancel();
    }
}