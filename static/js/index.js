var msgSuffix = "";
var conn;

// Operate const for exchange
var QUITprogram = "quit"
var	OPENport    = "open"
var CLOSEport   = "close"
var	WRITEport   = "write"
var	READport    = "read"
var	GETdevice   = "device"
var	DEFAULT     = "default"

// Status const
var	OK  = "ok"
var	NOK = "nok"

window.onload = function () {
  if (!window["WebSocket"]) {
    alert("Your browser does not support websocket!");
    return
  }
  var wsServer = 'ws://'+location.host+'/ws';
  conn = new WebSocket(wsServer);
  conn.onopen = function(evt) {
    onOpen();
  };
  conn.onclose = function (evt) {
    document.getElementById("output").value = "Connection closed!"
  };
  conn.onerror = function(evt) {
    alert("i'm error");
    alert(evt.data);
  };
  conn.onmessage = function(evt) {
    var exChangeData = toJson(evt.data);
    console.log(exChangeData);

    switch (exChangeData.Cmd) {
      case GETdevice:
        onGetDevice(exChangeData.Msg);
        break;
      case OPENport:
        onOpenPort(exChangeData.Msg);
        break;
      case WRITEport:
        onWritePort(exChangeData.Msg);
        break;
      case READport:
        onReadPort(exChangeData.Msg);
        break;
      default:

    }
  }
};

function onGetDevice(device) {
  var select = document.getElementById("my-select"); //get select object
  for (i in device) {
    var option = document.createElement("option"); //create option object
    option.value = device[i];
    option.innerHTML = device[i];
    select.appendChild(option);  // Add option to select
  }
};

function onOpenPort(msg){
  var prompt = "Port is opened!"
  if (msg != OK) {
    prompt = "Failed to open port!"
  }
  document.getElementById("output").value = prompt;
};

function onWritePort(msg){
  if (msg != OK) {
    alert("Failed to write to port!")
  }
};

function onReadPort(msg){
  var output = document.getElementById("output");
  output.appendChild(document.createTextNode(obj.Msg));
  output.scrollTop = output.scrollHeight;
};


function toJson(data) {
  var dataJson = data;
  while (typeof(dataJson) != "object"){
     dataJson = JSON.parse(dataJson);
  }
  return dataJson
};
function addonSwitch() {
  addon = document.getElementById("addon");
  var suffix = {"LF":"\n", "CR":"\r", "NULL":""};
  var arr = ["LF","CR","NULL"];
  var index = arr.indexOf(addon.innerText);
  var newindex = (index+1)%3;
  addon.innerText = arr[newindex];
  msgSuffix = suffix[arr[newindex]];
  console.log(msgSuffix);
};

function onOpen() {
  var exChangeData = new Object();
  exChangeData.Cmd = GETdevice;
  var exChangeJSON = JSON.stringify(exChangeData);
  conn.send(exChangeJSON);
  console.log(exChangeJSON);
};

function onClickSubmit() {
  var exChangeData = new Object();
  var input = document.getElementById("input");
  var select = document.getElementById("my-select");
  exChangeData.Cmd = WRITEport;
  exChangeData.Msg = input.value+msgSuffix;
  exChangeData.Target = select.value
  var exChangeJson = JSON.stringify(exChangeData);
  if (!conn) {
    return false;
  }
  conn.send(exChangeJson);
};

function onClickSelect() {
  var exChangeData = new Object();
  var select = document.getElementById("my-select");
  exChangeData.Cmd = OPENport;
  exChangeData.Target = select.value
  var exChangeJson = JSON.stringify(exChangeData);
  if (!conn) {
    return false;
  }
  conn.send(exChangeJson);
};
