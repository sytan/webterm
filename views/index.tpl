<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>Web terminal</title>

    <!-- Bootstrap core CSS -->
    <link href="../static/bootstrap-3.3.5/css/bootstrap.min.css" rel="stylesheet">

    <!-- Custom styles for this template -->
    <link href="../static/css/index.css" rel="stylesheet">

  </head>

  <body>
    <div class="container">
      <form role="form" id="form" onsubmit="return false;"> <!-- onsubmit to avoid submit the form when enter text input -->
        <div class="form-group">
          <label for="output">Output</label>
          <textarea id="output" class="form-control" rows="15" disabled="true"></textarea>
        </div>

          <div class="input-group">
              <span class="input-group-addon" id = "addon" onclick="addonSwitch()">NULL</span>
              <input type="text" class="form-control" placeholder="command" id="input" onkeydown= "onEnter(event)">
            <div class="input-group-btn">
              <!-- Button and dropdown menu -->
              <button type="button" class="btn btn-default dropdown-toggle" data-toggle="dropdown">
                <span class="caret"></span>
                <span class="sr-only">Toggle Dropdown</span>
              </button>
              <ul class="dropdown-menu dropdown-menu-right" role="menu">
                <li><a href="#">LR</a></li>
                <li><a href="#">CR</a></li>
                <li><a href="#">Something else here</a></li>
                <li class="divider"></li>
                <li><a href="#">Separated link</a></li>
              </ul>
              <button type="button" class="btn btn-default" onclick="onClickSubmit()">Submit</button>
            </div>
          </div><!-- /input-group -->
      </form>
      <select id="my-select" class="form-control" onclick="onClickSelect()">
      </select>
    </div><!-- /.container -->

    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="../static/js/jquery-3.1.0.js"></script>
    <script src="../static/bootstrap-3.3.5/js/bootstrap.min.js"></script>
    <script src="../static/js/index.js"></script>
  </body>
</html>
