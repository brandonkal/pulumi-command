Execute a Command and save it as a resource.

Each command can be specified as an object or a convenience array. If only `create` is specified, `update` will use the create definition. The `compare` property will be JSON serialized with a hash saved in the state. This is useful to ensure update is run if dependendent resources change.

An update will occur in these cases:
1. The `compare` hash or the `update` arguments change.
2. The specified `diff` command exits with an error.