@echo off

reflectsvc ^
  --keyfile host.key ^
  --certfile host.crt ^
  --destination localhost ^
  --fieldnames fieldnames.csv ^
  --debug ^
  --verbose ^
  --insecure ^
  --port 9090