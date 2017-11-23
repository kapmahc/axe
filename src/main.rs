
// #[get("/")]
// fn index() -> &'static str {
//     "Hello, world!"
// }
//
// #[error(404)]
// fn not_found() -> &'static str {
//     "404"
// }
//
// fn main() {
//     rocket::ignite()
//         .mount("/", routes![index])
//         .catch(errors![not_found])
//         .launch();
// }

extern crate axe;
extern crate env_logger;

fn main() {
    env_logger::init().unwrap();
    axe::env::run();
}
