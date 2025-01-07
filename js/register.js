document.addEventListener('DOMContentLoaded', function() {
    // 获取元素
    const registerButton = document.getElementById('register-button');
    const usernameInput = document.getElementById('username');
    const passwordInput = document.getElementById('password');

    // 检查输入是否为空
    function validateInput() {
        if (!usernameInput.value.trim()) {
            Swal.fire({
                title: '提示',
                text: '请输入用户名！',
                icon: 'warning',
                confirmButtonText: '确定',
                confirmButtonColor: '#3085d6'
            });
            return false;
        }
        if (!passwordInput.value.trim()) {
            Swal.fire({
                title: '提示',
                text: '请输入密码！',
                icon: 'warning',
                confirmButtonText: '确定',
                confirmButtonColor: '#3085d6'
            });
            return false;
        }
        return true;
    }

    // 注册按钮点击事件
    registerButton.addEventListener('click', async () => {
        if (!validateInput()) {
            return;
        }

        let username = usernameInput.value.trim();
        let password = passwordInput.value.trim();

        try {
            // 显示加载中动画
            Swal.fire({
                title: '请稍候',
                text: '正在处理您的注册请求...',
                allowOutsideClick: false,
                allowEscapeKey: false,
                allowEnterKey: false,
                showConfirmButton: false,
                didOpen: () => {
                    Swal.showLoading();
                }
            });

            let response = await fetch('http://10.190.246.22:8080/api/reg', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: `username=${encodeURIComponent(username)}&password=${encodeURIComponent(password)}`,
            });

            let data = await response.json();

            // 根据返回的 code 显示不同的提示
            if (data.code === "401") {
                Swal.fire({
                    title: '注册失败',
                    text: '该用户名已存在，请重新输入！',
                    icon: 'error',
                    confirmButtonText: '确定',
                    confirmButtonColor: '#3085d6'
                });
            } else if (data.code === "200") {
                await Swal.fire({
                    title: '注册成功',
                    text: '欢迎加入天理社区！',
                    icon: 'success',
                    confirmButtonText: '确定',
                    confirmButtonColor: '#3085d6'
                });
                // 注册成功后，清空输入框
                usernameInput.value = '';
                passwordInput.value = '';
            } else {
                Swal.fire({
                    title: '提示',
                    text: data.msg,
                    icon: 'info',
                    confirmButtonText: '确定',
                    confirmButtonColor: '#3085d6'
                });
            }
        } catch (error) {
            Swal.fire({
                title: '错误',
                text: '注册失败，请稍后重试！',
                icon: 'error',
                confirmButtonText: '确定',
                confirmButtonColor: '#3085d6'
            });
            console.error('Error:', error);
        }
    });
}); 