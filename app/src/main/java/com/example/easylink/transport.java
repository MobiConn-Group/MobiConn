package com.example.easylink;

import android.annotation.SuppressLint;
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;
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

        String url = getIntent().getStringExtra("url");
        String username = getIntent().getStringExtra("username");
        TextView textip = findViewById(R.id.textip);
        textip.setText("IP: "+url);
        TextView textusername = findViewById(R.id.textusername);
        textusername.setText("用户名: "+username);
        link_webserver();

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

    private void link_webserver() {

    }

}