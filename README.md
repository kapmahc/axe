# AXE

A complete open source e-commerce solution.

## Usage

```
  mvn package
  java -jar target/axe-*.jar db status config.yml
  java -jar target/axe-*.jar db migrate config.yml  
  java -jar target/axe-*.jar db rollback config.yml
  java -jar target/axe-*.jar server config.yml
  java -jar target/axe-*.jar check config.yml
```

## Notes

- Generate a random key

  ```
  openssl rand -base64 32
  ```

- ~/.npmrc

  ```
  prefix=${HOME}/.npm-packages
  ```

- Install ueditor

  ```
  npm install -g grunt-cli  
  cd node_modules/ueditor  
  npm install
  grunt
  ```

- Create database

  ```
  psql -U postgres
  CREATE DATABASE db-name WITH ENCODING = 'UTF8';
  CREATE USER user-name WITH PASSWORD 'change-me';
  GRANT ALL PRIVILEGES ON DATABASE db-name TO user-name;
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

- [dropwizard](http://www.dropwizard.io/)
- [bootstrap](http://getbootstrap.com/docs/4.0/getting-started/introduction/)
- [AdminLTE](https://github.com/almasaeed2010/AdminLTE)
- [For gmail smtp](http://stackoverflow.com/questions/20337040/gmail-smtp-debug-error-please-log-in-via-your-web-browser)
