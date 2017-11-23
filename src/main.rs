#![feature(plugin)]
#![plugin(rocket_codegen)]

extern crate rocket;

#[get("/")]
fn index() -> &'static str {
    "Hello, world!"
}

#[error(404)]
fn not_found() -> &'static str {
    "404"
}

fn main() {
    rocket::ignite()
        .mount("/", routes![index])
        .catch(errors![not_found])
        .launch();
}
