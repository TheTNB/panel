<!--
Name: 网站 - 全局设置
Author: 耗子
Date: 2022-11-30
-->
<script type="text/html" template lay-url="/api/panel/website/getDefaultSettings" lay-done="layui.data.sendParams(d.params)">
    <div class="layui-tab">
        <ul class="layui-tab-title">
            <li class="layui-this">默认页</li>
            <li>停止页</li>
        </ul>
        <div class="layui-tab-content">
            <div class="layui-tab-item layui-show">
                <!-- 默认页 -->
                <blockquote class="layui-elem-quote layui-quote-nm">
                    设置站点未找到时的提示页面。
                </blockquote>
                <div id="index-editor" style="height: -webkit-fill-available;">@{{ d.data.index }}</div>
            </div>
            <div class="layui-tab-item">
                <!-- 停止页 -->
                <blockquote class="layui-elem-quote layui-quote-nm">
                    设置站点停止时的提示页面。
                </blockquote>
                <div id="stop-editor" style="height: -webkit-fill-available;">@{{ d.data.stop }}</div>
            </div>
        </div>
    </div>
    <div class="layui-footer">
        <button id="save-website-default-settings" class="layui-btn">保存设置</button>
    </div>
</script>
<script>
    let indexEditor = '';
    let stopEditor = '';
    layui.data.sendParams = function(params) {
        layui.use(['admin', 'form', 'laydate', 'code'], function () {
            var $ = layui.$
                , admin = layui.admin
                , element = layui.element
                , layer = layui.layer
                , laydate = layui.laydate
                , code = layui.code;
            indexEditor = ace.edit("index-editor", {
                mode: "ace/mode/html",
                selectionStyle: "text"
            });
            stopEditor = ace.edit("stop-editor", {
                mode: "ace/mode/html",
                selectionStyle: "text"
            });

            $('#save-website-default-settings').click(function () {
                layer.load(2);
                admin.req({
                    url: '/api/panel/website/saveDefaultSettings'
                    , type: 'post'
                    , data: {
                        index: indexEditor.getValue(),
                        stop: stopEditor.getValue()
                    }
                    , success: function (res) {
                        layer.closeAll('loading');
                        if (res.code === 0) {
                            layer.msg('保存成功', {icon: 1});
                            setTimeout(function () {
                                admin.render();
                            }, 1000);
                        } else {
                            layer.msg(res.msg, {icon: 2});
                        }
                    }
                    , error: function (xhr, status, error) {
                        console.log('耗子Linux面板：ajax请求出错，错误' + error);
                    }
                });
            });

        });
    };
</script>
