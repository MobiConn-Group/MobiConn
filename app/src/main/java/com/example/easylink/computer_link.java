package com.example.easylink;//导入所需的库
import android.content.Intent;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.Toast;
import androidx.appcompat.app.AppCompatActivity;
import okhttp3.Call;
import okhttp3.Callback;
import okhttp3.FormBody;
import okhttp3.OkHttpClient;
import okhttp3.Request;
import okhttp3.RequestBody;
import okhttp3.Response;
import java.io.IOException;

//定义一个LoginActivity类，继承自AppCompatActivity
public class computer_link extends AppCompatActivity {

    //定义控件变量
    private EditText et_ip; //输入IP地址的文本框
    private EditText et_username; //输入用户名的文本框
    private EditText et_password; //输入密码的文本框

    //定义OkHttpClient对象
    private OkHttpClient client;

    //重写onCreate方法，初始化控件和事件
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_computer_link); //设置布局文件

        //初始化控件
        et_ip = findViewById(R.id.loginIp);
        et_username = findViewById(R.id.username);
        et_password = findViewById(R.id.userPassword);
        //登录按钮
        Button btn_login = findViewById(R.id.signInButton);

        //初始化OkHttpClient对象
        client = new OkHttpClient();

        //设置登录按钮的点击事件
        btn_login.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {

                Intent intent = new Intent();
                intent.setClass(computer_link.this, transport.class);
                startActivity(intent);

                //获取用户输入的IP地址，用户名和密码
//                String ip = et_ip.getText().toString().trim();
//                String username = et_username.getText().toString().trim();
//                String password = et_password.getText().toString().trim();
//                //判断输入是否为空
//                if (ip.isEmpty() || username.isEmpty()) {
//                    //提示用户输入不能为空
//                    Toast.makeText(computer_link.this, "请输入IP地址，用户名", Toast.LENGTH_SHORT).show();
//                } else {
//                    //执行登录操作
//                    login(ip, username, password);
//                }
            }
        });
    }

    //定义一个login方法，用于向服务器发送登录请求
    private void login(String ip, String username, String password) {
        //构造请求地址，假设服务器提供了一个login接口，接收IP，username和password参数
        String url = "http://" + ip + "/login";

        //构造请求体，使用FormBody携带键值对参数
        RequestBody requestBody = new FormBody.Builder()
                .add("username", username)
                .add("password", password)
                .build();

        //构造请求对象，使用POST方法
        Request request = new Request.Builder()
                .url(url)
                .post(requestBody)
                .build();

        //通过OkHttpClient对象构造Call对象
        Call call = client.newCall(request);

        //通过Call对象的enqueue方法发送异步请求，并注册回调函数
        call.enqueue(new Callback() {
            @Override
            public void onFailure(Call call, IOException e) {
                //请求失败时，回到主线程，更新UI
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        //提示用户请求失败
                        Toast.makeText(computer_link.this, "请求失败，请检查网络", Toast.LENGTH_SHORT).show();
                    }
                });
            }

            @Override
            public void onResponse(Call call, Response response) throws IOException {
                //请求成功时，回到主线程，更新UI
                runOnUiThread(new Runnable() {
                    @Override
                    public void run() {
                        //获取服务器返回的数据，假设是一个JSON字符串
                        try {
                            String result = response.body().string();
                            //解析JSON字符串，判断登录是否成功，这里省略解析过程
                            //如果登录成功，跳转到主界面，这里省略跳转过程
                            //如果登录失败，提示用户失败原因，这里省略提示过程
                            Intent intent = new Intent();
                            intent.setClass(computer_link.this, transport.class);
                            startActivity(intent);
                        } catch (IOException e) {
                            throw new RuntimeException(e);
                        }
                    }
                });
            }
        });
    }
}
