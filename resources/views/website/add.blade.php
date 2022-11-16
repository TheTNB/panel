<!--
Name: 网站 - 添加
Author: 耗子
Date: 2022-10-14
-->
<script type="text/html" template lay-done="layui.data.sendParams(d.params)">
    <form class="layui-form" action="" lay-filter="add-website-form">
        <div class="layui-form-item">
            <label class="layui-form-label">网站名</label>
            <div class="layui-input-block">
                <input type="text" name="name" lay-verify="required" placeholder="请输入网站名（英文，设置后不可修改）"
                       autocomplete="off" class="layui-input">
            </div>
        </div>
        <div class="layui-form-item layui-form-text">
            <label class="layui-form-label">域名</label>
            <div class="layui-input-block">
                <textarea name="domain" lay-verify="required"
                          placeholder="请输入域名，一行一个支持泛域名（格式：yourdomain.com:88 ，端口不填则默认80端口）"
                          class="layui-textarea"></textarea>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">PHP版本</label>
                <div class="layui-input-block">
                    <select name="php" lay-filter="add-website-php">
                        @{{# layui.each(d.params.php_version, function(index, item){ }}
                        @{{# if(index == "00"){ }}
                        <option value="@{{ item }}" selected="">@{{ item }}</option>
                        @{{# }else{ }}
                        <option value="@{{ item }}">@{{ item }}</option>
                        @{{# } }}
                        @{{# }); }}
                    </select>
                </div>
            </div>
            <div class="layui-inline">
                <label class="layui-form-label">数据库</label>
                <div class="layui-input-block">
                    <select name="db_type" lay-filter="add-website-db">
                        <option value="" selected=""></option>
                        @{{# layui.each(d.params.db_version, function(index, item){ }}
                        @{{# if(item){ }}
                        <option value="@{{ index }}">@{{ index }}</option>
                        @{{# } }}
                        @{{# }); }}
                    </select>
                </div>
            </div>
        </div>
        <div id="add-website-db-info" class="layui-form-item">
            <div class="layui-inline">
                <label class="layui-form-label">数据库名</label>
                <div class="layui-input-inline">
                    <input type="text" name="db_name" lay-verify="required" autocomplete="off" class="layui-input">
                </div>
            </div>
            <div class="layui-inline">
                <label class="layui-form-label">数据库用户</label>
                <div class="layui-input-inline">
                    <input type="text" name="db_username" lay-verify="required" autocomplete="off" class="layui-input">
                </div>
            </div>
            <div class="layui-inline">
                <label class="layui-form-label">数据库密码</label>
                <div class="layui-input-inline">
                    <input type="text" name="db_password" lay-verify="required" autocomplete="off" class="layui-input">
                </div>
            </div>
        </div>

        <div class="layui-form-item">
            <label class="layui-form-label">目录</label>
            <div class="layui-input-block">
                <input type="text" name="path"
                       placeholder="请输入网站根目录（不填默认为/www/wwwroot/网站名）"
                       autocomplete="off" class="layui-input">
            </div>
        </div>
        <div class="layui-form-item layui-form-text">
            <label class="layui-form-label">备注</label>
            <div class="layui-input-block">
                <textarea name="note" placeholder="请输入备注内容，可以为空。" class="layui-textarea"></textarea>
            </div>
        </div>
        <div class="layui-form-item">
            <div class="layui-input-block">
                <div class="layui-footer">
                    <button class="layui-btn" lay-submit="" lay-filter="add-website-submit">立即提交</button>
                    <button type="reset" class="layui-btn layui-btn-primary">重置</button>
                </div>
            </div>
        </div>
    </form>
</script>
<script>
    layui.data.sendParams = function (params) {
        layui.use(['admin', 'form', 'laydate'], function () {
            var $ = layui.$
                , admin = layui.admin
                , element = layui.element
                , layer = layui.layer
                , laydate = layui.laydate
                , table = layui.table
                , form = layui.form;

            $("#add-website-db-info").hide();
            form.render();

            form.on('select(add-website-db)', function (data) {
                console.log(data.value); //得到被选中的值
                if (data.value === "") {
                    $("#add-website-db-info").hide();
                    return false;
                }
                if (data.value === 'mysql') {
                    $("#add-website-db-info").show();
                    $('input[name="db_name"]').val($('input[name="name"]').val() + '_mysql');
                    $('input[name="db_user"]').val($('input[name="name"]').val() + '_mysql');
                    $('input[name="db_password"]').val($('input[name="name"]').val() + '_password');
                }
            });
            // 提交
            form.on('submit(add-website-submit)', function (data) {
                // 判断db_type是否为空
                if (data.field.db_type === "") {
                    data.field.db = 0;
                }else{
                    data.field.db = 1;
                }
                console.log(data.field);
                admin.req({
                    url: "/api/panel/website/add"
                    , method: 'post'
                    , data: data.field
                    , success: function (result) {
                        if (result.code !== 0) {
                            console.log('耗子Linux面板：新网站添加失败，接口返回' + result);
                            layer.msg('新网站添加失败，请刷新重试！')
                            return false;
                        }
                        layer.alert('新网站添加成功！', {
                            icon: 1
                            , title: '提示'
                            , btn: ['确定']
                            , yes: function (index) {
                                layer.close(index);
                                location.reload();
                            }
                        });
                    }
                    , error: function (xhr, status, error) {
                        console.log('耗子Linux面板：ajax请求出错，错误' + error)
                    }
                });
                return false;
            });
        });
    };
</script>