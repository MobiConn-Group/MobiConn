package com.mobisoft.mobiconn.ui.slideshow;

import android.os.Bundle;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.Button;
import android.widget.TextView;

import androidx.annotation.NonNull;
import androidx.fragment.app.Fragment;
import androidx.lifecycle.ViewModelProvider;

import com.google.android.material.textfield.TextInputEditText;
import com.mobisoft.mobiconn.MainActivity;
import com.mobisoft.mobiconn.api.MobiConnAPI;
import com.mobisoft.mobiconn.databinding.FragmentSlideshowBinding;

import java.util.Objects;

public class SlideshowFragment extends Fragment {

    private FragmentSlideshowBinding binding;

    public View onCreateView(@NonNull LayoutInflater inflater,
                             ViewGroup container, Bundle savedInstanceState) {
        SlideshowViewModel slideshowViewModel =
                new ViewModelProvider(this).get(SlideshowViewModel.class);

        binding = FragmentSlideshowBinding.inflate(inflater, container, false);
        View root = binding.getRoot();

        final TextView textView = binding.textSlideshow;
        slideshowViewModel.getText().observe(getViewLifecycleOwner(), textView::setText);

        final TextInputEditText textInputCursorText = binding.textInputCursorText;
        textInputCursorText.setEnabled(((MainActivity) requireActivity()).isConnected());

        final Button buttonSendCursorText = binding.buttonSendCursorText;
        buttonSendCursorText.setOnClickListener(this::buttonSendCursorTextOnClick);
        buttonSendCursorText.setEnabled(((MainActivity) requireActivity()).isConnected());

        return root;
    }

    @Override
    public void onDestroyView() {
        super.onDestroyView();
        binding = null;
    }

    public void buttonSendCursorTextOnClick(View view) {
        new Thread(() -> MobiConnAPI.sendCursorText(MainActivity.getServerUrl(), Objects.requireNonNull(binding.textInputCursorText.getText()).toString())).start();
    }
}