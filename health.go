package main

import (
	"net/http"
)

func health(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = rw.Write([]byte(`
<!DOCTYPE html>
<html lang="en">
	<head>
    	<meta charset="UTF-8">
    	<title>Spoty party bot</title>
	</head>
	<body>
		Hello i am ok
	</body>
</html>
		`))
	return
}
