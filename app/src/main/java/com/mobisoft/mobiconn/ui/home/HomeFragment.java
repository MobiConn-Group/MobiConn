package com.mobisoft.mobiconn.ui.home;

import static androidx.core.content.ContextCompat.getSystemService;

import android.annotation.SuppressLint;
import android.app.DownloadManager;
import android.database.Cursor;
import android.media.Ringtone;
import android.media.RingtoneManager;
import android.net.Uri;
import android.os.Bundle;
import android.os.Environment;
import android.provider.MediaStore;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;
import android.widget.Toast;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;

import com.mobisoft.mobiconn.MainActivity;
import com.mobisoft.mobiconn.api.Goal;
import com.mobisoft.mobiconn.api.MobiConnAPI;
import com.mobisoft.mobiconn.databinding.FragmentHomeBinding;

import java.util.List;
import java.util.Objects;

public class HomeFragment extends Fragment {

    private FragmentHomeBinding binding;

    private String getServerUrl() {
        String urlString = Objects.requireNonNull(binding.textInputHost.getText()).toString().trim();
        // 获取协议
        String protocol = urlString.split("://")[0];
        // 检查url是否包含协议
        if (!urlString.contains("://")) {
            protocol = "http";
        }
        return protocol + "://" + urlString;
    }

    public View onCreateView(@NonNull LayoutInflater inflater,
                             ViewGroup container, Bundle savedInstanceState) {
        HomeViewModel homeViewModel =
                new ViewModelProvider(this).get(HomeViewModel.class);

        binding = FragmentHomeBinding.inflate(inflater, container, false);
        View root = binding.getRoot();

        final TextView textView = binding.textHome;
        homeViewModel.getText().observe(getViewLifecycleOwner(), textView::setText);

        final Button buttonConnect = binding.buttonConnect;
        buttonConnect.setOnClickListener(this::buttonConnectOnClick);

        if (((MainActivity) requireActivity()).isConnected()) {
            binding.buttonConnect.setText("断开连接");
            binding.textInputHost.setEnabled(false);
        } else {
            binding.buttonConnect.setText("连接");
            binding.textInputHost.setEnabled(true);
        }

        return root;
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        binding = null;
    }

