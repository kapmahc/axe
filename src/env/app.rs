use std::fs::{File, create_dir_all, OpenOptions};
use std::os::unix::fs::OpenOptionsExt;
use std::path::{Path, PathBuf};
use std::io::Read;
use time;
use toml;
use rocket;
use rocket::config::ConfigBuilder;
use super::errors::{Error, Result};
use super::config;

pub struct App {
    name: String,
}


impl App {
    pub fn new(name: String) -> App {
        App { name: name }
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
        let cfg = try!(ConfigBuilder::from(try!(self.load_config())).finalize());
        let app = rocket::custom(cfg, false);
        // TOOD init router
        return Err(Error::from(app.launch()));
    }

    fn load_config(&self) -> Result<config::Config> {
        info!("read config from file {:?}", self.name);
        let mut file = try!(File::open(&self.name));
        let mut data = String::new();
        try!(file.read_to_string(&mut data));
        let cfg = try!(toml::from_str(&data));
        return Ok(cfg);
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
