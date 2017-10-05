# axe

A complete open source e-commerce solution.

## Install jdk

```
curl -s "https://get.sdkman.io" | zsh
sdk install java
sdk install gradle
sdk default java 8u144-zulu
```

## Build

```
git clone https://github.com/kapmahc/axe.git
cd axe
npm install
gradle build
```

### Install ueditor

```
npm install -g grunt-cli
git clone https://github.com/fex-team/ueditor.git node_modules/ueditor
cd node_modules/ueditor
git checkout v1.4.3.3
npm install
grunt
```

## Notes

- Create database

```
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)

  ```
  rabbitmq-plugins enable rabbitmq_management
  rabbitmqctl change_password guest change-me
  rabbitmqctl add_user who-am-i change-me
  rabbitmqctl set_user_tags who-am-i administrator
  rabbitmqctl list_vhosts
  rabbitmqctl add_vhost /v-host
  rabbitmqctl set_permissions -p /v-host who-am-i ".*" ".*" ".*"
  ```

- 'Peer authentication failed for user', open file "/etc/postgresql/9.5/main/pg_hba.conf" change line:

  ```
  local   all             all                                     peer  
  TO:
  local   all             all                                     md5
  ```

- Generate openssl certs

  ```
  openssl genrsa -out www.change-me.com.key 2048
  openssl req -new -x509 -key www.change-me.com.key -out www.change-me.com.crt -days 3650 # Common Name:*.change-me.com
  ```

## Documents

- [jce](http://www.oracle.com/technetwork/java/javase/downloads/jce8-download-2133166.html)
- [spring-boot](https://docs.spring.io/spring-boot/docs/2.0.0.M4/reference/html/index.html)
- [application.properties](https://docs.spring.io/spring-boot/docs/2.0.0.M4/reference/html/common-application-properties.html)
- [spring-security](https://docs.spring.io/spring-security/site/docs/5.0.0.M3/reference/htmlsingle/)
- [thymeleaf](http://www.thymeleaf.org/doc/tutorials/3.0/usingthymeleaf.html)
- [thymeleaf-layout](https://ultraq.github.io/thymeleaf-layout-dialect/Installation.html)
- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)
- [favicon.ico](http://icoconvert.com/)
