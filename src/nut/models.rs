use time;

#[derive(Queryable)]
pub struct Locale {
    pub id: usize,
    pub code: String,
    pub message: String,
    pub created_at: time::Timespec,
    pub updated_at: time::Timespec,
}
