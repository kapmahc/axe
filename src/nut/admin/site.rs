#[get("/info")]
pub fn get_info() -> &'static str {
    "site info"
}

#[get("/author")]
pub fn get_author() -> &'static str {
    "site author"
}
