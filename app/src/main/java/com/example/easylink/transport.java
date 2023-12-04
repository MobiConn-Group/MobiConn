package com.example.easylink;

import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.ToggleButton;

import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;

public class transport extends AppCompatActivity {

    @Override
    protected void onCreate(@Nullable Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_transport);

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


    }
}