# entco

[![Build Status](https://travis-ci.org/entco/entco.svg?branch=master)](https://travis-ci.org/entco/entco)
[![GoDoc](https://godoc.org/github.com/entco/entco?status.svg)](https://godoc.org/github.com/entco/entco)
[![Go Report Card](https://goreportcard.com/badge/github.com/entco/entco)](https://goreportcard.com/report/github.com/entco/entco)

components and tools for knockout applications

## Ent && Graphql Extensions

### Ent Type System and CodeGen:

- field.Decimal: base on Field.String to support Validate
           no matter Field.Other or Field.Float, it will be some Generator error in Ent Codegen
- Audit Fields: support audit fields 
- Soft Delete: support soft delete
- Tenant: support tenant and data isolation
- PrimaryKey: support snowflake id.
- Resolver Code: support resolver code generator

#### Ent Client:

OTEL Tracing: support opentelemetry tracing for ent client
Route Driver: support route driver for ent client

### Graphql:

- GloablID: support globalID for ent relay
- SimplePagination: support pagination by page limit and offset,but also base on cursor 