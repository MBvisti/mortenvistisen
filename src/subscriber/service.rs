use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};

pub fn generate_subcription_token() -> String {
    let mut rng = thread_rng();
    std::iter::repeat_with(|| rng.sample(Alphanumeric))
        .map(char::from)
        .take(25)
        .collect()
}
