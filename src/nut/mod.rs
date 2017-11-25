pub mod users;
pub mod admin;

#[get("/")]
pub fn home() -> &'static str {
    "home"
}
