extern crate axe;
extern crate env_logger;

fn main() {
    env_logger::init().unwrap();
    axe::env::run();
}
