#[get("/admin/site/info")]
pub fn get_info() -> &'static str {
    "site info"
}

#[get("/admin/site/author")]
pub fn get_author() -> &'static str {
    "site author"
}
