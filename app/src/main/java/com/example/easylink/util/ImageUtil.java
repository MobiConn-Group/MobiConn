package com.example.easylink.util;

import android.annotation.SuppressLint;
import android.content.ContentResolver;
import android.database.Cursor;
import android.net.Uri;
import android.provider.MediaStore;
import android.util.Log;

import java.io.File;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;

public class ImageUtil {
    public static void getImages(List<File> imageFiles, ContentResolver contentResolver) {
        // 定义一个Uri，表示媒体库中的图片
        Uri uri = MediaStore.Images.Media.EXTERNAL_CONTENT_URI;
        // 定义一个Cursor，用来查询图片信息
        Cursor cursor = contentResolver.query(uri, null, null, null, null);
        // 判断Cursor是否为空
        if (cursor != null) {
            // 遍历Cursor，获取每一张图片的路径，并创建File对象，添加到列表中
            while (cursor.moveToNext()) {
                // 获取图片的路径
                @SuppressLint("Range")
                String path = cursor.getString(cursor.getColumnIndex(MediaStore.Images.Media.DATA));
                // 创建File对象
                File file = new File(path);
                // 添加到列表中
                imageFiles.add(file);
                // 打印日志，显示图片的路径
                Log.d("MainActivity", "getImages: " + path);
            }
            // 关闭Cursor
            cursor.close();
        }
        // 打印日志，显示图片文件的数量
        Log.d("MainActivity", "getImages: imageFiles.size = " + imageFiles.size());
    }
}