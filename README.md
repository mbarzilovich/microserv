# microserv
Test task for LZD

#### System requirements:
1. Docker version 1.12+
2. Docker-compose version 1.9
3. GoLang version go1.7.4

#### Port mapping
Entry point for sending POST requests is listening on local port 4431 (HTTPS schema)
Bare receiver application is listening on local port 8082
Message broker admin UI avalilable on local port 8161 (Username: `admin`, password: `admin`)
Bare talkative applications is  available on local port 8092
Web server proxy for talkative application listens on local port 8091

#### Usage
* Use `microserv` wrapper to start/stop project

        ./microserv up
        ./microserv down
* In a separate console window run the loop which will produce messages

        i=1 ; while true ; do i=$((i+1)) ; curl -k -d 'text=foo'$i https://localhost:4431/ ; done

* Navigate to `localhost:8091` in you browser 
* See messages on the web page

