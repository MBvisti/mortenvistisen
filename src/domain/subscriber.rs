use serde::Deserialize;
use validator::validate_email;

#[derive(Debug, Deserialize)]
pub struct Email(String);

impl Email {
    pub fn parse(s: String) -> Result<Email, String> {
        match validate_email(&s) {
            true => Ok(Self(s)),
            false => Err(format!("{} not a valid email", s)),
        }
    }
}

impl AsRef<str> for Email {
    fn as_ref(&self) -> &str {
        &self.0
    }
}

