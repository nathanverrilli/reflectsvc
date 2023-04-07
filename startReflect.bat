@echo off

reflectsvc ^
  --destination localhost ^
  --fieldnames fieldnames.csv ^
  --debug ^
  --verbose ^
  --insecure ^
  --port 9090
