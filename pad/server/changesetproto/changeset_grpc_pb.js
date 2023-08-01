// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var changeset_pb = require('./changeset_pb.js');

function serialize_changset_GenerateReply(arg) {
  if (!(arg instanceof changeset_pb.GenerateReply)) {
    throw new Error('Expected argument of type changset.GenerateReply');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_changset_GenerateReply(buffer_arg) {
  return changeset_pb.GenerateReply.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_changset_GenerateRequest(arg) {
  if (!(arg instanceof changeset_pb.GenerateRequest)) {
    throw new Error('Expected argument of type changset.GenerateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_changset_GenerateRequest(buffer_arg) {
  return changeset_pb.GenerateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}


// The greeting service definition.
var ChangesetService = exports.ChangesetService = {
  // Sends a greeting
generate: {
    path: '/changset.Changeset/Generate',
    requestStream: false,
    responseStream: false,
    requestType: changeset_pb.GenerateRequest,
    responseType: changeset_pb.GenerateReply,
    requestSerialize: serialize_changset_GenerateRequest,
    requestDeserialize: deserialize_changset_GenerateRequest,
    responseSerialize: serialize_changset_GenerateReply,
    responseDeserialize: deserialize_changset_GenerateReply,
  },
};

exports.ChangesetClient = grpc.makeGenericClientConstructor(ChangesetService);
