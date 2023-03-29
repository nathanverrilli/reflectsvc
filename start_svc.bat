@echo on
reflectsvc ^
   --keyfile host.key ^
   --certfile host.crt ^
   --debug ^
   --verbose ^
   --insecure ^
   --destination localhost ^
   --header-key Authorization
   --header-value "bearer ***DUMMYTOKEN***"
   --header-key Content-Frog
   --header-value "poisoned dart frog"
