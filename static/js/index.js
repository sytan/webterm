window.onload = function () {
  var conn;
  var element = document.getElementById("input");
  document.getElementById("form").onsubmit = function () {
    console.log(element.value);
    if (!conn) {
      return false;
    }
    if (!element.value) {
      return false;
    }
    var exChangeMsg = new Object();
    exChangeMsg.Cmd = "plaintext";
    exChangeMsg.Msg = element.value;
    var exChangeJSON = JSON.stringify(exChangeMsg);
    console.log(exChangeJSON);
    conn.send(exChangeJSON);
    return false;
  };
  var wsServer = 'ws://'+location.host+'/ws';
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
        console.log(evt.data);
        var obj = evt.data;
        while (typeof(obj) != "object"){
           obj = JSON.parse(obj);
        }
        console.log(obj);
        var sltObj = document.getElementById("my-select"); //get select object
        if (obj.Cmd == "select") {
          for (x in obj.Msg) {
            var optionObj = document.createElement("option"); //create option object
            optionObj.value = obj.Msg[x];
            optionObj.innerHTML = obj.Msg[x];
            sltObj.appendChild(optionObj);  //添加到select
          }
        }else {
          output.appendChild(document.createTextNode(obj.Msg+'\n'));
        }
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
      alert("i'm not support");
  }
};
