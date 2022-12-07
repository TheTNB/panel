<!--
Name: 网站 - 编辑
Author: 耗子
Date: 2022-12-08
-->
<script type="text/html" template lay-done="layui.data.sendParams(d.params)">
    <div class="layui-tab" lay-filter="website-edit-tab">
        <ul class="layui-tab-title">
            <li class="layui-this">域名端口</li>
            <li>基本设置</li>
            <li>防火墙</li>
            <li>SSL</li>
            <li>伪静态</li>
            <li>配置原文</li>
            <li>访问日志</li>
        </ul>
        <div class="layui-tab-content">
            <div class="layui-tab-item layui-show">
                <!-- 域名绑定 -->
                <div class="layui-form layui-form-pane">
                    <div class="layui-form-item layui-form-text">
                        <label class="layui-form-label">域名</label>
                        <div class="layui-input-block">
                            <textarea name="domain" lay-verify="required"
                                      placeholder="请输入域名，一行一个支持泛域名"
                                      class="layui-textarea">@{{ d.params.config.domain }}</textarea>
                        </div>
                    </div>
                    <hr>
                    <div class="layui-form-item layui-form-text">
                        <label class="layui-form-label">端口</label>
                        <div class="layui-input-block">
                            <textarea name="port" lay-verify="required"
                                      placeholder="请输入访问端口，一行一个"
                                      class="layui-textarea">@{{ d.params.config.port }}</textarea>
                        </div>
                    </div>
                </div>
            </div>
            <div class="layui-tab-item">
                <!-- 基本设置 -->
                <div class="layui-form layui-form-pane">
                    <div class="layui-form-item">
                        <label class="layui-form-label">网站目录</label>
                        <div class="layui-input-block">
                            <input type="text" name="path" autocomplete="off" placeholder="请输入网站目录"
                                   class="layui-input" value="@{{ d.params.config.path }}">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">运行目录</label>
                        <div class="layui-input-block">
                            <input type="text" name="root" autocomplete="off" placeholder="请输入网站运行目录（Laravel等程序需要）"
                                   class="layui-input" value="@{{ d.params.config.root }}">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">默认文档</label>
                        <div class="layui-input-block">
                            <input type="text" name="index" autocomplete="off" placeholder="请输入默认文档，以空格隔开"
                                   class="layui-input" value="@{{ d.params.config.index }}">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">PHP版本</label>
                        <div class="layui-input-block">
                            <select name="php" lay-filter="website-php">
                                @{{# layui.each(d.params.php_version, function(index, item){ }}
                                @{{# if(item == d.params.config.php){ }}
                                <option value="@{{ item }}" selected="">@{{ item }}</option>
                                @{{# }else{ }}
                                <option value="@{{ item }}">@{{ item }}</option>
                                @{{# } }}
                                @{{# }); }}
                            </select>
                        </div>
                    </div>
                    <div class="layui-form-item" pane="">
                        <label class="layui-form-label">防跨站攻击</label>
                        <div class="layui-input-block">
                            <input type="checkbox" name="open_basedir" lay-skin="switch" lay-text="ON|OFF"
                                   @{{ d.params.config.open_basedir== 1 ? 'checked' : '' }} />
                        </div>
                    </div>
                </div>
            </div>
            <div class="layui-tab-item">
                <!-- 防火墙 -->
                <blockquote class="layui-elem-quote layui-quote-nm">
                    面板自带开源的 ngx_waf 防火墙<br>文档参考：<a
                            href="https://docs.addesp.com/ngx_waf/zh-cn/advance/directive.html"
                            target="_blank">https://docs.addesp.com/ngx_waf/zh-cn/advance/directive.html</a>
                </blockquote>
                <div class="layui-form layui-form-pane">
                    <div class="layui-form-item" pane="">
                        <label class="layui-form-label">总开关</label>
                        <div class="layui-input-inline">
                            <input type="checkbox" name="waf" lay-skin="switch" lay-text="ON|OFF"
                                   @{{ d.params.config.waf== 1 ? 'checked' : '' }} />
                        </div>
                        <div class="layui-form-mid layui-word-aux">只有打开了总开关，下面的设置才会生效！</div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">模式</label>
                        <div class="layui-input-block">
                            <input type="text" name="waf_mode" autocomplete="off" placeholder="DYNAMIC"
                                   class="layui-input" value="@{{ d.params.config.waf_mode }}">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">CC</label>
                        <div class="layui-input-block">
                            <input type="text" name="waf_cc_deny" autocomplete="off"
                                   placeholder="rate=1000r/m duration=60m"
                                   class="layui-input" value="@{{ d.params.config.waf_cc_deny }}">
                        </div>
                    </div>
                    <div class="layui-form-item">
                        <label class="layui-form-label">缓存</label>
                        <div class="layui-input-block">
                            <input type="text" name="waf_cache" autocomplete="off" placeholder="capacity=50"
                                   class="layui-input" value="@{{ d.params.config.waf_cache }}">
                        </div>
                    </div>
                </div>
            </div>
            <div class="layui-tab-item">
                <!-- SSL -->
                <div class="layui-form layui-form-pane">
                    <div class="layui-form-item" pane="">
                        <label class="layui-form-label">总开关</label>
                        <div class="layui-input-inline">
                            <input type="checkbox" name="ssl" lay-skin="switch" lay-text="ON|OFF"
                                   @{{ d.params.config.ssl== 1 ? 'checked' : '' }} />
                        </div>
                        <div class="layui-form-mid layui-word-aux">只有打开了总开关，下面的设置才会生效！</div>
                    </div>

                    <div class="layui-form-item">
                        <div class="layui-inline">
                            <label class="layui-form-label">HTTP跳转</label>
                            <div class="layui-input-block">
                                <input type="checkbox" name="http_redirect" lay-skin="switch" lay-text="ON|OFF"
                                       @{{ d.params.config.http_redirect== 1 ? 'checked' : '' }} />
                            </div>
                        </div>
                        <div class="layui-inline">
                            <label class="layui-form-label">HSTS</label>
                            <div class="layui-input-inline">
                                <input type="checkbox" name="hsts" lay-skin="switch" lay-text="ON|OFF"
                                       @{{ d.params.config.hsts== 1 ? 'checked' : '' }} />
                            </div>
                        </div>
                        <div class="layui-inline">
                            <div class="layui-input-inline">
                                <button id="issue-ssl" class="layui-btn layui-btn-sm">签发免费SSL证书</button>
                            </div>
                        </div>
                    </div>

                    <div class="layui-form-item layui-form-text">
                        @{{# if(d.params.config.ssl == 1){ }}
                        <label class="layui-form-label">证书 <span style="color: red; float: right;">剩余有效期：@{{ d.params.config.ssl_date }}天</span></label>
                        @{{# }else{ }}
                        <label class="layui-form-label">证书</label>
                        @{{# } }}
                        <div class="layui-input-block">
                            <textarea name="ssl_certificate" placeholder="请输入pem证书文件的内容"
                                      class="layui-textarea">@{{ d.params.config.ssl_certificate }}</textarea>
                        </div>
                    </div>
                    <div class="layui-form-item layui-form-text">
                        <label class="layui-form-label">私钥</label>
                        <div class="layui-input-block">
                            <textarea name="ssl_certificate_key" placeholder="请输入key私钥文件的内容"
                                      class="layui-textarea">@{{ d.params.config.ssl_certificate_key }}</textarea>
                        </div>
                    </div>

                </div>
            </div>
            <div class="layui-tab-item">
                <!-- 伪静态 -->
                <blockquote class="layui-elem-quote layui-quote-nm">
                    设置伪静态规则，填入 <code>location</code> 部分即可
                </blockquote>
                <div id="rewrite-editor" style="height: 400px;">@{{ d.params.config.rewrite }}</div>
            </div>
            <div class="layui-tab-item">
                <!-- 配置原文 -->
                <blockquote class="layui-elem-quote layui-quote-nm">
                    如果您不了解配置规则，请勿随意修改，否则可能会导致网站无法访问或面板功能异常！如果已经遇到问题，可尝试：
                    <button id="site-config-restore" class="layui-btn layui-btn-xs">重置配置</button>
                    <br>
                    如果你修改了原文，那么点击保存后，其余的修改将不会生效！
                </blockquote>
                <div id="config-editor" style="height: 400px;">@{{ d.params.config.config_raw }}</div>
            </div>
            <div class="layui-tab-item">
                <!-- 访问日志 -->
                <button id="clean-site-log" class="layui-btn">清空日志</button>
                <pre id="website-log" class="layui-code" lay-options="{about: '@{{ d.params.config.name }}.log'}">@{{ d.params.config.log }}</pre>
            </div>
        </div>
    </div>
    <div class="layui-footer">
        <button id="save-site-config" class="layui-btn">保存设置</button>
    </div>
</script>
<script>
    let rewriteEditor = '';
    let configEditor = '';
    layui.data.sendParams = function (params) {
        layui.use(['admin', 'form', 'laydate', 'code'], function () {
            var $ = layui.$
                , admin = layui.admin
                , layer = layui.layer
                , code = layui.code
                , form = layui.form
                , element = layui.element;
            form.render();
            element.render();
            element.on('tab(website-edit-tab)', function (data) {
                if (data.index === 6) {
                    // 隐藏保存按钮
                    $('.layui-footer').hide();
                } else {
                    $('.layui-footer').show();
                }
            });
            rewriteEditor = ace.edit("rewrite-editor", {
                mode: "ace/mode/nginx",
                selectionStyle: "text"
            });
            configEditor = ace.edit("config-editor", {
                mode: "ace/mode/nginx",
                selectionStyle: "text"
            });
            code({
                elem: '#website-log'
                , encode: true
                , about: false

            });

            $("#clean-site-log").click(function () {
                layer.confirm('确定要清空日志吗？', function (index) {
                    layer.close(index);
                    layer.load();
                    admin.req({
                        url: '/api/panel/website/clearSiteLog'
                        , type: 'post'
                        , data: {name: params.config.name}
                        , success: function (res) {
                            layer.closeAll('loading');
                            if (res.code === 0) {
                                layer.msg('已清空', {icon: 1});
                                setTimeout(function () {
                                    admin.render();
                                }, 1000);
                            } else {
                                layer.msg(res.msg, {icon: 2});
                            }
                        }
                        , error: function (xhr, status, error) {
                            layer.closeAll('loading');
                            console.log('耗子Linux面板：ajax请求出错，错误' + error);
                        }
                    });
                });
            });

            $('#save-site-config').click(function () {
                layer.load();
                var port = $('textarea[name="port"]').val();
                var reg = new RegExp(/\n443.*\n?/);
                // 如果开启了https，就自动添加443端口
                if ($('input[name="ssl"]').prop('checked') && !reg.test(port)) {
                    port = port + '\n443';
                }
                // 如果关闭了https，就自动删除443端口
                if (!$('input[name="ssl"]').prop('checked') && reg.test(port)) {
                    // 正则替换
                    port = port.replace(/443.*\n?/, '');
                }
                admin.req({
                    url: '/api/panel/website/saveSiteSettings'
                    , type: 'post'
                    , data: {
                        name: params.config.name,
                        config: {
                            domain: $('textarea[name="domain"]').val(),
                            port: port,
                            ssl: $('input[name="ssl"]').prop('checked') ? 1 : 0,
                            http_redirect: $('input[name="http_redirect"]').prop('checked') ? 1 : 0,
                            hsts: $('input[name="hsts"]').prop('checked') ? 1 : 0,
                            ssl_certificate: $('textarea[name="ssl_certificate"]').val(),
                            ssl_certificate_key: $('textarea[name="ssl_certificate_key"]').val(),
                            path: $('input[name="path"]').val(),
                            root: $('input[name="root"]').val(),
                            index: $('input[name="index"]').val(),
                            php: $('select[name="php"]').val(),
                            open_basedir: $('input[name="open_basedir"]').prop('checked') ? 1 : 0,
                            waf: $('input[name="waf"]').prop('checked') ? 1 : 0,
                            waf_mode: $('input[name="waf_mode"]').val(),
                            waf_cc_deny: $('input[name="waf_cc_deny"]').val(),
                            waf_cache: $('input[name="waf_cache"]').val(),
                            rewrite: rewriteEditor.getValue(),
                            config_raw: configEditor.getValue()
                        }
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

            // 重置配置
            $('#site-config-restore').click(function () {
                layer.confirm('高风险操作，网站配置重置后所有配置均需重新设置，确定要重置配置吗？', function (index) {
                    index = layer.msg('重置网站配置', {
                        icon: 16
                        , time: 0
                    });
                    admin.req({
                        url: '/api/panel/website/resetSiteConfig'
                        , type: 'post'
                        , data: {name: params.config.name}
                        , success: function (res) {
                            layer.close(index);
                            if (res.code === 0) {
                                layer.alert('重置成功，你需要重新添加域名/端口绑定，设置各配置参数！', function (index) {
                                    admin.render();
                                    layer.close(index);
                                });
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

            // 监听签发证书按钮
            $('#issue-ssl').click(function () {
                layer.confirm('确定要申请签发免费SSL证书吗？', function (index) {
                    index = layer.msg('正在签发证书，可能需要较长时间，请勿操作...', {
                        icon: 16
                        , time: 0
                    });
                    admin.req({
                        url: '/api/panel/website/issueSsl'
                        , type: 'post'
                        , data: {
                            name: params.config.name
                            , type: 'lets'
                        }
                        , success: function (res) {
                            layer.close(index);
                            if (res.code === 0) {
                                layer.msg('签发成功', {icon: 1});
                                setTimeout(function () {
                                    admin.render();
                                }, 1000);
                            } else {
                                layer.alert(res.msg, {icon: 2});
                            }
                        }
                        , error: function (xhr, status, error) {
                            layer.closeAll('loading');
                            console.log('耗子Linux面板：ajax请求出错，错误' + error);
                        }
                    });
                });
            });
        });
    };
</script>