//建立websocket连接
socket = io();
socket.emit('login', user);
socket.on('wechatmsg', function(objData) {
	var objArr = eval(objData);
	for(var i=0;i<objArr.length;i++){
		var obj = objArr[i];
		var _mywechat = obj.mywechat;
		var _friendwechat = obj.friendwechat;
		var content = obj.content;
		var time = new Date(parseInt(obj.time) * 1000).format("yyyy-MM-dd HH:mm:ss");
		var _hid = obj.hid;
		var type = obj.type;

		//新来的联系人消息
		if($(".people[data-friendwechat="+_friendwechat+"]")[0] == undefined){
			$(".prompt").hide();
			var peopleHtml = '<li class="people" data-mywechat="'+_mywechat+'" data-friendwechat="'+_friendwechat+'" data-hid="'+_hid+'"><img src="/static/images/unknow.jpg" width=30px height=30px style="border-radius: 50%;" /><span class="badge active">1</span><span>'+_friendwechat+'</span><i class="fa fa-times pull-right"></i></li>';
			$("#contact-list").prepend(peopleHtml);

			//创建聊天窗口
			var chatDivHtml = "<div class='chatDiv' data-friendwechat='"+_friendwechat+"' data-mywechat='"+_mywechat+"' data-hid='"+_hid+"'></div>";
			$("#chatWindow").append(chatDivHtml);

			//添加消息
			var msgHtml = getMsgHtml(type,content,time,0);
			$(".chatDiv[data-friendwechat="+_friendwechat+"][data-mywechat="+_mywechat+"]").append(msgHtml);
		}else{
			//添加消息
			var msgHtml = getMsgHtml(type,content,time,0);
			$(".chatDiv[data-friendwechat="+_friendwechat+"][data-mywechat="+_mywechat+"]").append(msgHtml);

			if($(".people[data-friendwechat="+_friendwechat+"]").hasClass("active")){//会话激活中
				//消息记录自动滚动到最底部
				var contentH = $(".chatDiv[data-friendwechat="+_friendwechat+"][data-mywechat="+_mywechat+"]").get(0).scrollHeight;
				$(".chatDiv[data-friendwechat="+_friendwechat+"][data-mywechat="+_mywechat+"]").animate({scrollTop:contentH+"px"},300);
			}else{//会话未激活
				var nowCnt = $(".people[data-friendwechat="+_friendwechat+"]").find(".badge").html();
				nowCnt++;
				$(".people[data-friendwechat="+_friendwechat+"]").find(".badge").html(nowCnt);

				if(!$(".people[data-friendwechat="+_friendwechat+"]").find(".badge").hasClass("active")){
					$(".people[data-friendwechat="+_friendwechat+"]").find(".badge").addClass("active");
				}

				//最新的消息最前显示
				var tempDom = $(".people[data-friendwechat="+_friendwechat+"]").detach();
				$("#contact-list").prepend(tempDom);
			}
		}
	}
});

socket.on("msgState",function(obj){
	//更改消息状态提示
	if(obj.state){
		$("#msg"+obj.msgId).find(".sendState").remove();
	}else{
		layer.msg("发送消息失败！");
		$("#msg"+obj.msgId).find(".sendState").removeClass("fa-spinner");
		$("#msg"+obj.msgId).find(".sendState").removeClass("fa-spin");
		$("#msg"+obj.msgId).find(".sendState").addClass("fa-repeat");
		$("#msg"+obj.msgId).find(".sendState").one("click",function(){
			//发送失败消息重发
			$(this).removeClass("fa-repeat");
			$(this).addClass("fa-spin");
			$(this).addClass("fa-spinner");

			var msgObj = {
				msgId:obj.msgId+"",
				msg:obj.msg
			}
			socket.emit('wechatmsg', msgObj);
		});
	}
});

//发送消息
function sendMsg(msg,mywechat,target,type,time){
	var obj = {
		msg:msg,
		mywechat:mywechat,
		target:target,
		type:type+"",
		time:time+"",
	}
	var msgObj = {
		msgId:msgId+"",
		msg:JSON.stringify(obj)
	}
	socket.emit('wechatmsg', msgObj);
	msgId++;
}

//联系人点击事件
$(document).on("click",".people",function(){
	mywechat = $(this).attr("data-mywechat");
	target = $(this).attr("data-friendwechat");
	hid = $(this).attr("data-hid");
	$(".chatDiv").hide();
	$(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").show();

	$(this).siblings().removeClass("active");

	$(this).addClass("active");
	$(this).find(".badge").removeClass("active");
	$(this).find(".badge").html("0");

	//自动滑到最底部
	var contentH = $(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").get(0).scrollHeight;
    $(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").animate({scrollTop:contentH+"px"},300);
});

//获取聊天消息的html
function getMsgHtml(type,content,time,sendFlag){//sendFlag 0收到 1发送
	var msgHtml = "";
	if(sendFlag == 0){
		msgHtml = "<div class='sender'><div><img src='/static/images/unknow.jpg'></div><div><div class='left_triangle'></div><p>" + time + "</p><span>" + content + "</span></div></div>";
	}else{
		msgHtml = "<div class='receiver' id='msg"+msgId+"'><div><img src='/static/images/salesman.jpg'></div><div><div class='right_triangle'></div><p>" + time + "</p><span>" + content + "</span></div><i class='sendState fa fa-spinner fa-spin'></i></div>";
	}
	return msgHtml;
}

//关闭会话
$(document).on("click",".people > i",function(e){
	var _mywechat = $(this).parent().attr("data-mywechat");
	var _friendwechat = $(this).parent().attr("data-friendwechat");
	$(".chatDiv[data-mywechat="+_mywechat+"][data-friendwechat="+_friendwechat+"]").remove();
	$(this).parent().remove();
	e.stopPropagation();
});

//发送消息
$("#sendMsgBtn").on("click",function(){
    var msg = $("#sendMsgText").val();
    if(msg == ""){
        alert("不能发送空消息!");
        return;
    }

    if(target == ""){
        alert("请选择联系人！");
        return;
    }

    var time = new Date().format("yyyy-MM-dd HH:mm:ss"); //发送消息时间
    var timestamp = Date.parse(new Date()) / 1000;
    var msgHtml = getMsgHtml(1,msg,time,1);
    $(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").append(msgHtml);//添加消息到页面
    var contentH = $(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").get(0).scrollHeight;
    $(".chatDiv[data-mywechat="+mywechat+"][data-friendwechat="+target+"]").animate({scrollTop:contentH+"px"},300);
    $("#sendMsgText").val("");//清空发送消息框

    sendMsg(msg,mywechat,target,1,timestamp);
});

//点击查看设备信息
$("#detailInfo").on("click",function(){
	if(target == ""){
		alert("请选择联系人！");
        return;
	}

	//页面层
	layer.open({
	    type: 1,
	    skin: 'layui-layer-rim', //加上边框
	    area: ['420px', '240px'], //宽高
	    content: '<div class="form-inline"><label>微信号：</label><span>'+mywechat+'</span></div><div class="form-inline"><label>好友微信号：</label><span>'+target+'</span></div><div class="form-inline"><label>设备号：</label><span>'+hid+'</span></div>'
	});
});