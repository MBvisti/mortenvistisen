use argon2::Config;

pub fn create_hashed_password_from_string(password_string: &str, salty: &str) -> String {
    let config = Config::default();
    argon2::hash_encoded(password_string.as_bytes(), salty.as_bytes(), &config).unwrap()
}

pub fn verify_hashed_password(hash: &str, password: &str) -> Result<bool, argon2::Error> {
    argon2::verify_encoded(hash, password.as_bytes())
}
