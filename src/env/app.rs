use std::fs::{read_dir, File, create_dir_all, OpenOptions};
use std::os::unix::fs::OpenOptionsExt;
use std::path::{Path, PathBuf};
use std::io::Read;
use std::env;
use std::result;
use time;
use rocket;
use mustache;
use postgres;
use log;
use redis::{self, Commands};
use super::errors::{Error, Result};
use super::config;
use super::super::nut;
use super::super::forum;
use super::super::survey;

pub struct App {
    name: String,
}


impl App {
    pub fn new(name: String) -> App {
        App { name: name }
    }

    // ---------- cache --------

    pub fn cache_list(&self) -> Result<()> {
        let con = try!(try!(try!(config::Config::load(&self.name)).redis.open()).get_connection());
        let keys: Vec<String> = try!(con.keys(format!("{}*", config::CACHE_PREFIX)));
        macro_rules! fmt {() => ("{:<32}{}")};
        println!(fmt!(), "KEY", "TTL(S)");
        for k in keys {
            let t: i64 = try!(redis::cmd("TTL").arg(&k).query(&con));
            println!(fmt!(), k, t);
        }
        return Ok(());
    }

    pub fn cache_clear(&self) -> Result<()> {
        let con = try!(try!(try!(config::Config::load(&self.name)).redis.open()).get_connection());
        let keys: Vec<String> = try!(con.keys(format!("{}*", config::CACHE_PREFIX)));
        let _: () = try!(con.del(keys));
        return Ok(());
    }


    // ---------- database --------

    pub fn database_migrate(&self) -> Result<()> {
        let root = self.get_database_migrations_root();
        let con = try!(self.get_database_connenction());
        let db = try!(con.transaction());
        let mut files = try!(try!(read_dir(&root)).collect::<result::Result<Vec<_>, _>>());
        files.sort_by_key(|d| d.path());
        for entry in files {
            let mig = entry.path();
            match mig.file_name() {
                None => (),
                Some(name) => {
                    match name.to_str() {
                        None => (),
                        Some(mig) => {
                            match mig.find('-') {
                                None => (),
                                Some(idx) => {
                                    log::info!("Find migration {}", mig);
                                    let version = &mig[..idx];
                                    let version = try!(version.parse::<i64>());
                                    let name = &mig[idx + 1..];
                                    let count: i64 = try!(db.query(
                                        "SELECT COUNT(*) FROM schema_migrations
                                        WHERE version = $1 AND name = $2",
                                        &[&version, &name],
                                    )).get(0)
                                        .get(0);
                                    if count == 0 {
                                        let mut file =
                                            try!(File::open(&root.join(mig).join("up.sql")));
                                        let mut buf = String::new();
                                        try!(file.read_to_string(&mut buf));
                                        try!(db.batch_execute(&buf));
                                        try!(db.execute(
                                            "INSERT INTO
                                            schema_migrations(version, name)
                                            VALUES($1, $2)",
                                            &[&version, &name],
                                        ));
                                    } else {
                                        log::info!("ingnore")
                                    }
                                }
                            };
                        }
                    }
                }
            }
        }
        try!(db.commit());
        return Ok(());
    }

    pub fn database_rollback(&self) -> Result<()> {
        let con = try!(self.get_database_connenction());
        let db = try!(con.transaction());
        for row in &try!(db.query(
            "SELECT version, name FROM schema_migrations ORDER BY id DESC LIMIT 1",
            &[],
        ))
        {
            let version: i64 = row.get(0);
            let name: String = row.get(1);
            log::info!("Rollback {} {}", version, name);
            let root = self.get_database_migrations_root();
            let mut file = try!(File::open(
                &root.join(format!("{}-{}", version, name)).join("down.sql"),
            ));
            let mut buf = String::new();
            try!(file.read_to_string(&mut buf));
            try!(db.batch_execute(&buf));
            try!(db.execute(
                "DELETE FROM schema_migrations WHERE version = $1 AND name = $2",
                &[&version, &name],
            ));
        }
        try!(db.commit());
        return Ok(());
    }

