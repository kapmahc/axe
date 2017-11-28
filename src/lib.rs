#![feature(plugin, custom_derive, use_extern_macros)]
#![plugin(rocket_codegen)]


extern crate rocket;
extern crate rocket_contrib;
extern crate url;
extern crate amqp;
extern crate base64;
extern crate clap;
#[macro_use(log)]
extern crate log;
extern crate postgres;
extern crate rand;
extern crate redis;
#[macro_use]
extern crate serde_derive;
extern crate serde;
extern crate toml;
extern crate time;
extern crate mustache;
extern crate rustc_serialize;
#[macro_use]
extern crate diesel;
#[macro_use]
extern crate diesel_codegen;
extern crate r2d2_diesel;
extern crate r2d2;

pub mod env;
pub mod nut;
pub mod forum;
pub mod survey;
pub mod reading;
pub mod erp;
pub mod mall;
pub mod ops;
pub mod pos;
