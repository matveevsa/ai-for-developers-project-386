using TypeSpec.Http;

@minLength(1)
@maxLength(100)
scalar NonEmptyString extends string;

@pattern("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
scalar EmailString extends string;

alias OwnerId = string;
alias EventTypeId = string;
alias SlotId = string;
alias BookingId = string;

model EventType {
  id: EventTypeId;
  ownerId: OwnerId;
  name: NonEmptyString;
  description?: string;
  @minValue(15)
  @maxValue(480)
  duration: int32;
}

model EventTypeCreate {
  name: NonEmptyString;
  description?: string;
  @minValue(15)
  @maxValue(480)
  duration: int32;
}

model Slot {
  id: SlotId;
  eventTypeId: EventTypeId;
  startTime: utcDateTime;
  endTime: utcDateTime;
}

model Booking {
  id: BookingId;
  slotId: SlotId;
  eventTypeId: EventTypeId;
  guestName: NonEmptyString;
  guestEmail: EmailString;
  guestNotes?: string;
}

model BookingCreate {
  slotId: SlotId;
  eventTypeId: EventTypeId;
  guestName: NonEmptyString;
  guestEmail: EmailString;
  guestNotes?: string;
}

model ErrorResponse {
  code: int32;
  message: string;
}
