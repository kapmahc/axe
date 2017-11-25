pub mod users;
pub mod admin;
pub mod middlewares;

use rocket::{Route, Request};

#[get("/")]
pub fn home() -> &'static str {
    "home"
}

#[error(404)]
fn not_found(req: &Request) -> &'static str {
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
