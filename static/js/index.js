var msgSuffix = "\n";
var conn;

// Operate const for exchange
var QUITprogram = "quit"
var	OPENport    = "open"
var CLOSEport   = "close"
var	WRITEport   = "write"
var	READport    = "read"
var	GETdevice   = "device"
var	DEFAULT     = "default"
var READonly = "readonly"

// Status const
var	OK  = "ok"
var	NOK = "nok"
var OWNER = "owner"

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
    document.getElementById("input").value = "Connection closed!"
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
        var select = document.getElementById("my-select");
        if (select.value == exChangeData.Target) {
          onReadPort(exChangeData.Msg);
        }
        break;
      default:

    }
  }
};

function onGetDevice(device) {
  var select = document.getElementById("my-select"); //get select object
  preValue = select.value;
  select.length=0;
  if (device[0] == "") { //when there's no device , the value of devie is [""]
    document.getElementById("input").value = "There is no serial port found!";
    select.disabled = true;
  }else{
    select.disabled = false;
    for (i in device) {
      var option = document.createElement("option"); //create option object
      option.value = device[i];
      option.innerHTML = device[i];
      select.appendChild(option);  // Add option to select
    }

    for(var i=0; i<select.options.length; i++) {
      if(select.options[i].innerHTML == preValue) {
        select.options[i].selected = true;
        break;
      }
    }
  }
};

function onOpenPort(msg){
  console.log(msg);
  var prompt = "Port is opened!"
  if (msg == NOK) {
    prompt = "Failed to open port!"
  }else if (msg == OWNER){
    console.log("im owner");
    prompt = "You get the write permission"
  }
  document.getElementById("input").value = prompt;
};

function onWritePort(msg){
  if (msg == READonly) {
    alert("You don't have the permission to write!")
  }
  if (msg == NOK) {
    alert("Failed to write!")
  }
};

function onReadPort(msg){
  var output = document.getElementById("output");
  output.appendChild(document.createTextNode(msg));
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
};

function onOpen() {
  var exChangeData = new Object();
  exChangeData.Cmd = GETdevice;
  var exChangeJSON = JSON.stringify(exChangeData);
  conn.send(exChangeJSON);
  console.log(exChangeJSON);
};

function onClickSubmit() {
  console.log("i'm submit");
  var exChangeData = new Object();
  var input = document.getElementById("input");
  var select = document.getElementById("my-select");
  exChangeData.Cmd = WRITEport;
  exChangeData.Msg = input.value+msgSuffix;
  exChangeData.Target = select.value
  console.log("submit :",exChangeData);
  var exChangeJson = JSON.stringify(exChangeData);
  if (!conn) {
    return false;
  }
  console.log("i'm sending");
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

function onEnter(evt) {
  if(evt.keyCode == 13) {
    onClickSubmit()
  }
};
