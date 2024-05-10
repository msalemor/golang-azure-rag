use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use reqwest::Client;
use serde::{Deserialize, Serialize};

#[get("/")]
async fn hello() -> impl Responder {
    HttpResponse::Ok().body("Hello world!")
}

#[post("/echo")]
async fn echo(req_body: String) -> impl Responder {
    HttpResponse::Ok().body(req_body)
}

async fn manual_hello() -> impl Responder {
    HttpResponse::Ok().body("Hey there!")
}

#[derive(Deserialize, Serialize)]
struct Payload {
    // Define your payload struct here
}

#[derive(Serialize)]
struct Response {
    // Define your response struct here
}

async fn CallAPI(payload: web::Json<Payload>) -> impl Responder {
    let client = Client::new();
    let res = client.post("").json(&payload).send().await.unwrap();
    //let response: Response = res.json::<Response>().await?;
    let response = Response {};
    web::Json(response)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(hello)
            .service(echo)
            .route("/hey", web::get().to(manual_hello))
    })
    .bind(("127.0.0.1", 8080))?
    .run()
    .await
}
