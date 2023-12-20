package com.example.easylink;

import androidx.activity.result.ActivityResult;
import androidx.activity.result.ActivityResultCallback;
import androidx.activity.result.ActivityResultLauncher;
import androidx.activity.result.contract.ActivityResultContracts;
import androidx.annotation.Nullable;
import androidx.appcompat.app.AppCompatActivity;
import androidx.core.app.ActivityCompat;
import androidx.core.content.ContextCompat;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.ContentResolver;
import android.content.Intent;
import android.content.pm.PackageManager;
import android.database.Cursor;
import android.net.Uri;
import android.os.Bundle;
import android.provider.MediaStore;
import android.provider.OpenableColumns;
import android.util.Log;
import android.view.View;
import android.widget.Button;
import android.widget.ImageButton;
import android.widget.TextView;
import android.widget.Toast;

import com.example.easylink.util.FileUtil;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.util.List;

import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.MediaType;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;

public class file_transport extends AppCompatActivity {

    private static final int REQUEST_PERMISSION_CODE = 100;
    private static final int REQUEST_CODE = 1;
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_file_transport);

        Button fileChooseButton = findViewById(R.id.fileChooseButton);

        // 检查是否有读取媒体存储的权限，如果没有，就请求权限
        if (ContextCompat.checkSelfPermission(this, android.Manifest.permission.READ_EXTERNAL_STORAGE) != PackageManager.PERMISSION_GRANTED) {
            ActivityCompat.requestPermissions(this, new String[]{android.Manifest.permission.READ_EXTERNAL_STORAGE}, REQUEST_PERMISSION_CODE);
        }

        Resource resource =  (Resource)getApplication();
        String ip = resource.getIp();
        String username = resource.getUsername();
        TextView textIp = findViewById(R.id.file_trans_ip);
        TextView textUsername = findViewById(R.id.file_trans_username);
        textIp.setText("IP: "+ip);
        textUsername.setText("用户名: "+username);
        String url = "http://" + ip + ":25236";
        fileChooseButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                Intent intent = new Intent();
                intent.setClass(file_transport.this, file_open.class);
                startActivity(intent);

                List<File> uploadFiles = resource.getUploadFiles();
                StringBuilder fileShowText = new StringBuilder("已添加的文件:\n");
                for (File file : uploadFiles) {
                    fileShowText.append(file.getName()).append("\n");
                }
                TextView fileShow = findViewById(R.id.file_show);
                fileShow.setText(fileShowText);
            }
        });

        Button fileTransButton = findViewById(R.id.fileTransButton);
        fileTransButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View view) {


            }
        });

        ImageButton back = findViewById(R.id.file_back);
        back.setOnClickListener(view -> file_transport.this.finish());
    }


}