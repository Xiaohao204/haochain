$(document).ready(function() {
  $('#btn').on('click',function() {
    var name = $('#name').val();
    var pwd = $('#pwd').val();
    if(name == '' || pwd == ''){
      $('#popout').css('display','block');
      $('.popout div').text('请输入用户名/密码');
    }else{
	  location.href='./bills.html'
	}
  })
})
