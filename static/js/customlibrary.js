var   selectedid=0;			//记录ID
var   Operationbtn=0;		//1添加，2编辑
var   checknum=3;		    //可选按钮个数
//显示层函数
function showfun(){
	var myAlert = document.getElementById("box01"); 
	 myAlert.style.display = "block"; 
	 myAlert.style.position = "absolute"; 
	 myAlert.style.top = "50%"; 
	 myAlert.style.left = "50%"; 
	 myAlert.style.marginTop = "-75px"; 
	 myAlert.style.marginLeft = "-150px";
	 
	 mybg = document.createElement("div"); 
	 mybg.setAttribute("id","mybg"); 
	 mybg.style.background = "#000"; 
	 mybg.style.width = "100%"; 
	 mybg.style.height = "100%"; 
	 mybg.style.position = "static"; 
	 mybg.style.top = "0"; 
	 mybg.style.left = "0"; 
	 mybg.style.zIndex = "0"; 
	 mybg.style.opacity = "0.3"; 
	 mybg.style.filter = "Alpha(opacity=30)"; 
	 document.body.appendChild(mybg);
	 mybg.style.display = "none";
	 
	 document.body.style.overflow = "hidden"; 
}

//隐藏层函数
function _Funreturn(){
	var myAlert = document.getElementById("box01"); 
	myAlert.style.display = "none"; 
	mybg.style.display = "none";
}


//自定义消息页面控制
function _deletecustom(id) {
    $.ajax({
        url: "/cstore/newslib/newslibdeletecustom",
        data: {"Id": id},
        type: "post",
        success: function (r) {
            if (r.status) {
                oTable.fnReloadAjax(oTable.fnSettings());
            } else {
                alert(r.info);
            }
        }, error: function (error) {
            console.log(error);
        }
    });
}


function _addcustom(){
	showfun();
	Operationbtn = 1;

	$("#showaddcustom").text("请填写自定义消息内容");
	$("#subbutton").val("添    加");


	$("textarea[name='customcontent']").val('');//清空内如编辑框
}

function _editcustom(id,Content) {
	showfun();
	selectedid =id;
	Operationbtn = 2;
	$("#showaddcustom").text("自定义信息修改");
	$("#subbutton").val("修    改");


	 if(Content != "" && Content != null){
	 	//$("#econtent").val(Content);
	 	content:$("textarea[name='customcontent']").val(Content);
	 }
}

//添加函数
function checkcustom() {
	var flag=true;
	var content=$("textarea[name='customcontent']").val();	
	if(content=="" || content == null) {
		alert("内容不能为空!");
		flag=false;
	}

	return flag;
}

function _subcustom(){
	if(!checkcustom()){
		return;
	}

	switch(Operationbtn)
	{
		case 1:
			$.ajax({
				async:true, //请勿改成异步，下面有些程序依赖此请数据
		  		type : "POST", 
		  		data:{
		  			content:$("textarea[name='customcontent']").val()
		  		},
		  		url: "/cstore/newslib/newslibaddcustom",
		  		dataType:'json',
				success:function(json){
					alert("添加成功!")
					$("textarea[name='customcontent']").val('');
					oTable.fnReloadAjax(oTable.fnSettings());
				},
				error: function(error){
					alert("添加失败"); 
					console.log(error);
				}  
		    },'json');
			break;
		case 2:
			$.ajax({
				async:true, //请勿改成异步，下面有些程序依赖此请数据
	      		type : "POST", 
	      		data:{
	      			id:selectedid,
	      			content:$("textarea[name='customcontent']").val()
	      		},
	      		url: "/cstore/newslib/newslibeditcustom",
	      		dataType:'json',
				success:function(json){
					alert("编辑成功!");
					oTable.fnReloadAjax(oTable.fnSettings());
				},
				error: function(error){
					alert("编辑失败"); 
					console.log(error);
				}  
			},'json');
			break;
		default:
			break;
	}
}



///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//打招呼页面控制
function _deletemanagement(id) {
    $.ajax({
        url: "/cstore/newslib/newslibdeletehello",
        data: {"Id": id},
        type: "post",
        success: function (r) {
            if (r.status) {
                oTable.fnReloadAjax(oTable.fnSettings());
            } else {
                alert(r.info);
            }
        }, error: function (error) {
            console.log(error);
        }
    });
}

