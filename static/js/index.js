window.onload = function () {
  var conn;
  var msg = document.getElementById("input");
  document.getElementById("form").onsubmit = function () {
      if (!conn) {
          return false;
      }
      if (!msg.value) {
          return false;
      }
      conn.send(msg.value);
      msg.value = "";
      return false;
  };
  var wsServer = 'ws://'+location.host+'/ws'
  if (window["WebSocket"]) {
      var output = document.getElementById("output");
      conn = new WebSocket(wsServer);
      conn.onopen = function(evt) {output.appendChild(document.createTextNode("Ready,Go!\n"))};
      conn.onclose = function (evt) {
          var item = document.createElement("div");
          item.innerHTML = "<b>Connection closed.</b>";
          appendLog(item);
      };
      conn.onmessage = function (evt) {
        var obj=JSON.parse(evt.data);
        obj = JSON.parse(obj);
        var sltObj = document.getElementById("my-select"); //get select object
        for (x in obj) {
            var optionObj = document.createElement("option"); //create option object
            optionObj.value = obj[x];
            optionObj.innerHTML = obj[x];
            sltObj.appendChild(optionObj);  //添加到select
        }
        output.appendChild(document.createTextNode(evt.data+'\n'));
          // var messages = evt.data.split('\n');
          // for (var i = 0; i < messages.length; i++) {
          //     var item = document.createElement("div");
          //     item.innerText = messages[i];
          //     appendLog(item);
          // }
      };
      conn.onerror = function(evt) {
        alert("i'm error");
        alert(evt.data);
      };
  } else {
      alert("i'm not support")
  }
};
