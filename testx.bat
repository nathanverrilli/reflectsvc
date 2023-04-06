curl --verbose --insecure ^
 --header "Content-Type: application/xml" ^
 --data-binary @body.xml ^
 --url http://localhost:9090/xml2json
