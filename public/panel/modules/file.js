layui.define(['jquery', 'layer', 'laypage'], function(exports) { //提示：模块也可以依赖其它模块，如：layui.define('layer', callback);
    var $ = layui.jquery,
        layer = layui.layer,
        laypage = layui.laypage;
    //外部接口
    var fm = {
            config:{'test':'test','thumb':{'nopic':'',width:100,height:100},icon_url:'ico/',btn_upload:true,btn_create:true}
            ,cache: {} //数据缓存
            ,index: layui.fm ? (layui.fm.index + 10000) : 0
            //设置全局项
            ,set: function(options){
                var that = this;
                that.config = $.extend({}, that.config, options);
                return that;
            }
            //事件监听
            ,on: function(events, callback){
                return layui.onevent.call(this, 'file', events, callback);
            }
            ,dirRoot:[{'path':'','name': '根目录'}]
            ,v:'1.0.1.2019.12.26'
        }
        //操作当前实例
        , thisFm = function() {
            var that = this,
                options = that.config,
                id = options.id || options.index;

            // console.log(id)
            if (id) {
                thisFm.that[id] = that; //记录当前实例对象
                thisFm.config[id] = options; //记录当前实例配置项
            }
            return {
                config: options,
                reload: function(options) {
                    that.reload.call(that, options);
                }
            }
        }
        //获取当前实例配置项
        ,getThisFmConfig = function(id){
            var config = thisFm.config[id];
            if(!config) hint.error('The ID option was not found in the fm instance');
            return config || null;
        }
        //构造器
        ,Class = function(options){
            var that = this;
            that.config = $.extend({}, that.config, fm.config, options);
            //记录所有实例
            thisFm.that = {}; //记录所有实例对象
            thisFm.config = {}; //记录所有实例配置项
            // console.log(that.config)
            that.render();
        };
    //渲染
    Class.prototype.render = function(){
        var that = this
            ,options = that.config;

        options.elem = $(options.elem);
        options.where = options.where || {};
        options.id = options.id || options.elem.attr('id') || that.index;

        //请求参数的自定义格式
        options.request = $.extend({
            pageName: 'page'
            ,limitName: 'limit'
        }, options.request)

        //响应数据的自定义格式
        options.response = $.extend({
            statusName: 'code'
            ,statusCode: 0
            ,msgName: 'msg'
            ,dataName: 'data'
            ,countName: 'count'
        }, options.response);

        //如果 page 传入 laypage 对象
        if(typeof options.page === 'object'){
            options.limit = options.page.limit || options.limit;
            options.limits = options.page.limits || options.limits;
            that.page = options.page.curr = options.page.curr || 1;
            delete options.page.elem;
            delete options.page.jump;
        }

        if(!options.elem[0]) return that;
        //渲染主体
        var _btn = ''
        if(options.btn_create){
            _btn +='<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="new_dir">建文件夹</button>';
        }
        if(options.btn_upload){
            _btn +='<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="uploadfile">上传文件</button>';
        }
        var _html = '<div class="layui-card" >' +
            '<div class="layui-card-body">' +
            '<div class="layui-btn-group tool_bar">' +
            _btn+
            '<button type="button" class="layui-btn layui-btn-primary layui-btn-sm" id="back"><i class="layui-icon layui-icon-left line"></i></button>' +
            '</div>' +
            '<div class="layui-inline path_bar" id="">' +
            '<a ><i class="layui-icon layui-icon-more-vertical line" ></i>根目录</a>' +
            '</div>' +
            '</div><hr><div class="layui-card-body">' +
            '<div class="file-body layui-form" style="">' +
            '<ul class="file layui-row fm_body layui-col-space10" >' +
            '</ul>' +
            '</div>' +
            '<hr><div ><div class="layui_page_'+options.id+'" id="layui_page_'+options.id+'"></div></div></div>';

        options.elem.html(_html);

        options.index = that.index;
        that.key = options.id || options.index;
        //各级容器
        that.layPage = options.elem.find('.layui_page_'+options.id);
        that.layBody = options.elem.find('.fm_body');
        that.layPathBar = options.elem.find('.path_bar');
        that.layToolBar = options.elem.find('.tool_bar');
        that.pullData(that.page); //请求数据
        that.events(); //事件
    }

    //页码
    Class.prototype.page = 1;

    //获得数据
    Class.prototype.pullData = function(curr) {
        var that = this,
            options = that.config,
            request = options.request,
            response = options.response,
            _status = false;

        that.startTime = new Date().getTime(); //渲染开始时间
        if (options.url) { //Ajax请求
            var params = {};
            params[request.pageName] = curr;
            params[request.limitName] = options.limit;

            //参数
            var data = $.extend(params, options.where);
            if (options.contentType && options.contentType.indexOf("application/json") == 0) { //提交 json 格式
                data = JSON.stringify(data);
            }

            that.loading();

            $.ajax({
                type: options.method || 'get',
                url: options.url,
                contentType: options.contentType,
                data: data,
                async: false,
                dataType: 'json',
                headers: options.headers || {},
                success: function(res) {
                    //如果有数据解析的回调，则获得其返回的数据
                    if (typeof options.parseData === 'function') {
                        res = options.parseData(res) || res;
                    }
                    //检查数据格式是否符合规范
                    if (res[response.statusName] != response.statusCode) {

                        that.errorView(
                            res[response.msgName] ||
                            ('返回的数据不符合规范，正确的成功状态码应为："' + response.statusName + '": ' + response.statusCode)
                        );
                    } else {
                        // console.log(res, curr, res[response.countName]);
                        that.renderData(res, curr, res[response.countName]);

                        options.time = (new Date().getTime() - that.startTime) + ' ms'; //耗时（接口请求+视图渲染）
                    }
                    typeof options.done === 'function' && options.done(res, curr, res[response.countName]);
                    _status = true;
                },
                error: function(e, m) {
                    that.errorView('数据接口请求异常：' + m);

                }
            });
        }
        return _status;
    };
    //数据渲染
    Class.prototype.renderData = function(res, curr, count){
        var that = this
            ,options = that.config
            ,data = res[options.response.dataName] || []

        //渲染数据
        var _content = ''
        layui.each(data,function(i,v){
            let _img,_type;
            _type = v.type;
            switch (v.type) {
                case 'directory':
                    _img = '<div  style="width:'+options.thumb['width']+'px;height:'+options.thumb['height']+'px;line-height:'+options.thumb['height']+'px"><img src="ico/dir.png" style="vertical-align:middle;"></div>';
                    _type = 'DIR';
                    break;
                default:

                    if (v.type == 'png' || v.type == 'gif' || v.type == 'jpg' || v.type == 'image') {
                        _img = '<img src="' + v.thumb + '" width="'+options.thumb['width']+'" height="'+options.thumb['height']+'" onerror=\'this.src="'+options.thumb['nopic']+'"\'  />';
                    } else {
                        _img = '<div  style="width:'+options.thumb['width']+'px;height:'+options.thumb['height']+'px;line-height:'+options.thumb['height']+'px"><img src="' + options.icon_url + v.type + '.png"  onerror=\'this.src="'+options.thumb['nopic']+'"\' /></div>';
                    }
                    break;
            }
            _content+='<li style="display:inline-block" data-type="'+_type+'" data-index="'+i+'">' +
                '<div class="content" align="center">'+
                _img +
                '<p class="layui-elip" title="' + v.name + '">' + v.name + ' </p>' +
                '</div>' +
                '</li>';
        });
        options.elem.find('.file').html(_content);
        fm.cache[options.id] = data; //记录数据
        //显示隐藏分页栏
        // console.log(that.layPage)
        that.layPage[(count == 0 || (data.length === 0 && curr == 1)) ? 'addClass' : 'removeClass']('layui-hide');
        if(data.length === 0){
            return that.errorView('空目录');
        } else {
            //that.layFixed.removeClass('layui-hide');
        }
        //同步分页状态
        if(options.page){
            // console.log(options,'layui_page_' + options.id)
            options.page = $.extend({
                elem: 'layui_page_' + options.id
                ,count: count
                ,limit: options.limit
                ,limits: options.limits || [10,20,30,40,50,60,70,80,90]
                ,groups: 3
                ,layout: ['prev', 'page', 'next', 'skip', 'count', 'limit']
                ,prev: '<i class="layui-icon">&#xe603;</i>'
                ,next: '<i class="layui-icon">&#xe602;</i>'
                ,jump: function(obj, first){
                    if(!first){
                        //分页本身并非需要做以下更新，下面参数的同步，主要是因为其它处理统一用到了它们
                        //而并非用的是 options.page 中的参数（以确保分页未开启的情况仍能正常使用）
                        that.page = obj.curr; //更新页码
                        options.limit = obj.limit; //更新每页条数

                        that.pullData(obj.curr);
                    }
                }
            }, options.page);
            options.page.count = count; //更新总条数
            laypage.render(options.page);
        }
    };
    //更新路径工具条
    Class.prototype.updatePathBar = function(){
        // console.log('updatePathBar',fm.dirRoot);
        var that = this
            ,options = that.config;
        //请求数据
        let dir_cur = fm.dirRoot[fm.dirRoot.length -1];
        options.where = {'path':dir_cur['path']}
        let _rs = that.pullData(1);
        // console.log(_rs)
        if(false == _rs) return;
        that.layPathBar.html('');

        fm.dirRoot.map(function(item,index,arr){
            let icon = index==0 ?'layui-icon-more-vertical':'layui-icon-right';
            let html = '<i class="layui-icon '+icon+'"></i>'+
                '<a  data-path="' + item.path + '" data-name="' + item.name + '" >' + item.name + '</a>'
            that.layPathBar.append(html);
        })


    }
    //事件处理
    Class.prototype.events = function(){
        var that = this
            ,options = that.config
            ,_BODY = $('body')
            ,dict = {}
            ,filter = options.elem.attr('lay-filter');
        //文件事件
        that.layBody.on('click', 'li', function(){ //单击行
            setPicEvent.call(this, 'pic');
        });
        //文件夹事件
        that.layBody.on('click', 'li[data-type=DIR]', function(){ //单击行
            var othis = $(this);
            var data =  fm.cache[options.id];
            var index = othis.data('index');
            data = data[index] || {};

            //导航图标
            fm.dirRoot.push({'path':data.path,'name': data.name});
            that.updatePathBar();
        });
        //返回上一级目录
        that.layToolBar.on('click', '#back', function(){
            var othis = $(this);
            if(fm.dirRoot.length == 1) return layer.msg('已经是根目录');

            fm.dirRoot.length >1 && fm.dirRoot.pop()
            that.updatePathBar();

            // console.log('back');
        });
        //上传文件
        that.layToolBar.on('click', '#uploadfile', function(){
            var othis = $(this);
            let eventType = 'uploadfile';
            layui.event.call(this,
                'file', eventType + '('+ filter +')'
                ,{obj:othis,path:fm.dirRoot[fm.dirRoot.length -1]['path']}
            );
            // console.log('uploadfile');
        });
        //新建文件夹
        that.layToolBar.on('click', '#new_dir', function(){
            var othis = $(this);
            let eventType = 'new_dir';
            layer.prompt({ title: '请输入新文件夹名字', formType: 0 }, function(name, index) {
                layer.close(index);
                //新建文件夹
                layui.event.call(this,
                    'file', eventType + '('+ filter +')'
                    ,{obj:othis,folder:name,path:fm.dirRoot[fm.dirRoot.length -1]['path']}
                );
            });
        });
        //创建点击文件事件监听
        var setPicEvent = function(eventType) {
            var othis = $(this);
            var data =  fm.cache[options.id];
            var index = othis.data('index');
            if (othis.data('type')=='DIR') return; //不触发事件
            data = data[index] || {};
            layui.event.call(this,
                'file', eventType + '('+ filter +')'
                ,{obj:othis,data:data}
            );
        };
    };

    //请求loading
    Class.prototype.loading = function(hide){
        var that = this
            ,options = that.config;
        if(options.loading){
            if(hide){
                that.layInit && that.layInit.remove();
                delete that.layInit;
                that.layBox.find(ELEM_INIT).remove();
            } else {
                that.layInit = $(['<div class="layui-table-init">'
                    ,'<i class="layui-icon layui-icon-loading layui-anim layui-anim-rotate layui-anim-loop"></i>'
                    ,'</div>'].join(''));
                that.layBox.append(that.layInit);
            }
        }
    };
    //异常提示
    Class.prototype.errorView = function(html){
        var that = this
        layer.msg(html);

    };
    //重载
    Class.prototype.reload = function(options){
        var that = this;

        options = options || {};
        delete that.haveInit;

        if(options.data && options.data.constructor === Array) delete that.config.data;
        that.config = $.extend(true, {}, that.config, options);

        that.render();
    };
    //重载
    fm.reload = function(id, options){
        var config = getThisFmConfig(id); //获取当前实例配置项
        if(!config) return;

        var that = thisFm.that[id];
        that.reload(options);

        return thisFm.call(that);
    };
    //核心入口
    fm.render = function(options){
        var inst = new Class(options);
        return thisFm.call(inst);
    };
    exports('file', fm);
});
