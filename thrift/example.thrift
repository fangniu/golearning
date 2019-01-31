namespace go example

struct ProtoRequest {
    1: i32 type
    2: binary content
    3: i64 sharding_id
}
