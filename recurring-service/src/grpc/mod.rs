use chrono::NaiveDate;
use std::str::FromStr;
use tonic::{Request, Response, Status};

use crate::models::events::{CreateEvent as CreateEventModel, UpdateEvent as UpdateEventModel};
use crate::services::events_service::EventsService;

pub mod pb {
    pub mod gobizz {
        pub mod events {
            pub mod v1 {
                tonic::include_proto!("gobizz.events.v1");
            }
        }
    }
}
use pb::gobizz::events::v1::*;

pub struct EventsGrpc {
    svc: EventsService,
}

impl EventsGrpc {
    pub fn new(svc: EventsService) -> Self {
        Self { svc }
    }
}

fn parse_uuid(id: &str) -> Result<uuid::Uuid, Status> {
    tracing::info!("Parsing UUID: {}", id);
    uuid::Uuid::from_str(id).map_err(|_| Status::invalid_argument("invalid UUID"))
}
fn parse_date(s: &str) -> Result<NaiveDate, Status> {
    NaiveDate::parse_from_str(s, "%Y-%m-%d")
        .map_err(|_| Status::invalid_argument("invalid date (expected YYYY-MM-DD)"))
}

fn event_to_proto(e: crate::models::events::Event) -> Event {
    Event {
        id: e.id.to_string(),
        customer_id: e.customer_id.to_string(),
        name: e.name,
        start_date: e.start_date.format("%Y-%m-%d").to_string(),
        interval_days: e.interval_days,
        stop_at: e.stop_at.map(|d| d.format("%Y-%m-%d").to_string()),
        created_at: Some(prost_types::Timestamp {
            seconds: e.created_at.timestamp(),
            nanos: e.created_at.timestamp_subsec_nanos() as i32,
        }),
        updated_at: Some(prost_types::Timestamp {
            seconds: e.updated_at.timestamp(),
            nanos: e.updated_at.timestamp_subsec_nanos() as i32,
        }),
    }
}

#[tonic::async_trait]
impl events_server::Events for EventsGrpc {
    async fn create_event(
        &self,
        request: Request<CreateEventRequest>,
    ) -> Result<Response<Event>, Status> {
        let r = request.into_inner();

        let customer_id = parse_uuid(&r.customer_id)?;
        let start = parse_date(&r.start_date)?;
        let stop_at = r.stop_at.map(|s| parse_date(&s)).transpose()?;

        let model = CreateEventModel {
            name: r.name,
            start_date: start,
            interval_days: r.interval_days,
            stop_at,
        };

        let ev = self
            .svc
            .create(customer_id, model)
            .await
            .map_err(|e| Status::invalid_argument(e.to_string()))?;
        Ok(Response::new(event_to_proto(ev)))
    }

    async fn get_event(
        &self,
        request: Request<GetEventRequest>,
    ) -> Result<Response<Event>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let id = parse_uuid(&r.id)?;
        match self.svc.get(customer_id, id).await {
            Ok(Some(ev)) => Ok(Response::new(event_to_proto(ev))),
            Ok(None) => Err(Status::not_found("event not found")),
            Err(e) => Err(Status::internal(e.to_string())),
        }
    }

    async fn list_events(
        &self,
        request: Request<ListEventsRequest>,
    ) -> Result<Response<ListEventsResponse>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let limit = if r.limit <= 0 { 20 } else { r.limit.min(100) } as i64;
        let offset = r.offset.max(0) as i64;

        let list = self
            .svc
            .list(customer_id, limit, offset)
            .await
            .map_err(|e| Status::internal(e.to_string()))?;
        Ok(Response::new(ListEventsResponse {
            events: list.into_iter().map(event_to_proto).collect(),
        }))
    }

    async fn update_event(
        &self,
        request: Request<UpdateEventRequest>,
    ) -> Result<Response<Event>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let id = parse_uuid(&r.id)?;

        let start_date = r
            .start_date
            .map(|s| parse_date(&s))
            .transpose()
            .map_err(|e| Status::invalid_argument(e.message().to_string()))?;
        let interval_days = r.interval_days.map(|v| v);

        let stop_at = match r.stop_at_change {
            Some(update_event_request::StopAtChange::StopAt(d)) => Some(Some(parse_date(&d)?)),
            Some(update_event_request::StopAtChange::ClearStopAt(true)) => Some(None),
            _ => None,
        };

        let patch = UpdateEventModel {
            name: r.name.map(|v| v),
            start_date,
            interval_days,
            stop_at,
        };

        match self.svc.update(customer_id, id, patch).await {
            Ok(Some(ev)) => Ok(Response::new(event_to_proto(ev))),
            Ok(None) => Err(Status::not_found("event not found")),
            Err(e) => Err(Status::invalid_argument(e.to_string())),
        }
    }

    async fn delete_event(
        &self,
        request: Request<DeleteEventRequest>,
    ) -> Result<Response<DeleteEventResponse>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let id = parse_uuid(&r.id)?;
        let deleted = self
            .svc
            .delete(customer_id, id)
            .await
            .map_err(|e| Status::internal(e.to_string()))?;
        Ok(Response::new(DeleteEventResponse { deleted }))
    }

    async fn cut_event_from(
        &self,
        request: Request<CutEventFromRequest>,
    ) -> Result<Response<Event>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let id = parse_uuid(&r.id)?;
        let from = parse_date(&r.from)?;
        match self.svc.cut_from(customer_id, id, from).await {
            Ok(Some(ev)) => Ok(Response::new(event_to_proto(ev))),
            Ok(None) => Err(Status::not_found("event not found")),
            Err(e) => Err(Status::invalid_argument(e.to_string())),
        }
    }

    async fn list_occurrences(
        &self,
        request: Request<ListOccurrencesRequest>,
    ) -> Result<Response<ListOccurrencesResponse>, Status> {
        let r = request.into_inner();
        let customer_id = parse_uuid(&r.customer_id)?;
        let start = parse_date(&r.start)?;
        let end = parse_date(&r.end)?;
        let name = r.name.map(|v| v);

        let occs = self
            .svc
            .list_occurrences(customer_id, start, end, name)
            .await
            .map_err(|e| Status::invalid_argument(e.to_string()))?;

        Ok(Response::new(ListOccurrencesResponse {
            start: start.format("%Y-%m-%d").to_string(),
            end: end.format("%Y-%m-%d").to_string(),
            occurrences: occs
                .into_iter()
                .map(|o| Occurrence {
                    event_id: o.event_id.to_string(),
                    name: o.name,
                    date: o.date.format("%Y-%m-%d").to_string(),
                })
                .collect(),
        }))
    }
}
