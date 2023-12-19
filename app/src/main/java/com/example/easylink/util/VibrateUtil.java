package com.example.easylink.util;

import android.app.Service;
import android.content.Context;
import android.os.VibrationEffect;
import android.os.Vibrator;

public class VibrateUtil {

    /**
     * 让手机振动milliseconds毫秒
     */
    public static void vibrate(Context context, long milliseconds) {
        Vibrator vib = (Vibrator) context.getSystemService(Service.VIBRATOR_SERVICE);
        if(vib.hasVibrator()){  //判断手机硬件是否有振动器
            //创建一个简单的振动效果，参数为持续时间和振动强度
            VibrationEffect effect = VibrationEffect.createOneShot(milliseconds, VibrationEffect.DEFAULT_AMPLITUDE);
            //使用振动效果来振动设备
            vib.vibrate(effect);
        }
    }



    /**
     * 取消震动
     */
    public static void virateCancle(Context context){
        //关闭震动
        Vibrator vib = (Vibrator) context.getSystemService(Service.VIBRATOR_SERVICE);
        vib.cancel();
    }
}

/*
// 开启震动
isVirating = true;
VirateUtil.virate(context, new long[]{100, 200, 100, 200}, 0)

//关闭震动
if (isVirating) {//防止多次关闭抛出异常，这里加个参数判断一下
    isVirating = false;
    VirateUtil.virateCancle(context);
}
 */
