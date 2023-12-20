package com.example.easylink;

import androidx.appcompat.app.AppCompatActivity;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;

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