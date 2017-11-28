#[get("/users/sign-in")]
pub fn get_sign_in() -> &'static str {
    "sign in"
}

#[post("/users/users/sign-in")]
pub fn post_sign_in() -> &'static str {
    "sign in"
}

#[get("/users/sign-up")]
pub fn get_sign_up() -> &'static str {
    "sign up"
}

#[post("/users/sign-up")]
pub fn post_sign_up() -> &'static str {
    "sign up"
}



#[get("/users/confirm")]
pub fn get_confirm() -> &'static str {
    "confirm"
}
