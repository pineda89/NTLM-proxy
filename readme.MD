# NTLM-proxy

This is born to call the navision API without fighting the NTLM. the idea is to have this service deployed (perhaps on the same machine) and use this proxy for communication with navision.

## NTLM-proxy Settings

environment vars:
- SERVER_PORT (default is 8080)
- DEFAULT_NTLM_URL (when NTLM_URL not informed in request, DEFAULT_NTLM_URL is used)
- DEFAULT_NTLM_USERNAME (when NTLM_USERNAME not informed in request, DEFAULT_NTLM_USERNAME is used)
- DEFAULT_NTLM_PASSWORD (when NTLM_PASSWORD not informed in request, DEFAULT_NTLM_PASSWORD is used)
- BASICAUTH_USERNAME (blank is no auth required)
- BASICAUTH_PASSWORD (blank is no auth required)

## Client usage

http request headers:
- NTLM_URL (optional, when not informed DEFAULT_NTLM_URL is used)
- NTLM_USERNAME (optional, when not informed DEFAULT_NTLM_USERNAME is used)
- NTLM_PASSWORD (optional, when not informed DEFAULT_NTLM_PASSWORD is used)

## Example

Start the server
```
$ set DEFAULT_NTLM_URL=http://NTLM-server:7047
$ set DEFAULT_NTLM_USERNAME=username
$ set DEFAULT_NTLM_PASSWORD=password
$ set BASICAUTH_USERNAME=clientuser
$ set BASICAUTH_PASSWORD=clientpass
$ NTLM-proxy.exe
```

Call the API
```
$ curl -u clientuser:clientpass http://localhost:8080/DynamicsNAV/WS/Services
```

## Using docker

```
$ docker build -t ntlm-proxy .
$ docker run -d -p 8080:8080 -e DEFAULT_NTLM_URL=http://NTLM-server:7047 -e DEFAULT_NTLM_USERNAME=username -e DEFAULT_NTLM_PASSWORD=password --name ntlm-proxy ntlm-proxy
```

