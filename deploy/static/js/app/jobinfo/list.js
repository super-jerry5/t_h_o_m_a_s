/**
 * Created by banxia on 16/5/17.
 */





var JobInfoList = {

    init:function(){
        JobInfoList.ctrl.initEvent();
    },
    ctrl:{

        initEvent:function(){

            $('.btn_operate_delete').on('click',function(){


                var va = $(this).attr("att");
                
                layer.confirm('Do you really want to delete this task', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
                    layer.close(index);
                    JobInfoList.ctrl.deleteJobInfo(va);
                });


            });
            
            $('.btn_operate_stop').on('click',function(){

                var va_arr = $(this).attr("att").split("|")
                var id = va_arr[0]
                var state = va_arr[1]
                if (state != "2" && state != "1")
                {
                	layer.confirm('Couldn\'t stop task not in state(STARTING, TRANSCODING) ', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
                        layer.close(index);
                    });
                }
                else
                {
                	layer.confirm('Do you really want to stop this task', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
                        layer.close(index);
                        JobInfoList.ctrl.operateJob(id, 3);
                    });
                }

            });
            
            $('.btn_operate_restart').on('click',function(){
            	var va_arr = $(this).attr("att").split("|")
                var id = va_arr[0]
                var state = va_arr[1]
                if (state != "2")
                {
                	layer.confirm('Couldn\'t restart task not in state(STARTING, TRANSCODING) ', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
                        layer.close(index);
                    });
                }
                else
                {
                	layer.confirm('Do you really want to restart this task', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
	                    layer.close(index);
	                    JobInfoList.ctrl.operateJob(id, 4);
	                });
                }

            });
            
            $('.btn_operate_start').on('click',function(){
            	var va_arr = $(this).attr("att").split("|")
                var id = va_arr[0]
                var state = va_arr[1]
                if (state != "5")
                {
                	layer.confirm('Couldn\'t start task not in state(STOPED) ', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
                        layer.close(index);
                    });
                }
                else
                {
                	layer.confirm('Do you really want to start this task', {icon: 3, title:'Message', btn:["OK", "Cancel"]}, function(index){
	                    layer.close(index);
	                    JobInfoList.ctrl.operateJob(id, 2);
	                });
                }

            });
           
        },

        // active =1 active 0
        operateJob:function(jobId, op){

            var dat = {"Id":jobId,"Op":op};
            layer.load(1);
            $.ajax({
                url:'/jobinfo/operate',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('Post failed', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                    	layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                            window.location.href ="/jobinfo/list";
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        },
        deleteJobInfo:function(id){

            var dat = {"Id":id};
            layer.load(1);
            $.ajax({
                url:'/jobinfo/delete',
                data:dat,
                dataType:'json',
                error:function(XMLHttpRequest, textStatus, errorThrown){
                    layer.closeAll();
                    layer.msg('Post failed', {
                        icon: 5,
                        time: 2000 //2秒关闭（如果不配置，默认是3秒）
                    }, function(){
                    });

                },
                success:function(data, textStatus, jqXHR){
                    layer.closeAll();
                    if(data.success == true) {
                        layer.msg(data.message, {
                            icon: 1,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                            window.location.href ="/jobinfo/list";
                        });

                    } else {
                        layer.msg(data.message, {
                            icon: 5,
                            time: 2000 //2秒关闭（如果不配置，默认是3秒）
                        }, function(){
                        });
                    }

                },
                type:'POST',
                cache:true

            });
        }
    }
};


$(function(){

    JobInfoList.init();
});