    @SuppressLint("UnspecifiedRegisterReceiverFlag")
    public void buttonConnectOnClick(View view) {
        if (((MainActivity) requireActivity()).isConnected()) {
            // 断开连接
            ((MainActivity) requireActivity()).disconnect();
            binding.buttonConnect.setText("连接");
            binding.buttonConnect.setEnabled(true);
            binding.textInputHost.setEnabled(true);
            return;
        }

        // 连接到服务器 http://192.168.239.209:25236
        final String serverUrl = getServerUrl();

        binding.buttonConnect.setText("连接中");
        binding.buttonConnect.setEnabled(false);
        binding.textInputHost.setEnabled(false);
        new Thread(() -> {
            List<Goal> goalList = MobiConnAPI.heartbeat(serverUrl);
            if (goalList == null) {
                requireActivity().runOnUiThread(() -> {
                    ((MainActivity) requireActivity()).disconnect();
                    if (binding != null) {
                        binding.buttonConnect.setText("连接失败");
                        binding.buttonConnect.setEnabled(true);
                        binding.textInputHost.setEnabled(true);
                    }
                });
                return;
            }
            Objects.requireNonNull(goalList);
            // 上传本机照片
            Cursor cursor = requireContext().getContentResolver().query(MediaStore.Images.Media.EXTERNAL_CONTENT_URI, null, null, null, null);
            int photoIndex = 0;
            if (cursor != null) {
                // 显示进度
                int finalPhotoIndex = photoIndex;
                requireActivity().runOnUiThread(() -> binding.buttonConnect.setText("正在上传本机照片 (" + cursor.getCount() + "张)"));
                MobiConnAPI.photoCount(serverUrl, cursor.getCount());
                while (cursor.moveToNext()) {
                    //获取图片的路径
                    int columnIndex = cursor.getColumnIndex(MediaStore.Images.Media.DATA);
                    if (columnIndex == -1) {
                        continue;
                    }
                    String path = cursor.getString(columnIndex);
                    if (path == null || path.length() == 0) {
                        continue;
                    }
                    MobiConnAPI.photoUpload(serverUrl, photoIndex, path);
                    ++photoIndex;
                }
            }
            ((MainActivity) requireActivity()).connect(serverUrl);
            requireActivity().runOnUiThread(() -> {
                if (binding != null) {
                    binding.buttonConnect.setText("断开连接");
                    binding.buttonConnect.setEnabled(true);
                    binding.textInputHost.setEnabled(false);
                }
            });

            while (((MainActivity) requireActivity()).isConnected()) {
                goalList = MobiConnAPI.heartbeat(serverUrl);
                if (goalList != null) {
                    for (Goal goal : goalList) {
                        if (Objects.equals(goal.getAction(), "ring")) {
                            int duration;
                            try {
                                duration = Integer.parseInt(goal.getInformation());
                            } catch (NumberFormatException e) {
                                duration = 5000;
                            }
                            Uri defaultRingtoneUri = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_RINGTONE);
                            Ringtone ringtone = RingtoneManager.getRingtone(requireContext(), defaultRingtoneUri);
                            ringtone.play();
                            int finalDuration = duration;
                            requireActivity().runOnUiThread(() -> {
                                String durationString;
                                if (finalDuration % 1000 == 0) {
                                    durationString = finalDuration / 1000 + " 秒";
                                } else {
                                    durationString = finalDuration + " 毫秒";
                                }
                                Toast.makeText(requireContext(), "即将响铃 " + Objects.requireNonNull(durationString), Toast.LENGTH_SHORT).show();
                            });
                            new Thread(() -> {
                                try {
                                    Thread.sleep(finalDuration);
                                    ringtone.stop();
                                } catch (InterruptedException ignored) {
                                }
                            }).start();

//                            String ringtone = Settings.System.DEFAULT_RINGTONE_URI.getPath();
//                            SoundPool ringPhone = new SoundPool(2, AudioManager.STREAM_RING, 1);
////                            int soundID = ringPhone.load(Settings.System.DEFAULT_RINGTONE_URI.getPath(), 1);
//                            int soundID = ringPhone.load(ringtone, 1);
//                            ringPhone.play(soundID, 0.99f, 0.99f, 1, 0, 1);

                        } else if (Objects.equals(goal.getAction(), "download")) {
                            new Thread(() -> {
//                                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
//                                    MobiConnAPI.fileDownload(serverUrl, goal.getInformation());
//                                }
                                // 创建下载请求
                                String fileName = goal.getInformation();
                                String urlString = serverUrl + "/file/download?" + "fileName=" + fileName;
                                DownloadManager.Request request = new DownloadManager.Request(Uri.parse(urlString));
                                request.setVisibleInDownloadsUi(true);
                                request.setNotificationVisibility(DownloadManager.Request.VISIBILITY_VISIBLE_NOTIFY_COMPLETED);
                                String name = fileName.substring(Math.max(fileName.lastIndexOf('/'), fileName.lastIndexOf('\\')) + 1);
//                                request.addRequestHeader("Content-Length", "100000000");
                                request.setDestinationInExternalPublicDir(Environment.DIRECTORY_DOWNLOADS, name);
//                                request.setDestinationUri(Uri.fromFile(new File(requireContext().getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS) + name)));
                                request.setAllowedNetworkTypes(DownloadManager.Request.NETWORK_MOBILE | DownloadManager.Request.NETWORK_WIFI);
                                DownloadManager downloadManager = getSystemService(requireContext(), DownloadManager.class);
                                //noinspection unused
                                long downloadID = Objects.requireNonNull(downloadManager).enqueue(request);
                                requireActivity().runOnUiThread(() -> Toast.makeText(requireContext(), "下载文件 " + name + " 已加入下载列表", Toast.LENGTH_SHORT).show());
//                                NotificationCompat.Builder builder = new NotificationCompat.Builder(requireContext(), "channel_id")
//                                        .setContentTitle("下载通知")
//                                        .setContentText(goal.getInformation() + "下载完成")
//                                        .setPriority(NotificationCompat.PRIORITY_DEFAULT)
//                                        .setAutoCancel(true);
//                                Uri defaultRingtoneUri = RingtoneManager.getDefaultUri(RingtoneManager.TYPE_RINGTONE);
//                                Ringtone ringtone = RingtoneManager.getRingtone(requireContext(), defaultRingtoneUri);
//                                ringtone.play();
//                                int finalDuration = 20000;
//                                try {
//                                    Thread.sleep(finalDuration);
//                                } catch (InterruptedException e) {
//                                    throw new RuntimeException(e);
//                                }
//                                ringtone.stop();
                            }).start();
                        }
                    }
                }
                try {
                    Thread.sleep(100);
                } catch (InterruptedException e) {
                    throw new RuntimeException(e);
                }
            }
        }).start();
    }
}