package function

import (
	// "fmt"
	// "io/ioutil"
	"net/http"

	"github.com/openfaas-incubator/go-function-sdk"
)

func Handle(req handler.Request) (handler.Response, error) {
	var err error

	return handler.Response{
		Body:       []byte(html()),
		StatusCode: http.StatusOK,
	}, err
}

func html() string {
	index := `
<!doctype html>
<html>
<head>
	<meta charset="utf-8">

	<title></title>

	<script src="https://cdn.jsdelivr.net/npm/canvg/dist/browser/canvg.min.js"></script>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.3.1/jquery.min.js"></script>
	<script type="text/javascript">
	function send(){
	  $.ajax({
	    type: "get",
	    url: "http://go3dprint.192.168.99.101.nip.io/function/mesh",
	    success:function(data)
	    {
	        // console.log(data);
			// document.getElementById("text").innerHTML = data;
			canvg(document.getElementById('canvas'), data);
	        setTimeout(function(){
	            send();
	        }, 1000);
	    }
	})};
	send()

	</script>

	<style type="text/css">
	</style>

</head>

<body>
	<h1 id="text">Test text.</h1>
	<canvas id="canvas" width="800px" height="600px"></canvas>
</body>

</html>
`
	return index
}
