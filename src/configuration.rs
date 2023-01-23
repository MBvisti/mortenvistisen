use config;
use secrecy::Secret;
use serde_aux::field_attributes::deserialize_number_from_string;

#[derive(serde::Deserialize)]
pub struct AppConfiguration {
    pub database: DatabaseConfig,
    pub server: ServerConfig,
    pub email_client: EmailClientConfig,
}

pub enum Env {
    Dev,
    Test,
    Prod,
}

impl Env {
    pub fn as_str(&self) -> &'static str {
        match self {
            Env::Dev => "dev",
            Env::Test => "test",
            Env::Prod => "prod",
        }
    }
}

impl TryFrom<String> for Env {
    type Error = String;

    fn try_from(value: String) -> Result<Self, Self::Error> {
        match value.to_lowercase().as_str() {
            "dev" => Ok(Self::Dev),
            "test" => Ok(Self::Test),
            "prod" => Ok(Self::Prod),
            other => Err(format!("{} is not a supported env", other)),
        }
    }
}

#[derive(serde::Deserialize)]
pub struct ServerConfig {
    pub host: String,
    #[serde(deserialize_with = "deserialize_number_from_string")]
    pub port: u16,
    pub base_url: String,
}

#[derive(serde::Deserialize)]
pub struct EmailClientConfig {
    pub api_base_url: String,
    pub app_base_url: String,
    pub sender: String,
    pub auth_token: Secret<String>,
}

#[derive(serde::Deserialize)]
pub struct DatabaseConfig {
    pub host: String,
    pub username: String,
    pub password: String,
    #[serde(deserialize_with = "deserialize_number_from_string")]
    pub port: u16,
    pub name: String,
    pub ssl_mode: String,
}

impl DatabaseConfig {
    pub fn get_connection_string(&self) -> String {
        format!(
            "postgresql://{}:{}@{}:{}/{}?sslmode={}",
            self.username, self.password, self.host, self.port, self.name, self.ssl_mode
        )
    }
}

pub fn get_config() -> Result<AppConfiguration, config::ConfigError> {
    let config = match config::Config::builder()
        .add_source(
            config::Environment::with_prefix("app")
                .prefix_separator("_")
                .separator("__"),
        )
        .build()
    {
        Ok(cfg) => cfg,
        Err(e) => {
            println!("{}", e);
            return Err(e);
        }
    };

    config.try_deserialize()
}
