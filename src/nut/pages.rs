use std::path::PathBuf;
use std::collections::HashMap;
use log;
use rocket_contrib::Template;

#[get("/pages/<path..>")]
pub fn show(path: PathBuf) -> Template {
    log::info!("{:?}", path);
    let mut ctx = HashMap::new();
    ctx.insert("path", "aaa");
    return Template::render("pages/show", &ctx);
}
