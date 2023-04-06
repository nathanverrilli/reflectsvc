curl --verbose --insecure ^
 --connect-timeout 750 ^
 --header "Content-Type: application/xml" ^
 --header "Accept: application/json" ^
 --data-binary @body.xml ^
 --verbose ^
 http://localhost:9090/xml2json
