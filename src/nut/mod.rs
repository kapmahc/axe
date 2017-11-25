pub mod users;
pub mod admin;

use rocket::Route;

#[get("/")]
pub fn home() -> &'static str {
    "home"
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
