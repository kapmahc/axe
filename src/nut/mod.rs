pub mod users;
pub mod admin;
pub mod middlewares;
pub mod models;

use log;
use rocket::{Route, Request};
use super::env::database::Db;
use super::env::i18n;

#[get("/")]
pub fn home(_db: Db, lng: i18n::Locale) -> &'static str {
    log::info!("{:?}", lng);
    "home"
}

#[error(404)]
fn not_found(_req: &Request) -> &'static str {
    "not found"
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
    ];
}
