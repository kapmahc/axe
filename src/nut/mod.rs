pub mod users;
pub mod admin;
pub mod pages;
pub mod middlewares;
pub mod models;

use std::collections::HashMap;
use log;
use rocket::{Route, Request};
use rocket_contrib::Template;
use super::env::database::Db;
use super::env::i18n;

pub const APPLICATION_LAYOUT: &'static str = "layouts/application/index";
pub const DASHBOARD_LAYOUT: &'static str = "layouts/dashboard/index";

#[get("/")]
pub fn home(_db: Db, lng: i18n::Locale) -> &'static str {
    log::info!("{:?}", lng);
    "home"
}

#[error(404)]
fn not_found(req: &Request) -> Template {
    let mut ctx = HashMap::new();
    ctx.insert("path", req.uri().as_str());
    ctx.insert("parent", APPLICATION_LAYOUT);
    return Template::render("errors/not_found", &ctx);
}

pub fn routes() -> Vec<Route> {
    return routes![
        home,
        users::get_sign_in,
        users::post_sign_in,
        users::get_sign_up,
        users::post_sign_up,
        admin::site::get_info,
        admin::site::get_author,
        users::get_confirm,
    ];
}
