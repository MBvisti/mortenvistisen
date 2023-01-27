use serde::{Deserialize, Serialize};
use validator::validate_email;

use super::handlers::SubscribeFormData;

#[derive(Debug, Deserialize)]
pub struct Email(String);

impl Email {
    pub fn parse(s: String) -> Result<Email, String> {
        match validate_email(&s) {
            true => Ok(Self(s)),
            false => Err(format!("{s} not a valid email")),
        }
    }
}

impl AsRef<str> for Email {
    fn as_ref(&self) -> &str {
        &self.0
    }
}

pub struct NewSubscriberPayload {
    pub email: Email,
    pub referer: String,
}

impl TryFrom<SubscribeFormData> for NewSubscriberPayload {
    type Error = String;

    fn try_from(value: SubscribeFormData) -> Result<Self, Self::Error> {
        let email = Email::parse(value.email)?;
        Ok(Self {
            email,
            referer: value.referer,
        })
    }
}

#[derive(Debug, Serialize)]
pub struct ConfirmSubEmail {
    pub app_base_url: String,
    pub token: String,
}
