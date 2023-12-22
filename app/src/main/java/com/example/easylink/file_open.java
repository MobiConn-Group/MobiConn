package com.example.easylink;

import android.annotation.SuppressLint;
import android.app.Activity;
import android.content.ContentResolver;
import android.content.Intent;
import android.database.Cursor;
import android.net.Uri;
import android.os.Bundle;
import android.provider.OpenableColumns;
import android.widget.TextView;
import android.widget.Toast;

import androidx.activity.result.ActivityResultLauncher;
import androidx.activity.result.contract.ActivityResultContracts;
import androidx.appcompat.app.AppCompatActivity;

import com.example.easylink.util.FileUtil;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.nio.file.Files;
import java.util.List;

public class file_open extends AppCompatActivity {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_file_open);

        Resource resource =  (Resource)getApplication();
        String ip = resource.getIp();
        String url = "http://" + ip + ":25236";

//        openFile(url);
        ActivityResultLauncher<Intent> launcher = registerForActivityResult( new ActivityResultContracts.StartActivityForResult(), result -> {
            if (result.getData() != null && result.getResultCode() == Activity.RESULT_OK) {
                final Intent data = result.getData();
                Uri uri = data.getData(); // 获取用户选择文件的URI
                // 通过ContentProvider查询文件路径
                ContentResolver resolver = getContentResolver();
                assert uri != null;
                Cursor cursor = resolver.query(uri, null, null, null, null);
                if (cursor != null && cursor.moveToFirst()) {
                    // 获取文件的原始名称
                    @SuppressLint("Range")
                    String displayName = cursor.getString(cursor.getColumnIndex(OpenableColumns.DISPLAY_NAME));
                    // 获取文件的MIME类型
                    String mimeType = resolver.getType(uri);

                    try {
                        InputStream inputStream = resolver.openInputStream(uri);
                        if (inputStream != null) {
                            File file = createFileFromInputStream(inputStream, displayName);
                            List<String> uploadFiles = resource.getUploadFiles();
                            uploadFiles.add(file.getName());
                            resource.setUploadFiles(uploadFiles);
                            FileUtil.uploadFile(file, url);
                        } else {
                            Toast.makeText(getApplicationContext(), "输入流获取失败", Toast.LENGTH_SHORT).show();
                        }
                    } catch (IOException e) {
                        e.printStackTrace();
                    }

                    cursor.close();
                } else {
                    Toast.makeText(getApplicationContext(), "文件信息查询失败", Toast.LENGTH_SHORT).show();
                }
            } else {
                Toast.makeText(getApplicationContext(), "请求失败", Toast.LENGTH_SHORT).show();
            }

        });

        Intent intent = new Intent(Intent.ACTION_GET_CONTENT);
        intent.addCategory(Intent.CATEGORY_OPENABLE);
        intent.setType("*/*");

        launcher.launch(intent);
        file_open.this.finish();
    }

    private void openFile(String url) {

    }

    // 从输入流创建临时文件
    private File createFileFromInputStream(InputStream inputStream, String displayName) throws IOException {
        File file = new File(getCacheDir(), displayName);
        OutputStream outputStream = Files.newOutputStream(file.toPath());
        byte[] buffer = new byte[4 * 1024];
        int bytesRead;
        while ((bytesRead = inputStream.read(buffer)) != -1) {
            outputStream.write(buffer, 0, bytesRead);
        }
        outputStream.flush();
        outputStream.close();
        inputStream.close();
        return file;
    }
}
