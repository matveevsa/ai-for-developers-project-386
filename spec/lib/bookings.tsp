using TypeSpec.Http;

@route("/api/owners/{ownerId}/bookings")
interface OwnerBookings {
  @get
  @summary("Предстоящие встречи владельца")
  op listUpcoming(
    @path ownerId: OwnerId,
  ): Booking[];
}
