use std::collections::HashMap;
use log;
use url::Url;
use time::Duration;
use rocket::request::{self, FromRequest};
use rocket::{Request, Outcome};
use rocket::http::Cookie;


#[derive(Debug)]
pub struct Locale(String);


impl<'a, 'r> FromRequest<'a, 'r> for Locale {
    type Error = ();

    fn from_request(req: &'a Request<'r>) -> request::Outcome<Locale, ()> {
        let key = "locale";
        // from query
        if let Ok(u) = Url::parse(&format!("http://localhost{}", req.uri().as_str())) {
            let params: HashMap<_, _> = u.query_pairs().into_owned().collect();
            if let Some(lng) = params.get(key) {
                req.cookies().add(
                    Cookie::build(key, lng.to_string())
                        .path("/")
                        .http_only(false)
                        .max_age(Duration::weeks(10))
                        .finish(),
                );
                return Outcome::Success(Locale(lng.to_string()));
            }
        }
        // from cookie
        if let Some(lng) = req.cookies().get(key) {
            return Outcome::Success(Locale(lng.value().to_string()));
        }
        // from header
        if let Some(lng) = req.headers().get_one("Accept-Language") {
            let langs: Vec<&str> = lng.split(',').collect();
            if let Some(lng) = langs.first() {
                return Outcome::Success(Locale(lng.to_string()));
            }
        }
        log::warn!("fail to detect language");
        return Outcome::Success(Locale("en-US".to_string()));
    }
}
