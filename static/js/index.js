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
        output.appendChild(document.createTextNode(evt.data+'\n'));
          // var messages = evt.data.split('\n');
          // for (var i = 0; i < messages.length; i++) {
          //     var item = document.createElement("div");
          //     item.innerText = messages[i];
          //     appendLog(item);
          // }
      };
      conn.onerror = function(evt) {
        alert(evt.data);
      };
  } else {
      alert("i'm not support")
  }
};
