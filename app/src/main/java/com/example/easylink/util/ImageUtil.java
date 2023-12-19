package com.example.easylink.util;

import java.io.File;
import java.util.Arrays;
import java.util.Comparator;

public class ImageUtil {
    public static File[] getImagesFromFolder(String folderPath) {
        File folder = new File(folderPath);
        File[] files = folder.listFiles();

        // 按照修改时间进行排序
        assert files != null;
        Arrays.sort(files, new Comparator<File>() {
            @Override
            public int compare(File file1, File file2) {
                long lastModified1 = file1.lastModified();
                long lastModified2 = file2.lastModified();
                return Long.compare(lastModified1, lastModified2);
            }
        });

        return files;
    }
}
