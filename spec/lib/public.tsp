using TypeSpec.Http;

namespace BookingSystem;

@route("/api/public/event-types")
interface PublicEventTypes {
  @get
  @summary("Публичный список видов брони")
  op list(): EventType[];
}

@route("/api/public/event-types/{eventTypeId}/slots")
interface PublicAvailableSlots {
  @get
  @summary("Свободные слоты на дату (до +14 дней)")
  op getAvailableSlots(
    @path eventTypeId: EventTypeId,
    @query date: utcDateTime,
  ): Slot[];
}

@route("/api/public/bookings")
interface PublicBookings {
  @post
  @summary("Создать бронирование (гость)")
  op create(
    @body body: BookingCreate,
  ): Booking | ErrorResponse;
}
