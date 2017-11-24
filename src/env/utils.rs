use base64;
use rand;
use rand::Rng;

pub fn random(len: usize) -> String {
    let secret: Vec<u8> = rand::thread_rng().gen_iter::<u8>().take(len).collect();
    return base64::encode(&secret);
}
