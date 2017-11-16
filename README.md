# AXE

A complete open source e-commerce solution by Go.

## Usage

```
go get -u github.com/kardianos/govendor
go get -d -u github.com/kapmahc/axe
cd $GOPATH/src/github.com/kapmahc/axe
govendor sync
npm install
make
```

## Atom plugins

- go-plus
- git-plus
- file-icons
- atom-beautify
- language-babel

## Notes

- Create database

```
psql -U postgres
CREATE DATABASE db-name WITH ENCODING = 'UTF8';
CREATE USER user-name WITH PASSWORD 'change-me';
GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
```

- ueditor

  ```
  cd node_modules/ueditor
  npm install grunt-cli -g
  npm install
  grunt
  ```

- Chrome browser: F12 => Console settings => Log XMLHTTPRequests

- Rabbitmq Management Plugin(<http://localhost:15612>)

  ```
  rabbitmq-plugins enable rabbitmq_management
  rabbitmqctl change_password guest change-me
  rabbitmqctl add_user who-am-i change-me
  rabbitmqctl set_user_tags who-am-i administrator
  rabbitmqctl list_vhosts
  rabbitmqctl add_vhost v-host
  rabbitmqctl set_permissions -p v-host who-am-i ".*" ".*" ".*"
  ```

- "RPC failed; HTTP 301 curl 22 The requested URL returned error: 301"

  ```
  git config --global http.https://gopkg.in.followRedirects true
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

- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)

- [favicon.ico](http://icoconvert.com/)

- [govendor](https://github.com/kardianos/govendor)

- [bootstrap](http://getbootstrap.com/docs/4.0/getting-started/introduction/)

- [AdminLTE](https://github.com/almasaeed2010/AdminLTE)

- [smver](http://semver.org/)
