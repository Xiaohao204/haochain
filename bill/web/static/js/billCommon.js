$(document).ready(function() {
  // 顶部按钮
  var divSpan = $(".bill .top div>span:nth-of-type(2)");
  $('.bill .top div span:nth-of-type(1)').on('click',function() {
    if(divSpan.css("display")=="none"){
      divSpan.show();
    }else{
      divSpan.hide();
    }
  });
  // 退出
  $('.exit').on('click',function() {
    $(location).attr('href','./login.html')
  })

  // 折叠菜单
  function filterSpaceNode(node){
		for(var i = 0;i<node.childNodes.length;i++){
				node.removeChild(node.childNodes[i]);
		}
	}

	var box = document.querySelector("#box");
	filterSpaceNode(box);
	var h3 = document.querySelector("#box h3");
	var div = document.querySelector("#box div");
		h3.onclick = function(){
			 none();
			 if(this.nextSibling.style.display=="block"){
				this.nextSibling.style.display="none";
        $('h3 span img').attr('src','../images/jt.png');
			 }else{
				this.nextSibling.style.display = "block";
        $('h3 span img').attr('src','../images/jt2.png');
			 }
	}
	function none(){
			div.style.display.display="none";
	}

  // 发布票据
  $('#list li:nth-of-type(1)').on('click',function() {
    $(location).attr('href','./issue.html');
  })
  // 我的票据
  $('#list li:nth-of-type(2)').on('click',function() {
    $(location).attr('href','./bills.html');
  })
  $('#list li:nth-of-type(3)').on('click',function() {
    $(location).attr('href','./waitEndorses.html');
  })


  // 弹框，提交票据 发布按钮
 
  $('.issueBtn').on('click',function() {
    
          $('#popout').css('display','block');
          $('.popout div').html('<p>提交成功</p><p>aaa111a1111111111f1111111r1112222222ffff22222222222222222ccvcv22222</p>');
       
  })
  

  // 发起背书按钮
  $('#endorseReq').on('click',function() {
      $('#popout').css('display','block');
      $('.popout div').html('<p>背书成功</p><p>aaa111a1111111111f1111111r1112222222ffff22222222222222222ccvcv22222</p>');
  })
  // 票据列表关闭按钮
  $('#myBillDetail>p:nth-of-type(3) span:nth-of-type(2)').on('click',function() {
    $(location).attr('href','./bills.html');
  })

  // 待签收票据详情按钮
  $('.myUnBill .detail').on('click',function() {
      $(location).attr('href','./waitEndorseInfo.html')
  })

  // 签收背书
  $('#myUnBillDetail>p span:nth-of-type(1)').on('click',function() {
    $('#popout').css('display','block');
    $('.popout div').html('<p>已签收成功</p><p>aaa111a1111111111f1111111r1112222222ffff22222222222222222ccvcv22222</p>');
  })
  // 拒绝背书
  $('#myUnBillDetail>p span:nth-of-type(2)').on('click',function() {
    $('#popout').css('display','block');
    $('.popout div').html('<p>已拒绝成功</p><p>aaa111a1111111111f1111111r1112222222ffff22222222222222222ccvcv22222</p>');
  })
  // 待背书列表关闭按钮
  $('#myUnBillDetail>p span:nth-of-type(3)').on('click',function() {
      $(location).attr('href','./waitEndorse.html')
  })
})
