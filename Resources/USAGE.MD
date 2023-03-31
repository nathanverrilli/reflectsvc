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
The `/xml2json` endpoint does the conversion, and sends
the data to the `--destination` endpoint.

**Please note that `/parsifal` endpoint is deprecated.
Please use the `/convert` endpoint instead.

<pre>
curl --verbose --insecure ^
  --header "Content-Type: application/xml" ^ 
  --data-binary @body.xml ^
  --request POST https://localhost:9090/convert`
</pre>


### /xml2json
This is the production endpoint; it accepts XML data,
converts it to JSON (as what `/convert` or `/parsifal`
would do), and sends it to the endpoint configured with
`--destination` as a proxy. It forwards some headers as
part of the proxy.
#### Forwarded Headers
1. Authorization
1. Accepts

#### Returned Headers



## Flags

### --servicename *`service`*
The name of the local service/system. Not well tested. 
Not using it binds to all local services.

### --port *`port`*
Local port to listen on. Default is `9090`.

### --certfile *`filename.crt`*
Certificate file for https support in the miniservice.
If this flag is missing, the service defaults to HTTP.

### --keyfile *`filename.key`*
Private key for https support in the miniservice.
If this flag is missing, the service defaults to HTTP.

### --insecure
Connect to the destination service (which may be an
HTTPS service) without worrying about the validity of
the remote destination certificate. **This is to 
support testing only; do not use this in production**.

This flag affects the `/xml2json` endpoint (only) when
it proxies the request forward to the `--destination`
endpoint.

### --debug
Enable debugging code and messages. If
`--quiet` is enabled, then this debugging output
is still generated, but goes only to the log.

### --verbose
Enable additional messages about input and output. If
`--quiet` is enabled, then these additional messages are
still generated, but go only into the log.

### --quiet
Suppress most log output to `STDERR`. 
*Log messages are written to the logfile 
regardless of this flag*. This flag ***only***
affects writes to `STDERR`.

### --destination *`endpoint`*
The /xml2json endpoint will send the
converted JSON information on to the specified destination,
which must be complete.  'localhost' is a **special** value
that sends the data to the service's `/reflect` endpoint with
the port, so with the default `--port 9090` is the equivalent
of specifying `--destination https://localhost:9090/reflect`.

This flag affects the `/xml2json` endpoint.

### --fieldNames *`filename`*

The `/xml2json` endpoint can of incoming XML field names to outgoing JSON
field names as part of the xml2json endpoint. These to / from strings
are held in the file specifed by `--fieldNames <file>`. <file> should
be a plain unicode file. Lines beginning with `#` are ignored as 
comments, and all comments must begin on a new line.
Empty lines are ignored. Field name replacements are specified as:  
&nbsp;&nbsp;&nbsp;&nbsp;&lsqb;`From XML Field Name`&rsqb;&lsqb;`toJsonFieldName`&rsqb;.  
Within the square brackets, both case and spaces are significant.
No space is permitted between the two tokens.

This flag affects the `/xml2json`, `/convert`, and the
`/parsifal` endpoints.

### --omitEmpty

When converting XML to JSON and sending it, if an 
XML field has a name but no value (or an empty value)
then omit the empty field in the JSON output. Omitted
JSON values are logged.

This flag affects the `/xml2json`, `/convert`, and the
`/parsifal` endpoints.

