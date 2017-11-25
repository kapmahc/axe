use std::fs::{File, create_dir_all, OpenOptions};
use std::os::unix::fs::OpenOptionsExt;
use std::path::{Path, PathBuf};
use std::io::Read;
use std::env;
use time;
use rocket;
use mustache;
use super::errors::{Error, Result};
use super::config;

pub struct App {
    name: String,
}


impl App {
    pub fn new(name: String) -> App {
        App { name: name }
    }

    pub fn database_migrate(&self) -> Result<()> {
        return Ok(());
    }

    pub fn database_rollback(&self) -> Result<()> {
        return Ok(());
    }

    pub fn database_version(&self) -> Result<()> {
        return Ok(());
    }

    pub fn generate_nginx_conf(&self, ssl: bool) -> Result<()> {
        let file = Path::new("tmp").join("nginx.conf");
        match file.parent() {
            None => {}
            Some(d) => try!(create_dir_all(d)),
        }
        let cfg = try!(config::Config::load(&self.name));
        let scheme = if ssl { "https" } else { "http" };

        info!(
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

    pub fn start_server(&self) -> Result<()> {
        let cfg = try!(
            rocket::config::ConfigBuilder::from(try!(config::Config::load(&self.name)))
                .finalize()
        );
        let app = rocket::custom(cfg, false);
        // TODO init router
        return Err(Error::from(app.launch()));
    }


    fn create_empty_file(&self, name: PathBuf) -> Result<()> {
        info!("generate file {:?}", name.as_os_str());
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
