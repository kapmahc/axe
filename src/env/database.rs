use std::ops::Deref;
use diesel::pg::PgConnection;
use r2d2;
use r2d2_diesel::ConnectionManager;
use rocket::http::Status;
use rocket::request::{self, FromRequest};
use rocket::{Request, State, Outcome};
use super::errors::Result;


type Pool = r2d2::Pool<ConnectionManager<PgConnection>>;


pub fn open(url: String) -> Result<Pool> {
    let cfg = r2d2::Config::default();
    let cmg = ConnectionManager::<PgConnection>::new(url);
    let pool = try!(r2d2::Pool::new(cfg, cmg));
    return Ok(pool);
}


pub struct Db(pub r2d2::PooledConnection<ConnectionManager<PgConnection>>);

impl<'a, 'r> FromRequest<'a, 'r> for Db {
    type Error = ();

    fn from_request(request: &'a Request<'r>) -> request::Outcome<Db, ()> {
        let pool = request.guard::<State<Pool>>()?;
        match pool.get() {
            Ok(conn) => Outcome::Success(Db(conn)),
            Err(_) => Outcome::Failure((Status::ServiceUnavailable, ())),
        }
    }
}

impl Deref for Db {
    type Target = PgConnection;

    fn deref(&self) -> &Self::Target {
        &self.0
    }
}
