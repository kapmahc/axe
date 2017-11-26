use std::fs::{File, create_dir_all, OpenOptions};
use std::os::unix::fs::OpenOptionsExt;
use std::path::{Path, PathBuf};
use std::io::Read;
use std::env;
use time;
use rocket;
use mustache;
use log;
use redis::{self, Commands};
use diesel::{migrations, pg, Connection};
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
        let root = try!(migrations::find_migrations_directory());
        log::info!("find migrations in {:?}", root);
        let db = try!(self.get_database_connenction());
        try!(migrations::run_pending_migrations(&db));
        return Ok(());
    }

    pub fn database_rollback(&self) -> Result<()> {
        let root = try!(migrations::find_migrations_directory());
        log::info!("find migrations in {:?}", root);
        let db = try!(self.get_database_connenction());
        try!(migrations::revert_latest_migration(&db));
        return Ok(());
    }

    pub fn database_show(&self) -> Result<()> {
        let root = try!(migrations::find_migrations_directory());
        log::info!("find migrations in {:?}", root);
        let db = try!(self.get_database_connenction());
        let items = try!(migrations::mark_migrations_in_directory(&db, &root));
        for (it, ok) in items {
            match it {
                None => (),
                Some(f) => {
                    match f.file_name() {
                        None => (),
                        Some(f) => {
                            match f.to_str() {
                                None => (),
                                Some(f) => {
                                    println!("{:<32} {}", f, ok);
                                    ()
                                }
                            }
                        }
                    }
                }
            };

        }




        return Ok(());
    }

    // fn get_database_migrations_root(&self) -> PathBuf {
    //     return Path::new("db").join("migrations");
    // }

    fn get_database_connenction(&self) -> Result<pg::PgConnection> {
        let db = try!(pg::PgConnection::establish(
            &try!(config::Config::load(&self.name)).database.url(),
        ));
        return Ok(db);
        // let db = try!(try!(config::Config::load(&self.name)).database.open());
        // try!(db.execute(
        //     "CREATE TABLE IF NOT EXISTS schema_migrations (
        //             id SERIAL PRIMARY KEY,
        //             version BIGINT NOT NULL,
        //             name VARCHAR(255) NOT NULL,
        //             created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now()
        //           )",
        //     &[],
        // ));
        // try!(db.execute(
        //     "CREATE UNIQUE INDEX IF NOT EXISTS idx_schema_migrations_version_name
        //         ON schema_migrations(version, name)",
        //     &[],
        // ));
        // return Ok(db);
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
        let root = Path::new("migrations").join(format!(
            "{}_{}",
            try!(time::strftime("%Y%m%d%H%M%S", &time::now_utc())),
            name
        ));
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