function _editmanagement(id,Content,newstype) {
	showfun();
	selectedid =id;
	Operationbtn = 2;
	$("#showaddmanagement").text("打招呼信息修改");
	$("#subbutton").val("修    改");


	//设置可选按钮状态
	var i=0,j=0;
    for(i=1;i<checknum+1;i++){
    	if(newstype.length<3){
	    	$("input[name=selectcheck"+i+"]").prop({"checked":false});
	    	$("input[name=selectcheck"+i+"]").val("0")
    	}else {
	    	j=parseInt(newstype[i-1]); 
	    	if(1==j){
	    		$("input[name=selectcheck"+i+"]").prop({"checked":true});
	    		$("input[name=selectcheck"+i+"]").val("1")
	    	}else{
	    		$("input[name=selectcheck"+i+"]").prop({"checked":false});
	    		$("input[name=selectcheck"+i+"]").val("0")
	    	}
    	}
	}

	//加载内容编辑框
	if(Content != "" && Content != null) {
	 	$("textarea[name='contenttext']").val(Content);
	}else{
		$("textarea[name='contenttext']").val("");
	}
}


function _CheckType(){
	var attrdata = $(this).val();
	alert(attrdata);



	// var context  =$("textarea[name='contenttext']").val();
	// var showtext = context;//要显示的内容
	// var convalue = new Array(checknum);
	// convalue[0]  = '{Nickname}';
	// convalue[1]  = '{Realname}';
	// convalue[2]  = '{Area}';

	// var checkvalue=new RegExp(convalue[num-1]);
	// if($("input[name=selectcheck"+num+"]").prop("checked")){
	// 	if(checkvalue.test(showtext)){
	// 		showtext =convalue[num-1]+"已存在，是否继续！";
	// 		if(confirm(showtext)){
	// 			showtext = context+convalue[num-1];
	// 		}else{
	// 			showtext = context;
	// 		}
	// 	}else{
	// 		showtext = context+convalue[num-1];
	// 	}
	// 	$("input[name=selectcheck"+num+"]").val("1")
	// }else{
	// 	showtext ="是否删除存在的"+convalue[num-1];
	// 	if(confirm(showtext)){
	// 		showtext = context;
	// 		while(checkvalue.test(showtext)){
	// 		  showtext = showtext.replace(convalue[num-1],"");
	// 		}
	// 	}else{
	// 		showtext = context;
	// 	}
	// 	$("input[name=selectcheck"+num+"]").val("0")
	// }
	// $("textarea[name='contenttext']").val(showtext);
}

function _addmanagement(){
	showfun();
	Operationbtn = 1;

	$("#showaddmanagement").text("请填写打招呼信息");
	$("#subbutton").val("添    加");


    //初始化可选按钮
    for(var i=1;i<checknum+1;i++){
    	$("input[name=selectcheck"+i+"]").prop({"checked":false});
    	$("input[name=selectcheck"+i+"]").val("0")
	}
	$("textarea[name='contenttext']").val('');//清空内如编辑框
}

//检查编辑框
function checkmanagement() {
	var flag=true;

	var content=$("textarea[name='contenttext']").val();
	if(content == "" || content == null) {
		alert("内容不能为空，请输入！");
		flag=false;
	}
	return flag;
}

//表单提交函数
function _submanagement(){
	if(!checkmanagement()){
		return;
	}

	var i=0,typevalue="";
	for(i=1;i<checknum+1;i++){
		typevalue+=$("input[name=selectcheck"+i+"]").val();
	}
	switch(Operationbtn){
		case 1:
			$.ajax({
				async:true, //请勿改成异步，下面有些程序依赖此请数据
	      		type : "POST", 
	      		data:{
	      			content:$("textarea[name='contenttext']").val(),
	      			newstype:typevalue
	      		},
	      		url: "/cstore/newslib/newslibaddhello",
	      		dataType:'json',
				success:function(json){
					alert("添加成功");
					oTable.fnReloadAjax(oTable.fnSettings());

				    //初始化可选按钮
				    for(i=1;i<checknum+1;i++){
				    	$("input[name=selectcheck"+i+"]").prop({"checked":false});
				    	$("input[name=selectcheck"+i+"]").val("0")
					}
					$("textarea[name='contenttext']").val('');//清空内如编辑框
				},
				error: function(error){
					alert("添加失败"); 
					console.log(error);
				}  
			},'json');
			break;
		case 2:
			$.ajax({
				async:true, //请勿改成异步，下面有些程序依赖此请数据
	      		type : "POST", 
	      		data:{
	      			id:selectedid,
	      			content:$("textarea[name='contenttext']").val(),
	      			newstype:typevalue
	      		},
	      		url: "/cstore/newslib/newslibedithello",
	      		dataType:'json',
				success:function(json){
					alert("编辑成功");
					oTable.fnReloadAjax(oTable.fnSettings());
				},
				error: function(error){
					alert("编辑失败"); 
					console.log(error);
				}  
			},'json');
			break;
		default:
			break;
	}
}