    pub fn database_show(&self) -> Result<()> {
        let db = try!(self.get_database_connenction());
        macro_rules! fmt {() => ("{:<16}{:<24}{}")};
        println!(fmt!(), "VERSION", "NAME", "CREATED AT");
        for row in &try!(db.query(
            "SELECT version, name, created_at FROM schema_migrations",
            &[],
        ))
        {
            let version: i64 = row.get(0);
            let name: String = row.get(1);
            let created_at: time::Timespec = row.get(2);
            println!(fmt!(), version, name, time::at(created_at).rfc822());
        }
        return Ok(());
    }

    fn get_database_migrations_root(&self) -> PathBuf {
        return Path::new("db").join("migrations");
    }

    fn get_database_connenction(&self) -> Result<postgres::Connection> {
        let db = try!(try!(config::Config::load(&self.name)).database.open());
        try!(db.execute(
            "CREATE TABLE IF NOT EXISTS schema_migrations (
                    id SERIAL PRIMARY KEY,
                    version BIGINT NOT NULL,
                    name VARCHAR(255) NOT NULL,
                    created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
                  )",
            &[],
        ));
        try!(db.execute(
            "CREATE UNIQUE INDEX IF NOT EXISTS idx_schema_migrations_version_name
                ON schema_migrations(version, name)",
            &[],
        ));
        return Ok(db);
    }

    // ------- generate --------

    pub fn generate_nginx_conf(&self, ssl: bool) -> Result<()> {
        let file = Path::new("tmp").join("nginx.conf");
        match file.parent() {
            None => {}
            Some(d) => try!(create_dir_all(d)),
        }
        let cfg = try!(config::Config::load(&self.name));
        let scheme = if ssl { "https" } else { "http" };

        log::info!(
            "generate file {:?} for {}://{}",
            file.as_os_str(),
            scheme,
            cfg.http.name
        );
        let mut tpf = try!(File::open(Path::new("templates").join("nginx.conf")));
        let mut buf = String::new();
        try!(tpf.read_to_string(&mut buf));
        let tpl = try!(mustache::compile_str(&buf));
        let data = try!(mustache::MapBuilder::new().insert("port", &cfg.http.port))
            .insert_bool("ssl", ssl)
            .insert_str("theme", cfg.http.theme)
            .insert_str("name", cfg.http.name)
            .insert_str("root", try!(env::current_dir()).display())
            .build();
        let mut tpd = try!(
            OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o644)
                .open(file)
        );
        try!(tpl.render_data(&mut tpd, &data));
        return Ok(());
    }

    pub fn generate_locale(&self, name: String) -> Result<()> {
        let mut _fn = Path::new("locales").join(name);
        _fn.set_extension("yaml");
        return self.create_empty_file(_fn);
    }

    pub fn generate_migration(&self, name: String) -> Result<()> {
        let fmt = "%Y%m%d%H%M%S".to_string() + "-" + &name;
        let tsf = try!(time::strftime(&fmt, &time::now_utc()));
        let root = Path::new("db").join("migrations").join(tsf);
        let ext = "sql";

        let mut up = root.join("up");
        up.set_extension(ext);
        try!(self.create_empty_file(up));

        let mut down = root.join("down");
        down.set_extension(ext);
        try!(self.create_empty_file(down));

        Ok(())
    }

    // --------- start server --------

    pub fn start_server(&self) -> Result<()> {
        let cfg = try!(
            rocket::config::ConfigBuilder::from(try!(config::Config::load(&self.name)))
                .finalize()
        );
        let err = rocket::custom(cfg, false)
            .manage(nut::middlewares::Status::new())
            .mount("/", nut::routes())
            .mount("/forum", forum::routes())
            .mount("/survey", survey::routes())
            .catch(errors![nut::not_found])
            .launch();
        return Err(Error::from(err));
    }

    // ------------

    fn create_empty_file(&self, name: PathBuf) -> Result<()> {
        log::info!("generate file {:?}", name.as_os_str());
        match name.parent() {
            None => {}
            Some(d) => try!(create_dir_all(d)),
        }
        try!(
            OpenOptions::new()
                .write(true)
                .create_new(true)
                .mode(0o600)
                .open(name)
        );
        Ok(())
    }
}
