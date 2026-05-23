import "@typespec/http";
import "@typespec/openapi3";
import "./lib/models.tsp";
import "./lib/event-types.tsp";
import "./lib/bookings.tsp";
import "./lib/public.tsp";

using TypeSpec.Http;

@service(#{
  title: "Booking System API",
})
@server("http://localhost:8080", "Local development")
namespace BookingSystem;
