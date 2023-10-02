use actix_web::{get, post, web,  Responder};
use serde::{Deserialize, Serialize};
use sqlx::PgPool;
use uuid::Uuid;

use crate::entities::{Email, NewSubscriberPayload};
use crate::views::{render_internal_error_tmpl, render_template};
use crate::{repository, views};
use rand::distributions::Alphanumeric;
use rand::{thread_rng, Rng};

fn generate_subcription_token() -> String {
    let mut rng = thread_rng();
    std::iter::repeat_with(|| rng.sample(Alphanumeric))
        .map(char::from)
        .take(25)
        .collect()
}

#[derive(Debug, Serialize)]
struct SubscribeMetaData {
    has_error: bool,
    error_msg: Option<String>,
    title: Option<String>,
}

#[derive(Debug, Deserialize)]
pub struct SubscribeFormData {
    pub email: String,
    pub referer: String,
}

impl TryInto<NewSubscriberPayload> for SubscribeFormData {
    type Error = String;

    fn try_into(self) -> Result<NewSubscriberPayload, Self::Error> {
        let email = Email::parse(self.email)?;
        Ok(NewSubscriberPayload {
            email,
            referer: self.referer,
        })
    }
}

#[post("/subscribe")]
pub async fn subscribe_to_newsletter(
    // req: HttpRequest,
    form: web::Form<SubscribeFormData>,
    pool: web::Data<PgPool>,
    // email_client: web::Data<EmailClient>,
) -> impl Responder {
    let payload: NewSubscriberPayload = match form.0.try_into() {
        Ok(p) => p,
        Err(e) => {
            println!("there was an error getting new sub payload: {e:?}");
            let sub_res_view = views::SubscribeResponse::new(views::SubscribeResponseData {
                has_error: true,
                error_msg: "could not subscribe you, please try again".to_string(),
            });

            return render_template(sub_res_view);
        }
    };

    let email_exists = match repository::does_email_exists(&pool, &payload.email).await {
        Ok(exists) => exists,
        Err(e) => {
            println!("there was an error checking if email exists: {e:?}");
            return render_internal_error_tmpl(None);
        }
    };

    if email_exists {
        let sub_res_view = views::SubscribeResponse::new(views::SubscribeResponseData {
            has_error: true,
            error_msg: "that email is already registered".to_string(),
        });

        return render_template(sub_res_view);
    }

    let new_sub_id = Uuid::new_v4();
    if let Err(e) =
        repository::create_new_subscriber(&pool, &payload.email, &payload.referer, new_sub_id).await
    {
        println!("{e:?}");
        return render_internal_error_tmpl(None);
    };

    let new_sub_token = generate_subcription_token();
    if let Err(e) = repository::create_new_subscriber_token(&pool, new_sub_id, &new_sub_token).await
    {
        println!("{e:?}");
        return render_internal_error_tmpl(None);
    };

    // email_context.insert(
    //     "meta_data",
    //     &ConfirmSubEmail {
    //         app_base_url: email_client.app_base_url.clone(),
    //         token: new_sub_token,
    //     },
    // );
    // let html_content = match render_template("confirm_sub_email.html", &email_context) {
    //     Ok(t) => t,
    //     Err(_) => {
    //         return HttpResponse::InternalServerError()
    //             .content_type("text/html")
    //             .body(render_internal_error_tmpl(None))
    //     }
    // };
    // match email_client
    //     .send_email(
    //         payload.email,
    //         "Thanks for subscribing to my newsletter!",
    //         &html_content,
    //     )
    //     .await
    // {
    //     Ok(_) => (),
    //     Err(e) => {
    //         println!("{e:?}");
    //         return render_subscribe_err(
    //             "I fucked up somehow, sorry. Please try again".to_string(),
    //             None,
    //             context,
    //             "_subscribe_response.html",
    //         )
    //         .await;
    //     }
    // };

    let sub_res_view = views::SubscribeResponse::new(views::SubscribeResponseData {
        has_error: false,
        error_msg: "".to_string(),
    });

    render_template(sub_res_view)
}

#[derive(serde::Deserialize)]
pub struct Parameters {
    token: String,
}

#[tracing::instrument(
    name = "verifying subscriber", 
    skip(pool, params)
    fields(
        token = %params.token
    )
)]
#[get("/subscribe/verify")]
pub async fn verify_subscription(
    pool: web::Data<PgPool>,
    params: web::Query<Parameters>,
) -> impl Responder {
    let id = match repository::get_subscriber_id_from_token(&pool, &params.token).await {
        Ok(id) => id,
        Err(e) => {
            println!("{e:?}");
            return render_internal_error_tmpl(None);
        }
    };

    let is_verified = match repository::is_user_verified(&pool, id).await {
        Ok(is_verified) => is_verified,
        Err(e) => {
            println!("{e:?}");
            return render_internal_error_tmpl(None);
        }
    };

    if is_verified {
        let subscribe_verify_view = views::SubscribeVerify::new(views::SubscribeVerifyData {
            has_error: false,
            already_verified: true,
            error_msg: None,
        });
        return render_template(subscribe_verify_view);
    }

    match repository::update_email_to_verified(&pool, id).await {
        Ok(_) => (),
        Err(e) => {
            println!("{e:?}");
            return render_internal_error_tmpl(None);
        }
    };

    let sub_verify_view = views::SubscribeVerify::new(views::SubscribeVerifyData {
        has_error: false,
        already_verified: false,
        error_msg: None,
    });
    render_template(sub_verify_view)
}

// #[get("/subscribe/delete")]
// pub async fn delete_subscriber(
//     // tmpl: web::Data<tera::Tera>,
//     pool: web::Data<PgPool>,
//     params: web::Query<Parameters>,
// ) -> impl Responder {
//     let context = tera::Context::new();

//     let id = match repository::get_subscriber_id_from_token(&pool, &params.token).await {
//         Ok(id) => id,
//         Err(e) => {
//             println!("{e:?}");
//             match e {
//                 // TODO: add message that it has already been deleted
//                 sqlx::Error::RowNotFound => {
//                     println!("{e:?}");
//                     return render_subscribe_err(
//                         "Email already deleted".to_string(),
//                         Some(String::from("You're all set")),
//                         context,
//                         "confirm_subscription.html",
//                     )
//                     .await;
//                 }
//                 _ => {
//                     return render_subscribe_err(
//                         "I fucked up somehow, sorry. Please try again".to_string(),
//                         None,
//                         context,
//                         "delete_subscription.html",
//                     )
//                     .await
//                 }
//             }
//         }
//     };

//     match repository::delete_subscriber_token(&pool, &params.token).await {
//         Ok(_) => (),
//         Err(e) => {
//             println!("{e:?}");
//             return render_subscribe_err(
//                 "I fucked up somehow, sorry. Please try again".to_string(),
//                 None,
//                 context,
//                 "delete_subscription.html",
//             )
//             .await;
//         }
//     };

//     match repository::delete_subscriber(&pool, id).await {
//         Ok(_) => (),
//         Err(e) => {
//             println!("{e:?}");
//             return render_subscribe_err(
//                 "I fucked up somehow, sorry. Please try again".to_string(),
//                 None,
//                 context,
//                 "delete_subscription.html",
//             )
//             .await;
//         }
//     };

//     render_subscribe_err(
//         "Email deleted".to_string(),
//         Some(String::from("You're all set")),
//         context,
//         "confirm_subscription.html",
//     )
//     .await
// }
