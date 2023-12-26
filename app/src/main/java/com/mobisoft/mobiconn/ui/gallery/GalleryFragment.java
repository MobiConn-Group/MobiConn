package com.mobisoft.mobiconn.ui.gallery;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;

import com.mobisoft.mobiconn.MainActivity;
import com.mobisoft.mobiconn.api.MobiConnAPI;
import com.mobisoft.mobiconn.databinding.FragmentGalleryBinding;

public class GalleryFragment extends Fragment {

    private FragmentGalleryBinding binding;

    public View onCreateView(@NonNull LayoutInflater inflater,
                             ViewGroup container, Bundle savedInstanceState) {
        GalleryViewModel galleryViewModel =
                new ViewModelProvider(this).get(GalleryViewModel.class);

        binding = FragmentGalleryBinding.inflate(inflater, container, false);
        View root = binding.getRoot();

        final TextView textView = binding.textGallery;
        galleryViewModel.getText().observe(getViewLifecycleOwner(), textView::setText);

        final Button buttonStitchLeft = binding.buttonSwitchLeft;
        buttonStitchLeft.setOnClickListener(this::buttonSwitchLeftOnClick);
        if (((MainActivity) requireActivity()).isConnected()) {
            buttonStitchLeft.setEnabled(true);
        }

        final Button buttonStitchRight = binding.buttonSwitchRight;
        buttonStitchRight.setOnClickListener(this::buttonSwitchRightOnClick);
        if (((MainActivity) requireActivity()).isConnected()) {
            buttonStitchRight.setEnabled(true);
        }

        return root;
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        binding = null;
    }

    public void buttonSwitchLeftOnClick(View view) {
        new Thread(() -> MobiConnAPI.powerPoint(MainActivity.getServerUrl(), 0)).start();
    }

    public void buttonSwitchRightOnClick(View view) {
        new Thread(() -> MobiConnAPI.powerPoint(MainActivity.getServerUrl(), 1)).start();
    }
}