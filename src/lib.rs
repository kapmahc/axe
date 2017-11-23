#![feature(plugin)]
#![plugin(rocket_codegen)]
extern crate rocket;

extern crate amqp;
extern crate base64;
extern crate clap;
#[macro_use]
extern crate log;
extern crate postgres;
extern crate rand;
extern crate redis;
#[macro_use]
extern crate serde_derive;
extern crate toml;
extern crate time;
extern crate mustache;
extern crate rustc_serialize;


pub mod env;
pub mod nut;
pub mod forum;
pub mod survey;
pub mod reading;
pub mod erp;
pub mod mall;
pub mod ops;
pub mod pos;
