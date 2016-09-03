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
      <form role="form" id="form">
        <div class="form-group">
          <label for="output">Output</label>
          <textarea id="output" class="form-control" rows="10" disabled="true"></textarea>
        </div>
        <div class="form-group">
          <input type="text" class="form-control" placeholder="command" id="input">
          <button type="submit" class="btn btn-default">Submit</button>
        </div>
      </form>
    </div><!-- /.container -->


    <!-- Bootstrap core JavaScript
    ================================================== -->
    <!-- Placed at the end of the document so the pages load faster -->
    <script src="../static/js/jquery-3.1.0.js"></script>
    <script src="../static/bootstrap-3.3.5/js/bootstrap.min.js"></script>
    <script src="../static/js/index.js"></script>
  </body>
</html>
