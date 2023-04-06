# --keyfile host.key ^
# --certfile host.crt ^

@echo on
reflectsvc ^
   --debug ^
   --verbose ^
   --insecure ^
   --destination localhost ^
   --header-key Authorization ^
   --header-value 'bearer ***DUMMYTOKEN***' ^
   --header-key Content-Frog ^
   --header-value 'poisoned dart frog'
