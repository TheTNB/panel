<!--
Name: 面板设置模版
Author: 耗子
Date: 2022-12-08
-->
<title>面板设置</title>
<div class="layui-fluid">
    <div class="layui-row layui-col-space15">
        <div class="layui-col-md12">
            <div class="layui-card">
                <div class="layui-card-header">设置</div>
                <div class="layui-card-body">

                    <div class="layui-form" lay-filter="panel_setting">
                        <div class="layui-form-item">
                            <label class="layui-form-label">API 开关</label>
                            <div class="layui-input-inline">
                                <input type="checkbox" name="api" lay-skin="switch" lay-text="ON|OFF"/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">开启后将提供面板API接口的访问支持</div>
                        </div>
                        <div id="setting-api-token" class="layui-form-item">
                            <label class="layui-form-label">API Token</label>
                            <div class="layui-input-inline">
                                <input type="text" name="api_token" value="获取中ing..." class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">API Token，用于携带访问面板接口</div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">多设备登录</label>
                            <div class="layui-input-inline">
                                <input type="checkbox" name="multi_login" lay-skin="switch" lay-text="ON|OFF"/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">
                                开启后将允许多设备同时登录面板，可能具有一定安全隐患
                            </div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">面板名称</label>
                            <div class="layui-input-inline">
                                <input type="text" name="name" value="获取中ing..." class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">修改面板的显示名称</div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">面板用户名</label>
                            <div class="layui-input-inline">
                                <input type="text" name="username" value="获取中ing..." class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">修改面板的登录用户名</div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">面板密码</label>
                            <div class="layui-input-inline">
                                <input type="password" name="password" value="" class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">修改面板的登录密码（留空不修改）</div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">面板邮箱</label>
                            <div class="layui-input-inline">
                                <input type="text" name="email" value="获取中ing..." class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">修改面板账号的邮箱，目前用于签发免费SSL证书</div>
                        </div>
                        <div class="layui-form-item">
                            <label class="layui-form-label">面板端口</label>
                            <div class="layui-input-inline">
                                <input type="text" name="port" value="" class="layui-input" disabled/>
                            </div>
                            <div class="layui-form-mid layui-word-aux">修改面板的访问端口（<b style="color: red;">保存后需要手动修改浏览器地址栏的端口为新端口以访问面板</b>）
                            </div>
                        </div>
                        <div class="layui-form-item">
                            <div class="layui-input-block">
                                <button class="layui-btn" lay-submit lay-filter="panel_setting_submit">确认修改</button>
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </div>
</div>

<script>
    layui.define(['form', 'upload'], function () {
        var $ = layui.$
            , layer = layui.layer
            , admin = layui.admin
            , form = layui.form;

        // 渲染表单
        form.render();
        $('#setting-api-token').hide();

        // ajax获取设置项并赋值
        admin.req({
            url: "/api/panel/setting/get"
            , method: 'get'
            , success: function (result) {
                if (result.code !== 0) {
                    console.log('耗子Linux面板：系统信息获取失败，接口返回' + result);
                    layer.msg('系统信息获取失败，请刷新重试！')
                    return false;
                }
                form.val("panel_setting",
                    result.data
                );
                $('input').attr('disabled', false);
                if (result.data.api === 1) {
                    $('#setting-api-token').show();
                    $('#setting-api-token input').attr('readonly', true);
                }
            }
            , error: function (xhr, status, error) {
                console.log('耗子Linux面板：ajax请求出错，错误' + error);
            }
        });

        // 面板设置
        form.on('submit(panel_setting_submit)', function (obj) {
            // 面板API
            if (obj.field.api === "on") {
                obj.field.api = 1;
            } else {
                obj.field.api = 0;
            }
            // 多设备登录
            if (obj.field.multi_login === "on") {
                obj.field.multi_login = 1;
            } else {
                obj.field.multi_login = 0;
            }
            // 提交修改
            admin.req({
                url: "/api/panel/setting/save"
                , method: 'post'
                , data: obj.field
                , success: function (result) {
                    if (result.code !== 0) {
                        console.log('耗子Linux面板：设置保存失败，接口返回' + result);
                        layer.msg('面板设置保存失败，请刷新重试！')
                        return false;
                    }
                    admin.render();
                    layer.msg('面板设置保存成功！');
                }
                , error: function (xhr, status, error) {
                    console.log('耗子Linux面板：ajax请求出错，错误' + error);
                }
            });
            return false;
        });
    });
</script>