use std::sync::atomic::AtomicUsize;

#[derive(Debug)]
pub struct Status {
    hit_count: AtomicUsize,
}

impl Status {
    pub fn new() -> Status {
        return Status { hit_count: AtomicUsize::new(0) };
    }
}
