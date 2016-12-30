window.onload = function(){
  startTime();
}
function startTime()
  {
  var today=new Date()
  var y = today.getFullYear()
  var mouth = today.getMonth()+1
  var d = today.getDate()
  var h=today.getHours()
  var m=today.getMinutes()
  var s=today.getSeconds()
  // add a zero in front of numbers<10

  h=checkTime(h)
  m=checkTime(m)
  s=checkTime(s)
  $('.showtime').text(y+"-"+mouth +"-"+d+" "+h+":"+m+":"+s)
  t=setTimeout('startTime()',500)
}

function checkTime(i){
  if (i<10) 
    {i="0" + i}
    return i
}
