use rocket::{config, custom, Rocket};
use super::super::nut;

pub fn mount(cfg: config::Config) -> Rocket {
    let app = custom(cfg, false);
    return app.mount("/", routes![nut::home])
        .mount(
            "/users",
            routes![
            nut::users::get_sign_in,
            nut::users::post_sign_in,
            nut::users::get_sign_up,
            nut::users::post_sign_up,
        ],
        )
        .mount(
            "/admin/site",
            routes![
        nut::admin::site::get_info,
        nut::admin::site::get_author,
        ],
        );
}
