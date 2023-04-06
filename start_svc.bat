# --keyfile host.key ^
# --certfile host.crt ^

@echo on
reflectsvc ^
   --fieldnames fieldnames.csv ^
   --debug ^
   --verbose ^
   --insecure ^
   --destination localhost
