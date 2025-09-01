use std::net::{Ipv6Addr, SocketAddr};
use tonic::transport::Server;
use tonic_health::server::health_reporter;

use crate::config::AppConfig;
use crate::db::connect;
use crate::repositories::events_repository::EventsRepository;
use crate::services::events_service::EventsService;

mod config;
mod db;
mod grpc;
mod models;
mod repositories;
mod services;

use grpc::EventsGrpc;
use grpc::pb::gobizz::events::v1::events_server::EventsServer;
use tracing_subscriber::{EnvFilter, fmt};

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    dotenvy::dotenv().ok();

    let filter = EnvFilter::try_from_default_env().unwrap_or_else(|_| EnvFilter::new("info"));
    fmt()
        .with_env_filter(filter)
        .with_target(true)
        .compact()
        .init();

    let cfg = AppConfig::from_env()?;
    let pool = connect(&cfg.database_url).await?;

    let repo = EventsRepository::new(pool.clone());
    let svc = EventsService::new(repo);
    let events_grpc = EventsGrpc::new(svc);
    tracing::info!("DB connected");

    let (health_reporter, health_service) = health_reporter();
    health_reporter
        .set_serving::<EventsServer<EventsGrpc>>()
        .await;

    let addr: SocketAddr = SocketAddr::from((Ipv6Addr::UNSPECIFIED, cfg.port));
    tracing::info!("gRPC listening on {}", addr);

    Server::builder()
        .add_service(health_service)
        .add_service(EventsServer::new(events_grpc))
        .serve(addr)
        .await?;

    Ok(())
}
