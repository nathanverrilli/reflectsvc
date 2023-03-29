# REFLECTSV Usage

Reflectsvc is a mini http service supporting
test and one primary endpoint. All endpoints 
expect POST HTTP(s) requests.

`reflectsvc` logs to `reflectsvc.log` and it
truncates the log file when it starts (data from prior runs
of the service are discarded on start).

The miniservice listens for `HTTPS` connections if it has 
a valid certificate file and private key file; otherwise it
reverts to `HTTP`.

## Endpoints and function
### /reverse
This is a trivial string reversal service, to test
that the server is up and running. 
#### Usage
`curl --header "Content-Type: application/json" --data-binary "{\"s\":\"drawkcaB\"}" ^
--request POST http://localhost:9090/reverse`

### /reflect
This reflects the contents of a POST request to the user,
as well as writing the POST contents to the log and stdout
for examination.

`curl --header "Content-Type: application/json" 
--data-binary @body.json 
--request POST 
http://localhost:9090/reflect`

### /convert && /parsifal
These endpoints take XML data, show the XML data, and then
the results of converting the XML fields to JSON.

`curl --insecure 
--header "Content-Type: application/xml" 
--data-binary @body.xml 
--request POST
https://localhost:9090/convert`


### /xml2json
This is the production endpoint; it accepts XML data,
converts it to JSON (as what `/convert` or `/parsifal`
would do), and sends it to the endpoint configured with
`--destination`.

## Flags
### --servicename *`service`*
The name of the local service/system. Not well tested. 
Not using it binds to all local services.

### --port *`port`*
Local port to listen on. Default is `9090`.

### --certfile *`filename.crt`*
Certificate file for https support in the miniservice.

### --keyfile *`filename.key`*
Private key for https support in the miniservice.

### --insecure
Connect to the destination service (which may be an
HTTPS service) without worrying about the validity of
the remote destination certificate. **This is to 
support testing only; do not use this in production**.

### --debug
Enable debugging code and messages

### --verbose
Enable additional messages about input and output

### --quiet
Suppress most log output to `STDERR`. 
*Log messages are written to the logfile 
regardless of this flag*. This flag ***only***
affects writes to `STDERR`.

### --destination
The /xml2json endpoint will send the
converted JSON information on to the specified destination,
which must be complete.  'localhost' is a **special** value
that sends the data to the service's /reflect endpoint.

