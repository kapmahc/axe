use rocket::Route;

#[get("/")]
pub fn home() -> &'static str {
    "forum home"
}

pub fn routes() -> Vec<Route> {
    return routes![home];
}
