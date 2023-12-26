package com.mobisoft.mobiconn.ui.exhibition;

import androidx.lifecycle.LiveData;
import androidx.lifecycle.MutableLiveData;
import androidx.lifecycle.ViewModel;

public class ExhibitionViewModel extends ViewModel {

    private final MutableLiveData<String> mText;

    public ExhibitionViewModel() {
        mText = new MutableLiveData<>();
        mText.setValue("上传文件");
    }

    public LiveData<String> getText() {
        return mText;
    }
}