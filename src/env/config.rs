use std::io::{BufWriter, Write, Read};
use std::os::unix::fs::OpenOptionsExt;
use std::fs::{File, OpenOptions};
use std::path::Path;
use postgres;
use redis;
use amqp;
use base64;
use toml;
use log;
use rocket::config::{self, Environment};
use super::errors::Result;
use super::utils;

pub const CACHE_PREFIX: &'static str = "cache://";
pub const TASK_PREFIX: &'static str = "task://";


#[derive(Serialize, Deserialize, Debug)]
pub struct Config {
    env: String,
    workers: u16,
    pub(crate) database: PostgreSQL,
    pub(crate) redis: Redis,
    pub(crate) rabbitmq: RabbitMQ,
    pub(crate) http: HTTP,
}


impl From<Config> for config::ConfigBuilder {
    fn from(c: Config) -> Self {
        let env = match c.env.parse::<Environment>() {
            Ok(v) => v,
            Err(_) => {
                log::error!("bad enviroment {}", c.env);
                config::Environment::Development
            }
        };

        let lev = if env.is_prod() {
            config::LoggingLevel::Normal
        } else {
            config::LoggingLevel::Debug
        };

        let views = match Path::new("themes")
            .join(c.http.theme)
            .join("views")
            .to_str() {
            Some(v) => v.to_string(),
            None => "views".to_string(),
        };
        return config::Config::build(env)
            .address("localhost")
            .port(c.http.port)
            .workers(c.workers)
            .limits(
                config::Limits::new()
                    .limit("forms", c.http.limits * (1 << 20))
                    .limit("json", c.http.limits * (1 << 20)),
            )
            .secret_key(c.http.secret)
            .log_level(lev)
            .extra("database", c.database.url())
            .extra("redis", c.redis.url())
            .extra("rabbitmq", c.rabbitmq.url())
            .extra("template_dir", views);
    }
}

impl Config {
    pub fn load(name: &str) -> Result<Config> {
        log::info!("read config from file {:?}", name);
        let mut file = try!(File::open(&name));
        let mut data = String::new();
        try!(file.read_to_string(&mut data));
        let cfg = try!(toml::from_str(&data));
        return Ok(cfg);
    }

    pub fn new() -> Config {
        return Config {
            env: "development".to_string(),
            workers: 6,
            database: PostgreSQL::new(),
            redis: Redis::new(),
            rabbitmq: RabbitMQ::new(),
            http: HTTP::new(),
        };
    }
    pub fn write(&self, name: String) -> Result<()> {
        log::info!("generate file {}", name);
        let fd = try!(
            OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o600)
                .open(name)
        );
        let mut wrt = BufWriter::new(&fd);
        try!(wrt.write_all(&try!(toml::to_vec(self))));
        Ok(())
    }
}


#[derive(Serialize, Deserialize, Debug)]
pub struct PostgreSQL {
    host: String,
    port: u16,
    user: String,
    password: String,
    name: String,
}

impl PostgreSQL {
    pub fn new() -> PostgreSQL {
        let name = env!("CARGO_PKG_NAME");
        PostgreSQL {
            host: "localhost".to_string(),
            port: 5432,
            user: "postgres".to_string(),
            password: "".to_string(),
            name: name.to_string(),
        }
    }

    pub fn host(&mut self, host: String) -> &mut PostgreSQL {
        self.host = host;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut PostgreSQL {
        self.port = port;
        self
    }
    pub fn name(&mut self, name: String) -> &mut PostgreSQL {
        self.name = name;
        self
    }
    pub fn user(&mut self, user: String) -> &mut PostgreSQL {
        self.user = user;
        self
    }
    pub fn password(&mut self, password: String) -> &mut PostgreSQL {
        self.password = password;
        self
    }
    pub fn url(&self) -> String {
        format!(
            "postgres://{}:{}@{}:{}/{}",
            self.user,
            self.password,
            self.host,
            self.port,
            self.name
        )
    }


    pub fn open(&self) -> Result<postgres::Connection> {
        log::info!("open database {}", self.url());
        let con = try!(postgres::Connection::connect(
            self.url(),
            postgres::TlsMode::None,
        ));
        return Ok(con);
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Redis {
    host: String,
    port: u16,
    db: i64,
}

impl Redis {
    pub fn new() -> Redis {
        Redis {
            host: "localhost".to_string(),
            port: 6379,
            db: 0,
        }
    }
    pub fn host(&mut self, host: String) -> &mut Redis {
        self.host = host;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut Redis {
        self.port = port;
        self
    }
    pub fn db(&mut self, db: i64) -> &mut Redis {
        self.db = db;
        self
    }
    pub fn url(&self) -> String {
        format!("redis://{}:{}/{}", self.host, self.port, self.db)
    }
    pub fn open(&self) -> Result<redis::Client> {
        log::info!("open {}", self.url());
        let con = try!(redis::Client::open(redis::ConnectionInfo {
            addr: Box::new(redis::ConnectionAddr::Tcp(self.host.to_string(), self.port)),
            db: self.db,
            passwd: None,
        }));
        return Ok(con);
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct RabbitMQ {
    host: String,
    port: u16,
    user: String,
    password: String,
    #[serde(rename = "virtual")]
    _virtual: String,
}

impl RabbitMQ {
    pub fn new() -> RabbitMQ {
        let name = env!("CARGO_PKG_NAME");
        RabbitMQ {
            host: "localhost".to_string(),
            port: 5432,
            user: "postgres".to_string(),
            password: "".to_string(),
            _virtual: name.to_string(),
        }
    }

    pub fn host(&mut self, host: String) -> &mut RabbitMQ {
        self.host = host;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut RabbitMQ {
        self.port = port;
        self
    }
    pub fn name(&mut self, _virtual: String) -> &mut RabbitMQ {
        self._virtual = _virtual;
        self
    }
    pub fn user(&mut self, user: String) -> &mut RabbitMQ {
        self.user = user;
        self
    }
    pub fn password(&mut self, password: String) -> &mut RabbitMQ {
        self.password = password;
        self
    }
    pub fn url(&self) -> String {
        format!(
            "amqp://{}@{}:{}/{}",
            self.user,
            self.host,
            self.port,
            self._virtual
        )
    }
    pub fn open(&self) -> Result<amqp::Session> {
        log::info!("open rabbitmq {}", self.url());
        let url = format!(
            "amqp://{}:{}@{}:{}/{}",
            self.user,
            self.password,
            self.host,
            self.port,
            self._virtual
        );
        let con = try!(amqp::Session::open_url(&url));
        return Ok(con);
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct HTTP {
    pub(crate) port: u16,
    pub(crate) name: String,
    pub(crate) theme: String,
    secret: String,
    limits: u64,
}

impl HTTP {
    pub fn new() -> HTTP {
        HTTP {
            port: 8080,
            // a 256-bit base64 encoded string (44 characters)
            secret: utils::random(256 / 8),
            theme: "moon".to_string(),
            name: "www.change-me.com".to_string(),
            limits: 20,
        }
    }
    pub fn secret(&mut self, secret: String) -> &mut HTTP {
        self.secret = secret;
        self
    }
    pub fn secret_key(&self) -> Result<Vec<u8>> {
        let key = try!(base64::decode(&self.secret));
        return Ok(key);
    }
    pub fn name(&mut self, name: String) -> &mut HTTP {
        self.name = name;
        self
    }
    pub fn theme(&mut self, theme: String) -> &mut HTTP {
        self.theme = theme;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut HTTP {
        self.port = port;
        self
    }
    pub fn limits(&mut self, limits: u64) -> &mut HTTP {
        self.limits = limits;
        self
    }
}
