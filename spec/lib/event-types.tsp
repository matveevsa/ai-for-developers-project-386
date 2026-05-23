using TypeSpec.Http;

@route("/api/event-types")
interface EventTypes {
  @post
  @summary("Создать тип события")
  op create(
    @body body: EventTypeCreate,
  ): EventType | ErrorResponse;

  @get
  @summary("Список всех типов событий")
  op list(): EventType[];
}
