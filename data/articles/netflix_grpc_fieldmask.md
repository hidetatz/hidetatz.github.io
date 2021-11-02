type: input
timestamp: 2021-11-02 23:58:20
url: https://netflixtechblog.com/practical-api-design-at-netflix-part-1-using-protobuf-fieldmask-35cfdc606518
lang: en
---


* We sometimes want to specify which fields are needed/not needed when calling remote API- like GraphQL.
* It is important because remote calls are not free. It imposes extra latency, error probability is increased, and consumes network bandwidth
* In JSON API, we can use sparse fieldsets
    * https://jsonapi.org/format/#fetching-sparse-fieldsets
* In gRPC, we can use FieldMask

```proto
// feld_mask.proto
message FieldMask {
  // The set of field mask paths.
  repeated string paths = 1;
}
```

* Let's say we want to call `GerProduction` API but we don't need full response from the API, we can use FieldMask like this:

```
import "google/protobuf/field_mask.proto";

message GetProductionRequest {
  string production_id = 1;
  google.protobuf.FieldMask field_mask = 2;
}
```

* If a client needs only `title` and `format`, it can build a request like this:

```
FieldMask fieldMask = FieldMask.newBuilder()
    .addPaths("title")
    .addPaths("format")
    .build();

GetProductionRequest request = GetProductionRequest.newBuilder()
    .setProductionId(LA_CASA_DE_PAPEL_PRODUCTION_ID)
    .setFieldMask(fieldMask)
    .build();
```

* Serverside implementation:

```
private static final String FIELD_SEPARATOR_REGEX = "\\.";
private static final String MAX_FIELD_NESTING = 2;
private static final String SCHEDULE_FIELD_NAME =                                // (1)
    Production.getDescriptor()
    .findFieldByNumber(Production.SCHEDULE_FIELD_NUMBER).getName();

@Override
public void getProduction(GetProductionRequest request, 
                          StreamObserver<GetProductionResponse> response) {

    FieldMask canonicalFieldMask =                                               
        FieldMaskUtil.normalize(request.getFieldMask());                         // (2) 

    boolean scheduleFieldRequested =                                             // (3)
        canonicalFieldMask.getPathsList().stream()
            .map(path -> path.split(FIELD_SEPARATOR_REGEX, MAX_FIELD_NESTING)[0])
            .anyMatch(SCHEDULE_FIELD_NAME::equals);

    if (scheduleFieldRequested) {
        ProductionSchedule schedule = 
            makeExpensiveCallToScheduleService(request.getProductionId());       // (4)
        ...
    }

    ...
}
```

