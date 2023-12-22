package com.example.easylink;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.ImageButton;
import android.widget.TextView;
import android.widget.Toast;

import java.io.IOException;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.FormBody;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class text_transport extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_text_transport);

        Resource resource =  (Resource)getApplication();
        String ip = resource.getIp();
        String username = resource.getUsername();
        TextView textIp = findViewById(R.id.text_trans_ip);
        TextView textUsername = findViewById(R.id.text_trans_username);
        textIp.setText("IP: "+ip);
        textUsername.setText("用户名: "+username);
        String url = "http://" + ip + ":25236";

        Button textTransButton = findViewById(R.id.textTransButton);
        textTransButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                EditText editText = findViewById(R.id.text_input);
                String text = editText.getText().toString();

                //创建一个OkHttpClient对象
                OkHttpClient client = new OkHttpClient();

                //创建一个请求体，指定字符串和媒体类型
                RequestBody body = new FormBody.Builder()
                        .add("text", text)
                        .build();

                //创建一个请求，指定url和请求方法
                Request request = new Request.Builder()
                        .url(url)
                        .post(body)
                        .build();

                //执行请求，获取响应
                client.newCall(request).enqueue(new Callback() {
                    @Override
                    public void onFailure(Call call, IOException e) {
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                //请求失败时的处理
                                e.printStackTrace();
                                Toast.makeText(getApplicationContext(), "文字传输失败", Toast.LENGTH_LONG).show();
                            }
                        });
                    }

                    @Override
                    public void onResponse(Call call, Response response) throws IOException {
                        runOnUiThread(new Runnable() {
                            @Override
                            public void run() {
                                //请求成功时的处理
                                if (response.isSuccessful()) {
                                    //获取响应体的字符串
                                    try {
                                        String result = response.body().string();
                                    } catch (IOException e) {
                                        throw new RuntimeException(e);
                                    }
                                    editText.setText("");
                                } else {
                                    Toast.makeText(getApplicationContext(), "请求失败："+response.code(), Toast.LENGTH_LONG).show();
                                }
                            }
                        });
                    }
                });

            }
        });

        ImageButton back = findViewById(R.id.text_back);
        back.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                text_transport.this.finish();
            }
        });
    }
}