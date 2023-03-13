REM Create the private/public key pair
openssl genrsa -out host.pem 4096
REM Create the public key
openssl rsa -in host.pem -pubout -out host.pub
REM Create the private key
openssl rsa -in host.pen -out host.key
REM Create the SSL certificate
 openssl req -new -x509 -nodes -sha256 -days 999 -key host.key -out host.crt
