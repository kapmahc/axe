#[get("/sign-in")]
pub fn get_sign_in() -> &'static str {
    "sign in"
}

#[post("/sign-in")]
pub fn post_sign_in() -> &'static str {
    "sign in"
}

#[get("/sign-up")]
pub fn get_sign_up() -> &'static str {
    "sign up"
}

#[post("/sign-up")]
pub fn post_sign_up() -> &'static str {
    "sign up"
}
