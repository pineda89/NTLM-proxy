# NTLM-proxy

This is born to call a NTLM API without fighting the NTLM. the idea is to have this service deployed (perhaps on the same machine) and use this proxy for communication with NTLM.

## NTLM-proxy Settings

### command line args:
- port
- ntlm_url
- ntlm_username
- ntlm_password
- basicauth_username
- basicauth_password

### environment vars:
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
$ NTLM-proxy.exe -port 8080 -ntlm_url=http://NTLM-server:7047 -ntlm_username=username -ntlm_password=password -basicauth_username=clientuser -basicauth_password=clientpass
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

## Deploy as a windows service

```
$ sc.exe create NTLM-proxy binPath= "C:\filepath\NTLM-proxy.exe -port 8081 -ntlm_url=http://NTLM-server:7047 -ntlm_username=username -ntlm_password=password -basicauth_username=clientuser -basicauth_password=clientpass"
$ sc.exe start NTLM-proxy 
```