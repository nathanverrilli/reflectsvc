@echo on

reflectsvc ^
  --keyfile host.key ^
  --certfile host.crt ^
  --destination localhost ^
   --debug ^
   --verbose ^
   --insecure ^
   --port 9090