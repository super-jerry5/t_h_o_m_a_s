/**
 * Created by banxia on 16/5/17.
 */


/**
 * Created by banxia on 16/5/16.
 */


var EditJob = {

    init:function(){

        EditJob.ctrl.initEvent();
        EditJob.ctrl.initValidate();
    },
    ctrl:{
        initEvent:function(){


        },
        initValidate:function(){

            $('#editForm').bootstrapValidator({
                message: 'Invalid Parameters',
                feedbackIcons: {
                    valid: 'glyphicon glyphicon-ok',
                    invalid: 'glyphicon glyphicon-remove',
                    validating: 'glyphicon glyphicon-refresh'
                },
                fields:{

                	Name:{
                        message: 'Task Name Invalid',
                        validators: {
                            notEmpty: {
                                message: 'Task Name couldn\'t be empty'
                            },
                            stringLength: {
                                min: 3,
                                max: 30,
                                message: 'Task Name string length should between 3 and 30'
                            }

                        }
                    },
                    InputUrl:{
                        message: 'InputUrl Invalid',
                        validators: {
                            notEmpty: {
                                message: 'InputUrl couldn\'t be empty'
                            },
                            stringLength: {
                                min: 3,
                                max: 512,
                                message: 'InputUrl string length shoudl between 3 and 512'
                            }

                        }
                    },
                    InputUrlBackup:{
                        message: 'InputUrlBackup Invalid',
                        validators: {
                        	stringLength: {
                                min: 3,
                                max: 512,
                                message: 'InputUrlBackup string length shoudl between 3 and 512'
                            }

                        }
                    },
                    OutputUrl:{
                        message: 'OutputUrl Invalid',
                        validators: {
                            notEmpty: {
                                message: 'OutputUrl couldn\'t be empty'
                            },
                            stringLength: {
                                min: 3,
                                max: 512,
                                message: 'InputUrlBackup string length shoudl between 3 and 512'
                            }

                        }
                    },
                    Scale:{
                        message: 'Scale Invalid',
                        validators: {
                            notEmpty: {
                                message: 'Scale couldn\'t be empty'
                            }
                        }
                    },
                    VcodecBitrate:{
                        message: 'VcodecBitrate Invalid',
                        validators: {
                            notEmpty: {
                                message: 'VcodecBitrate couldn\'t be empty'
                            }
                        }
                    },
                    AudioParam:{
                        message: 'AudioParam Invalid',
                        validators: {
                            notEmpty: {
                                message: 'AudioParam couldn\'t be empty'
                            }
                        }
                    },
                    VcodecType:{
                        message: 'VcodecType Invalid',
                        validators: {
                            notEmpty: {
                                message: 'VcodecType couldn\'t be empty'
                            }
                        }
                    }
                }
            }).on('success.form.bv', function(e) {
                // Prevent form submission
                e.preventDefault();

                // Get the form instance
                var $form = $(e.target);
                // Get the BootstrapValidator instance
                var bv = $form.data('bootstrapValidator');

                var dat = $form.serialize();
                layer.load(1);
                $.ajax({
                    url:'/jobinfo/edit',
                    data:dat,
                    dataType:'json',
                    error:function(XMLHttpRequest, textStatus, errorThrown){
                        layer.closeAll();
                        $form
                            .bootstrapValidator('disableSubmitButtons', false)
                    },
                    success:function(data, textStatus, jqXHR){
                        layer.closeAll();
                        $form
                            .bootstrapValidator('disableSubmitButtons', false)
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


            });
        }
    }
};
$(function(){


    EditJob.init();

});