///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//系统限制设置
//检查控件
function checkset() {
	var flag=true;

	var inumlimit 		=$("input[name='numlimit']").val();
	var ifrelimit 		=$("input[name='frelimit']").val();
	var ifrelimittime	=$("input[name='frelimittime']").val();


	if(inumlimit == "" || inumlimit == null) {
		alert("次数限制不能为空！");
		flag=false;
	}

	if(ifrelimit == "" || ifrelimit == null) {
		alert("频率限制不能为空！");
		flag=false;
	}
	
	if(ifrelimittime == "" || ifrelimittime == null) {
		alert("时间设置不能为空！");
		flag=false;
	}
	return flag;
}


function _Funsystemset(){
		if(!checkset()){
			return;
		}
		$.ajax({
			async:true, //请勿改成异步，下面有些程序依赖此请数据
      		type : "POST", 
      		data:{
      			inumlimit:$("input[name='numlimit']").val(),
      			ifrelimit:$("input[name='frelimit']").val(),
      			ifrelimittime:$("input[name='frelimittime']").val(),
      			bsetnum:$("input[name='setnumtype']:checked").val(),
      			bsetfre:$("input[name='setfretype']:checked").val()
      		},
      		url: "/cstore/newslib/newslieditbset",
      		dataType:'json',
			success:function(json){
				alert("设置成功"); 
			},
			error: function(error){
				alert("设置失败"); 
				console.log(error);
			}  
		},'json');
	}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
//以下判断控件是输入数字，字母，汉字
  $.fn.aa = function () {
      $(this).keypress(function (event) {
         var eventObj = event || e;
          var keyCode = eventObj.keyCode || eventObj.which;
         if ((keyCode >= 48 && keyCode <= 57))
             return true;
         else
             return false;
     }).focus(function () {
     //禁用输入法
         this.style.imeMode = 'disabled';
     }).bind("paste", function () {
     //获取剪切板的内容
         var clipboard = window.clipboardData.getData("Text");
        if (/^\d+$/.test(clipboard))
             return true;
         else
             return false;
     });
 };

  // <summary>
  // 限制只能输入字母
  // </summary>
  // ----------------------------------------------------------------------
  $.fn.bb = function () {
      $(this).keypress(function (event) {
          var eventObj = event || e;
          var keyCode = eventObj.keyCode || eventObj.which;
         if ((keyCode >= 65 && keyCode <= 90) || (keyCode >= 97 && keyCode <= 122))
             return true;
         else
             return false;
     }).focus(function () {
         this.style.imeMode = 'disabled';
         }).bind("paste", function () {
         var clipboard = window.clipboardData.getData("Text");
         if (/^[a-zA-Z]+$/.test(clipboard))
             return true;
         else
             return false;
     });
 };

  // ----------------------------------------------------------------------
  // <summary>
  // 限制只能输入数字和字母
  // </summary>
  // ----------------------------------------------------------------------
  $.fn.cc = function () {
      $(this).keypress(function (event) {
          var eventObj = event || e;
          var keyCode = eventObj.keyCode || eventObj.which;
          if ((keyCode >= 48 && keyCode <= 57) || (keyCode >= 65 && keyCode <= 90) || (keyCode >= 97 && keyCode <= 122))
             return true;
         else
             return false;
     }).focus(function () {
         this.style.imeMode = 'disabled';
     }).bind("paste", function () {
         var clipboard = window.clipboardData.getData("Text");
         if (/^(\d|[a-zA-Z])+$/.test(clipboard))
             return true;
         else
             return false;
     });
};

/*
  <ul>
         <li>只能输入数字：<input type="text" class="onlyNum" /></li>
         <li>只能输入字母：<input type="text" class="onlyAlpha" /></li>
         <li>只能输入数字和字母：<input type="text" class="onlyNumAlpha" /></li>
  </ul>
  */

 $(function () {
     // 限制使用了onlyNum类样式的控件只能输入数字
     $("input[name='numlimit']").aa();
     $("input[name='frelimit']").aa();
     $("input[name='frelimittime']").aa();

     //限制使用了onlyAlpha类样式的控件只能输入字母

     //$("#emicrosignal").bb();

     // 限制使用了onlyNumAlpha类样式的控件只能输入数字和字母

     //$("#emicrosignal").cc();

    });