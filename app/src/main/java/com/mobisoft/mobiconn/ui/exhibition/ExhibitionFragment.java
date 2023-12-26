package com.mobisoft.mobiconn.ui.exhibition;

import android.content.ContentResolver;
import android.database.Cursor;
import android.net.Uri;
import android.os.Bundle;
import android.provider.MediaStore;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;

import androidx.activity.result.ActivityResultLauncher;
import androidx.activity.result.contract.ActivityResultContracts;
import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;

import com.mobisoft.mobiconn.MainActivity;
import com.mobisoft.mobiconn.api.MobiConnAPI;
import com.mobisoft.mobiconn.databinding.FragmentExhibitionBinding;

import java.io.IOException;
import java.io.InputStream;
import java.util.Objects;

public class ExhibitionFragment extends Fragment {
    private static ExhibitionFragment instance = null;

    private FragmentExhibitionBinding binding;

    ActivityResultLauncher<String> mGetContent;

    public View onCreateView(@NonNull LayoutInflater inflater, ViewGroup container, Bundle savedInstanceState) {
        instance = this;

        ExhibitionViewModel ExhibitionViewModel = new ViewModelProvider(this).get(ExhibitionViewModel.class);

        binding = FragmentExhibitionBinding.inflate(inflater, container, false);
        View root = binding.getRoot();

        mGetContent = registerForActivityResult(new ActivityResultContracts.GetContent(), uri -> {
            if (uri == null) {
                return;
            }

            String displayName = getDisplayName(uri);
            binding.buttonUploadFile.setText("正在上传");
            MainActivity.setUploading(true);
            binding.buttonUploadFile.setEnabled(false);
            new Thread(() -> {
                try (InputStream inputStream = requireContext().getContentResolver().openInputStream(uri)) {
                    MobiConnAPI.fileUpload(MainActivity.getServerUrl(), displayName, Objects.requireNonNull(inputStream));
                } catch (IOException e) {
                    MainActivity.setUploading(false);
                    if (instance == null) {
                        return;
                    }
                    instance.requireActivity().runOnUiThread(() -> {
                        instance.binding.buttonUploadFile.setText("上传失败");
                        instance.binding.buttonUploadFile.setEnabled(true);
                    });
                    return;
                }
                MainActivity.setUploading(false);
                if (instance == null) {
                    return;
                }
                instance.requireActivity().runOnUiThread(() -> {
                    instance.binding.buttonUploadFile.setText("上传成功");
                    instance.binding.buttonUploadFile.setEnabled(true);
                });
            }).start();
        });


        final TextView textView = binding.textExhibition;
        ExhibitionViewModel.getText().observe(getViewLifecycleOwner(), textView::setText);

        final Button buttonUploadFile = binding.buttonUploadFile;
        buttonUploadFile.setOnClickListener(this::buttonUploadFileOnClick);

        if (MainActivity.isUploading()) {
            buttonUploadFile.setText("正在上传");
            buttonUploadFile.setEnabled(false);
        } else {
            buttonUploadFile.setText("上传文件");
            buttonUploadFile.setEnabled(true);
        }

        return root;
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        instance = null;
        binding = null;
    }

    public void buttonUploadFileOnClick(View view) {
        try {
            mGetContent.launch("*/*");
        } catch (android.content.ActivityNotFoundException e) {
            binding.buttonUploadFile.setText(e.getMessage());
        }
    }

    private String getDisplayName(Uri uri) {
        ContentResolver resolver = requireContext().getContentResolver();
        Cursor cursor = resolver.query(uri, null, null, null, null);
        if (cursor == null) {
            // 未查询到，说明为普通文件，可直接通过URI获取文件路径
            return uri.getPath();
        }
        if (cursor.moveToFirst()) {
            // 多媒体文件，从数据库中获取文件的真实路径
            String[] s = cursor.getColumnNames();
            int index = cursor.getColumnIndex(MediaStore.Files.FileColumns.DISPLAY_NAME);
            if (index == -1) {
                return null;
            }
            return cursor.getString(index);
        }
        cursor.close();
        return null;
    }
}