use reqwest::Client;
use secrecy::{ExposeSecret, Secret};

use crate::entities::Email;

pub struct EmailClient {
    http_client: Client,
    api_base_url: String,
    sender: Email,
    auth_token: Secret<String>,
    pub app_base_url: String,
}

impl EmailClient {
    pub fn new(
        api_base_url: String,
        app_base_url: String,
        sender: Email,
        auth_token: Secret<String>,
    ) -> Self {
        Self {
            http_client: Client::builder()
                .timeout(std::time::Duration::from_secs(10))
                .build()
                .unwrap(),
            api_base_url,
            sender,
            auth_token,
            app_base_url,
        }
    }
    pub async fn send_email(
        &self,
        recipient: Email,
        subject: &str,
        html_content: &str,
    ) -> Result<(), reqwest::Error> {
        let url = format!("{}/email", self.api_base_url);
        let req_body = SendEmailReq {
            from: self.sender.as_ref(),
            to: recipient.as_ref(),
            subject,
            html_body: html_content,
        };
        match self
            .http_client
            .post(&url)
            .header("X-Postmark-Server-Token", self.auth_token.expose_secret())
            .json(&req_body)
            .send()
            .await?
            .error_for_status()
        {
            Ok(_) => (),
            Err(e) => {
                println!("this is the err: {e}");
                return Err(e);
            }
        };
        Ok(())
    }
}

#[derive(serde::Serialize)]
#[serde(rename_all = "PascalCase")]
struct SendEmailReq<'a> {
    from: &'a str,
    to: &'a str,
    subject: &'a str,
    html_body: &'a str,
}
