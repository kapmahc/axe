#!/bin/sh
rm -r vendor
govendor init
govendor fetch golang.org/x/crypto/bcrypt
govendor fetch github.com/gin-gonic/gin
govendor fetch github.com/SermoDigital/jose/jwt
govendor fetch github.com/SermoDigital/jose/jws
govendor fetch github.com/SermoDigital/jose/crypto
govendor fetch github.com/streadway/amqp
govendor fetch github.com/google/uuid
govendor fetch github.com/spf13/viper
govendor fetch github.com/go-pg/migrations
govendor fetch github.com/urfave/cli
govendor fetch golang.org/x/text/language
govendor fetch github.com/facebookgo/inject
govendor fetch github.com/garyburd/redigo/redis
govendor fetch github.com/go-ini/ini
govendor fetch github.com/sirupsen/logrus
govendor fetch github.com/sirupsen/logrus/hooks/syslog
govendor fetch github.com/BurntSushi/toml
govendor fetch gopkg.in/gomail.v2
govendor fetch github.com/gin-contrib/sessions
govendor fetch github.com/aws/aws-sdk-go/aws/session
govendor fetch github.com/aws/aws-sdk-go/service/s3
govendor fetch gopkg.in/russross/blackfriday.v2
