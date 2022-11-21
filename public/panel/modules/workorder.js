/**
 * workorder demo
 */


layui.define(['table', 'form', 'element'], function(exports){
  var $ = layui.$
  ,admin = layui.admin
  ,view = layui.view
  ,table = layui.table
  ,form = layui.form
  ,element = layui.element;

  table.render({
    elem: '#LAY-app-workorder'
    ,url: './res/json/workorder/demo.js' //模拟接口
    ,cols: [[
      {type: 'numbers', fixed: 'left'}
      ,{field: 'orderid', width: 100, title: '工单号', sort: true}
      ,{field: 'attr', width: 100, title: '业务性质'}
      ,{field: 'title', width: 100, title: '工单标题', width: 300}
      ,{field: 'progress', title: '进度', width: 200, align: 'center', templet: '#progressTpl'}
      ,{field: 'submit', width: 100, title: '提交者'}
      ,{field: 'accept', width: 100, title: '受理人员'}
      ,{field: 'state', title: '工单状态', templet: '#buttonTpl', minWidth: 80, align: 'center'}
      ,{title: '操作', align: 'center', fixed: 'right', toolbar: '#table-system-order'}
    ]]
    ,page: true
    ,limit: 10
    ,limits: [10, 15, 20, 25, 30]
    ,text: '对不起，加载出现异常！'
    ,done: function(){
      element.render('progress');
    }
  });

  //工具条
  table.on('tool(LAY-app-workorder)', function(obj){
    var data = obj.data;
    if(obj.event === 'edit'){
      admin.popup({
        title: '编辑工单'
        ,area: ['450px', '450px']
        ,id: 'LAY-popup-workorder-add'
        ,success: function(layero, index){
          view(this.id).render('app/workorder/listform').done(function(){
            form.render(null, 'layuiadmin-form-workorder');
            
            //提交
            form.on('submit(LAY-app-workorder-submit)', function(data){
              var field = data.field; //获取提交的字段

              //提交 Ajax 成功后，关闭当前弹层并重载表格
              //$.ajax({});
              layui.table.reload('LAY-app-workorder'); //重载表格
              layer.close(index); //执行关闭 
            });
          });
        }
      });
    }
  });

  exports('workorder', {})
});