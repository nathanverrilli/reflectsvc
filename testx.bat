curl --insecure ^
 --header "Content-Type: application/xml" ^
 --data-binary @body.xml ^
 --request POST ^
 --verbose ^
 https://localhost:9090/convert
