curl --insecure ^
 --header "Content-Type: application/xml" ^
 --data-binary @body.xml ^
 --request POST ^
 https://localhost:9090/parsifal
