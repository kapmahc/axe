use std::io::{BufWriter, Write};
use std::os::unix::fs::OpenOptionsExt;
use std::fs::OpenOptions;
use std::boxed::Box;
use postgres;
use redis;
use amqp;
use base64;
use rand;
use toml;
use rand::Rng;
use super::errors::Result;

#[derive(Serialize, Deserialize, Debug)]
pub struct Config {
    env: &'static str,
    database: PostgreSQL,
    redis: Redis,
    rabbitmq: RabbitMQ,
    http: HTTP,
}

impl Config {
    pub fn new() -> Config {
        Config {
            env: "development",
            database: PostgreSQL::new(),
            redis: Redis::new(),
            rabbitmq: RabbitMQ::new(),
            http: HTTP::new(),
        }
    }
    pub fn write(&self, name: String) -> Result<()> {
        info!("generate file {}", name);
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
    // pub fn read() -> Result<Config> {
    //     let cfg: Config = try!(toml::from_str(""));
    //     Ok(cfg)
    // }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct PostgreSQL {
    host: &'static str,
    port: u16,
    user: &'static str,
    password: &'static str,
    name: &'static str,
}

impl PostgreSQL {
    pub fn new() -> PostgreSQL {
        let name = env!("CARGO_PKG_NAME");
        PostgreSQL {
            host: "localhost",
            port: 5432,
            user: "postgres",
            password: "",
            name: name,
        }
    }

    pub fn host(&mut self, host: &'static str) -> &mut PostgreSQL {
        self.host = host;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut PostgreSQL {
        self.port = port;
        self
    }
    pub fn name(&mut self, name: &'static str) -> &mut PostgreSQL {
        self.name = name;
        self
    }
    pub fn user(&mut self, user: &'static str) -> &mut PostgreSQL {
        self.user = user;
        self
    }
    pub fn password(&mut self, password: &'static str) -> &mut PostgreSQL {
        self.password = password;
        self
    }
    pub fn url(&self) -> String {
        format!(
            "postgres://{}@{}:{}/{}",
            self.user,
            self.host,
            self.port,
            self.name
        )
    }
    pub fn open(&self) -> Result<postgres::Connection> {
        info!("open database {}", self.url());
        let con = try!(postgres::Connection::connect(
            format!(
                "postgres://{}:{}@{}:{}/{}",
                self.user,
                self.password,
                self.host,
                self.port,
                self.name
            ),
            postgres::TlsMode::None,
        ));
        return Ok(con);
    }
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Redis {
    host: &'static str,
    port: u16,
    db: i64,
}

impl Redis {
    pub fn new() -> Redis {
        Redis {
            host: "localhost",
            port: 6379,
            db: 0,
        }
    }
    pub fn host(&mut self, host: &'static str) -> &mut Redis {
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
    host: &'static str,
    port: u16,
    user: &'static str,
    password: &'static str,
    #[serde(rename = "virtual")]
    _virtual: &'static str,
}

impl RabbitMQ {
    pub fn new() -> RabbitMQ {
        let name = env!("CARGO_PKG_NAME");
        RabbitMQ {
            host: "localhost",
            port: 5432,
            user: "postgres",
            password: "",
            _virtual: name,
        }
    }

    pub fn host(&mut self, host: &'static str) -> &mut RabbitMQ {
        self.host = host;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut RabbitMQ {
        self.port = port;
        self
    }
    pub fn name(&mut self, _virtual: &'static str) -> &mut RabbitMQ {
        self._virtual = _virtual;
        self
    }
    pub fn user(&mut self, user: &'static str) -> &mut RabbitMQ {
        self.user = user;
        self
    }
    pub fn password(&mut self, password: &'static str) -> &mut RabbitMQ {
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
        info!("open rabbitmq {}", self.url());
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
    port: u16,
    name: &'static str,
    theme: &'static str,
    secret: String,
}

impl HTTP {
    pub fn new() -> HTTP {
        let name = env!("CARGO_PKG_NAME");
        let secret: String = rand::thread_rng().gen_iter::<char>().take(32).collect();
        HTTP {
            port: 8080,
            secret: base64::encode(&secret),
            theme: "bootstrap",
            name: name,
        }
    }
    pub fn secret(&mut self, secret: &'static str) -> &mut HTTP {
        self.secret = secret.to_string();
        self
    }
    pub fn secret_key(&self) -> Result<Vec<u8>> {
        let key = try!(base64::decode(&self.secret));
        return Ok(key);
    }
    pub fn name(&mut self, name: &'static str) -> &mut HTTP {
        self.name = name;
        self
    }
    pub fn theme(&mut self, theme: &'static str) -> &mut HTTP {
        self.theme = theme;
        self
    }
    pub fn port(&mut self, port: u16) -> &mut HTTP {
        self.port = port;
        self
    }
}
