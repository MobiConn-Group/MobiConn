package com.example.easylink;

import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import android.content.Context;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;

import com.example.easylink.util.ImageUtil;
import com.example.easylink.util.MediaUtil;
import com.example.easylink.util.VibrateUtil;

import java.io.File;
import java.util.ArrayList;
import java.util.List;

public class MainActivity extends AppCompatActivity {
    // 定义一个常量，表示请求权限的代码
    private static final int REQUEST_PERMISSION_CODE = 100;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);

        // 检查是否有读取媒体存储的权限，如果没有，就请求权限
        if (ContextCompat.checkSelfPermission(this, android.Manifest.permission.READ_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
            ActivityCompat.requestPermissions(this, new String[]{android.Manifest.permission.READ_EXTERNAL_STORAGE}, REQUEST_PERMISSION_CODE);
        }
        // 定义一个列表，用来存储图片文件
        List<File> imageFiles = new ArrayList<>();

        ImageUtil.getImages(imageFiles, getContentResolver());
        Resource resource =  (Resource)getApplication();
        resource.setImages(imageFiles);

        Button button = findViewById(R.id.jumpButton);
        button.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {
                Intent intent = new Intent();
                intent.setClass(MainActivity.this, computer_link.class);
                startActivity(intent);

            }
        });
    }


}