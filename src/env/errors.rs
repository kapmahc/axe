use std::io;
use std::fmt;
use std::error;
use std::result;
use std::num;
use toml;
use time;
use postgres;
use redis;
use base64;
use amqp;
use rocket;
use mustache;
use diesel;
use r2d2;

pub type Result<T> = result::Result<T, Error>;

#[derive(Debug)]
pub enum Error {
    Io(io::Error),
    NumParseIntError(num::ParseIntError),
    TomlSer(toml::ser::Error),
    TomlDe(toml::de::Error),
    TimeParse(time::ParseError),
    Base64Decode(base64::DecodeError),
    AMQP(amqp::AMQPError),
    Redis(redis::RedisError),
    Postgres(postgres::Error),
    RocketConfig(rocket::config::ConfigError),
    RocketLaunchError(rocket::error::LaunchError),
    MustacheError(mustache::Error),
    MustacheEncoderError(mustache::encoder::Error),
    DieselMigrationError(diesel::migrations::MigrationError),
    DieselRunMigrationsError(diesel::migrations::RunMigrationsError),
    DieselConnectionError(diesel::ConnectionError),
    R2d2InitializationError(r2d2::InitializationError),
}

impl fmt::Display for Error {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match *self {
            Error::Io(ref err) => err.fmt(f),
            Error::NumParseIntError(ref err) => err.fmt(f),
            Error::TomlSer(ref err) => err.fmt(f),
            Error::TomlDe(ref err) => err.fmt(f),
            Error::TimeParse(ref err) => err.fmt(f),
            Error::Base64Decode(ref err) => err.fmt(f),
            Error::AMQP(ref err) => err.fmt(f),
            Error::Redis(ref err) => err.fmt(f),
            Error::Postgres(ref err) => err.fmt(f),
            Error::RocketConfig(ref err) => err.fmt(f),
            Error::RocketLaunchError(ref err) => err.fmt(f),
            Error::MustacheError(ref err) => err.fmt(f),
            Error::MustacheEncoderError(ref err) => err.fmt(f),
            Error::DieselMigrationError(ref err) => err.fmt(f),
            Error::DieselRunMigrationsError(ref err) => err.fmt(f),
            Error::DieselConnectionError(ref err) => err.fmt(f),
            Error::R2d2InitializationError(ref err) => err.fmt(f),
        }
    }
}

impl error::Error for Error {
    fn description(&self) -> &str {
        match *self {
            Error::Io(ref err) => err.description(),
            Error::NumParseIntError(ref err) => err.description(),
            Error::TomlSer(ref err) => err.description(),
            Error::TomlDe(ref err) => err.description(),
            Error::TimeParse(ref err) => err.description(),
            Error::Base64Decode(ref err) => err.description(),
            Error::AMQP(ref err) => err.description(),
            Error::Redis(ref err) => err.description(),
            Error::Postgres(ref err) => err.description(),
            Error::RocketConfig(ref err) => err.description(),
            Error::RocketLaunchError(ref err) => err.description(),
            Error::MustacheError(ref err) => err.description(),
            Error::MustacheEncoderError(ref err) => err.description(),
            Error::DieselMigrationError(ref err) => err.description(),
            Error::DieselRunMigrationsError(ref err) => err.description(),
            Error::DieselConnectionError(ref err) => err.description(),
            Error::R2d2InitializationError(ref err) => err.description(),
        }
    }

    fn cause(&self) -> Option<&error::Error> {
        match *self {
            Error::Io(ref err) => Some(err),
            Error::NumParseIntError(ref err) => Some(err),
            Error::TomlSer(ref err) => Some(err),
            Error::TomlDe(ref err) => Some(err),
            Error::TimeParse(ref err) => Some(err),
            Error::Base64Decode(ref err) => Some(err),
            Error::AMQP(ref err) => Some(err),
            Error::Redis(ref err) => Some(err),
            Error::Postgres(ref err) => Some(err),
            Error::RocketConfig(ref err) => Some(err),
            Error::RocketLaunchError(ref err) => Some(err),
            Error::MustacheError(ref err) => Some(err),
            Error::MustacheEncoderError(ref err) => Some(err),
            Error::DieselMigrationError(ref err) => Some(err),
            Error::DieselRunMigrationsError(ref err) => Some(err),
            Error::DieselConnectionError(ref err) => Some(err),
            Error::R2d2InitializationError(ref err) => Some(err),
        }
    }
}

impl From<io::Error> for Error {
    fn from(err: io::Error) -> Error {
        Error::Io(err)
    }
}

impl From<num::ParseIntError> for Error {
    fn from(err: num::ParseIntError) -> Error {
        Error::NumParseIntError(err)
    }
}

impl From<toml::ser::Error> for Error {
    fn from(err: toml::ser::Error) -> Error {
        Error::TomlSer(err)
    }
}

impl From<toml::de::Error> for Error {
    fn from(err: toml::de::Error) -> Error {
        Error::TomlDe(err)
    }
}

impl From<time::ParseError> for Error {
    fn from(err: time::ParseError) -> Error {
        Error::TimeParse(err)
    }
}


impl From<postgres::Error> for Error {
    fn from(err: postgres::Error) -> Error {
        Error::Postgres(err)
    }
}

impl From<redis::RedisError> for Error {
    fn from(err: redis::RedisError) -> Error {
        Error::Redis(err)
    }
}

impl From<base64::DecodeError> for Error {
    fn from(err: base64::DecodeError) -> Error {
        Error::Base64Decode(err)
    }
}

impl From<amqp::AMQPError> for Error {
    fn from(err: amqp::AMQPError) -> Error {
        Error::AMQP(err)
    }
}

impl From<rocket::config::ConfigError> for Error {
    fn from(err: rocket::config::ConfigError) -> Error {
        Error::RocketConfig(err)
    }
}

impl From<rocket::error::LaunchError> for Error {
    fn from(err: rocket::error::LaunchError) -> Error {
        Error::RocketLaunchError(err)
    }
}

impl From<mustache::Error> for Error {
    fn from(err: mustache::Error) -> Error {
        Error::MustacheError(err)
    }
}

impl From<mustache::encoder::Error> for Error {
    fn from(err: mustache::encoder::Error) -> Error {
        Error::MustacheEncoderError(err)
    }
}

impl From<diesel::migrations::MigrationError> for Error {
    fn from(err: diesel::migrations::MigrationError) -> Error {
        Error::DieselMigrationError(err)
    }
}

impl From<diesel::migrations::RunMigrationsError> for Error {
    fn from(err: diesel::migrations::RunMigrationsError) -> Error {
        Error::DieselRunMigrationsError(err)
    }
}

impl From<diesel::ConnectionError> for Error {
    fn from(err: diesel::ConnectionError) -> Error {
        Error::DieselConnectionError(err)
    }
}

impl From<r2d2::InitializationError> for Error {
    fn from(err: r2d2::InitializationError) -> Error {
        Error::R2d2InitializationError(err)
    }
}
