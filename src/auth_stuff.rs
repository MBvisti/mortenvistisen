use actix_identity::Identity;
use argon2::{
    password_hash::{PasswordHasher, PasswordVerifier, SaltString},
    Argon2, ParamsBuilder, PasswordHash,
};
//use rand::rngs::OsRng;
use secrecy::{ExposeSecret, Secret};

lazy_static::lazy_static! {
    pub static ref SECRET_KEY: String = std::env::var("APP_SERVER__PASSWORD_HASH_KEY").unwrap_or_else(|_| "supersecretkey".to_string());
    pub static ref SALT: String = std::env::var("APP_SERVER__PASSWORD_SALT").unwrap_or_else(|_| "salty".to_string());
}

fn build_argon2_instance() ->  Argon2<'static> {
    Argon2::new_with_secret(
        SECRET_KEY.as_bytes(),
        argon2::Algorithm::Argon2id,
        argon2::Version::V0x13,
        ParamsBuilder::new().params().unwrap(),
    )
    .unwrap()
}

pub fn hash_password(password: Secret<String>) -> String {
    let salt = SaltString::b64_encode(SALT.as_bytes()).unwrap();

    build_argon2_instance()
        .hash_password(password.expose_secret().as_bytes(), &salt)
        .unwrap()
        .to_string()
}

pub fn verify_password(password: &str, hashed_password: &str) -> bool {
    let hp = PasswordHash::new(hashed_password).unwrap();
    build_argon2_instance()
        .verify_password(password.as_bytes(), &hp)
        .is_ok()
}

pub fn validate_cookie_identity(identity: Option<Identity>) -> Result<String, String> {
    let id = match identity {
        Some(id) => id,
        None => return Err("not validated".to_string()),
    };

    match id.id() {
        Ok(id) => Ok(id),
        Err(e) => Err(e.to_string()),
    }
}
