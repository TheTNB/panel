/**
 * forum demo
 */


layui.define(['table', 'form'], function(exports){
  var $ = layui.$
  ,admin = layui.admin
  ,view = layui.view
  ,table = layui.table
  ,form = layui.form;

  //帖子管理
  table.render({
    elem: '#LAY-app-forum-list'
    ,url: './res/json/forum/list.js' //模拟接口
    ,cols: [[
      {type: 'checkbox', fixed: 'left'}
      ,{field: 'id', width: 100, title: 'ID', sort: true}
      ,{field: 'poster', title: '发帖人'}
      ,{field: 'avatar', title: '头像', width: 100, templet: '#imgTpl'}
      ,{field: 'content', title: '发帖内容'}
      ,{field: 'posttime', title: '发帖时间', sort: true}
      ,{field: 'top', title: '置顶', templet: '#buttonTpl', minWidth: 80, align: 'center'}
      ,{title: '操作', width: 150, align: 'center', fixed: 'right', toolbar: '#table-forum-list'}
    ]]
    ,page: true
    ,limit: 10
    ,limits: [10, 15, 20, 25, 30]
    ,text: '对不起，加载出现异常！'
  });
  
  //工具条
  table.on('tool(LAY-app-forum-list)', function(obj){
    var data = obj.data;
    if(obj.event === 'del'){
      layer.confirm('确定删除此条帖子？', function(index){
        obj.del();
        layer.close(index);
      });
    } else if(obj.event === 'edit'){
      admin.popup({
        title: '编辑帖子'
        ,area: ['550px', '450px']
        ,id: 'LAY-popup-forum-edit'
        ,resize: false
        ,success: function(layero, index){
          view(this.id).render('app/forum/listform', data).done(function(){
            form.render(null, 'layuiadmin-form-list');
            
            //提交
            form.on('submit(layuiadmin-app-forum-submit)', function(data){
              var field = data.field; //获取提交的字段

              //提交 Ajax 成功后，关闭当前弹层并重载表格
              //$.ajax({});
              layui.table.reload('LAY-app-forum-list'); //重载表格
              layer.close(index); //执行关闭 
            });
          });
        }
      });
    }
  });

  //回帖管理
  table.render({
    elem: '#LAY-app-forumreply-list'
    ,url: './res/json/forum/replys.js' //模拟接口
    ,cols: [[
      {type: 'checkbox', fixed: 'left'}
      ,{field: 'id', width: 100, title: 'ID', sort: true}
      ,{field: 'replyer', title: '回帖人'}
      ,{field: 'cardid', title: '回帖ID', sort: true}
      ,{field: 'avatar', title: '头像', width: 100, templet: '#imgTpl'}
      ,{field: 'content', title: '回帖内容', width: 200}
      ,{field: 'replytime', title: '回帖时间', sort: true}
      ,{title: '操作', width: 150, align: 'center', fixed: 'right', toolbar: '#table-forum-replys'}
    ]]
    ,page: true
    ,limit: 10
    ,limits: [10, 15, 20, 25, 30]
    ,text: '对不起，加载出现异常！'
  });
  
  //工具条
  table.on('tool(LAY-app-forumreply-list)', function(obj){
    var data = obj.data;
    if(obj.event === 'del'){
      layer.confirm('确定删除此条评论？', function(index){
        obj.del();
        layer.close(index);
      });
    } else if(obj.event === 'edit'){
      admin.popup({
        title: '编辑回帖'
        ,area: ['550px', '400px']
        ,id: 'LAY-popup-forum-edit'
        ,resize: false
        ,success: function(layero, index){
          view(this.id).render('app/forum/replysform', data).done(function(){
            form.render(null, 'layuiadmin-app-forum-reply');
            
            //提交
            form.on('submit(layuiadmin-app-forumreply-submit)', function(data){
              var field = data.field; //获取提交的字段

              //提交 Ajax 成功后，关闭当前弹层并重载表格
              //$.ajax({});
              layui.table.reload('LAY-app-forumreply-list'); //重载表格
              layer.close(index); //执行关闭 
            });
          });
        }
      });
    }
  });
  
  exports('forum', {})
});