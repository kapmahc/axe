use std::fs::{create_dir_all, OpenOptions};
use std::io::Write;
use std::os::unix::fs::OpenOptionsExt;
use std::path::{Path, PathBuf};
use time;
use super::errors::Result;

pub struct App<'a> {
    config: &'a str,
}

impl<'a> App<'a> {
    pub fn new(name: &'a str) -> App {
        App { config: name }
    }

    pub fn generate_locale(&self, name: &'a str) -> Result<()> {
        let mut _fn = Path::new("locales").join(name);
        _fn.set_extension("yaml");
        return self.create_empty_file(_fn);
    }

    pub fn generate_migration(&self, name: &'a str) -> Result<()> {
        let fmt = "%Y%m%d%H%M%S".to_string() + "-" + name;
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
