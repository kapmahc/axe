

use rocket::Route;

#[get("/")]
pub fn home() -> &'static str {
    "survey home"
}

pub fn routes() -> Vec<Route> {
    return routes![home];
}
